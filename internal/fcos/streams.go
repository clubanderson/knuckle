// Package fcos provides Fedora CoreOS stream metadata fetching.
package fcos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// StreamsBaseURL is the public URL prefix for FCOS stream metadata.
	StreamsBaseURL = "https://builds.coreos.fedoraproject.org/streams/"
	streamTimeout  = 15 * time.Second

	// maxStreamResponseSize caps the stream metadata response at 2MB.
	maxStreamResponseSize = 2 << 20
)

// streamMetadata is the top-level JSON returned by the FCOS streams endpoint.
// Only the fields needed for version extraction are modeled.
type streamMetadata struct {
	Architectures map[string]struct {
		Artifacts map[string]struct {
			Release string `json:"release"`
		} `json:"artifacts"`
	} `json:"architectures"`
}

// FetchStreamFedoraVersion returns the Fedora major version number for a given
// FCOS stream (e.g. "stable" → 44). The major version is parsed from the first
// component of the FCOS release version string (e.g. "44.20260510.3.1" → 44).
//
// This is needed to filter fedora-sysexts/community assets by the correct
// Fedora version, since asset filenames embed the target Fedora major version.
func FetchStreamFedoraVersion(ctx context.Context, stream string) (int, error) {
	url := StreamsBaseURL + stream + ".json"
	client := &http.Client{Timeout: streamTimeout}
	return fetchStreamFedoraVersion(ctx, url, client)
}

// FetchStreamFedoraVersionFromURL fetches the Fedora major version from a
// custom URL. Used for testing with httptest servers.
func FetchStreamFedoraVersionFromURL(ctx context.Context, url string, client *http.Client) (int, error) {
	return fetchStreamFedoraVersion(ctx, url, client)
}

func fetchStreamFedoraVersion(ctx context.Context, url string, client *http.Client) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("creating FCOS stream request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "knuckle/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetching FCOS stream metadata: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("FCOS stream metadata returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxStreamResponseSize))
	if err != nil {
		return 0, fmt.Errorf("reading FCOS stream response: %w", err)
	}

	var meta streamMetadata
	if err := json.Unmarshal(body, &meta); err != nil {
		return 0, fmt.Errorf("parsing FCOS stream JSON: %w", err)
	}

	// Use x86_64 as the reference architecture for the version; all architectures
	// share the same Fedora major version within a stream.
	archData, ok := meta.Architectures["x86_64"]
	if !ok {
		return 0, fmt.Errorf("FCOS stream metadata has no x86_64 architecture data")
	}

	metalArtifact, ok := archData.Artifacts["metal"]
	if !ok {
		return 0, fmt.Errorf("FCOS stream metadata has no metal artifact")
	}

	return ParseFedoraVersion(metalArtifact.Release)
}

// ParseFedoraVersion extracts the Fedora major version from an FCOS release
// string. The format is "<fedora-major>.<date>.<minor>.<patch>" — only the
// first dot-separated component is used (e.g. "44.20260510.3.1" → 44).
func ParseFedoraVersion(release string) (int, error) {
	if release == "" {
		return 0, fmt.Errorf("empty FCOS release version")
	}
	major := release
	if idx := strings.IndexByte(release, '.'); idx > 0 {
		major = release[:idx]
	}
	v, err := strconv.Atoi(major)
	if err != nil {
		return 0, fmt.Errorf("invalid Fedora major version %q in release %q: %w", major, release, err)
	}

	const (
		minFedoraVersion = 30
		maxFedoraVersion = 99
	)
	if v < minFedoraVersion || v > maxFedoraVersion {
		return 0, fmt.Errorf("Fedora major version %d out of expected range [%d, %d]", v, minFedoraVersion, maxFedoraVersion)
	}
	return v, nil
}
