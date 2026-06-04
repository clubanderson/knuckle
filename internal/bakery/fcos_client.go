package bakery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/projectbluefin/knuckle/internal/model"
	"github.com/projectbluefin/knuckle/internal/validate"
)

const (
	// FCOSCatalogURL is the GitHub Releases API for the fedora-sysexts/community repo.
	FCOSCatalogURL = "https://api.github.com/repos/fedora-sysexts/community/releases?per_page=100"

	// maxFCOSCatalogPages caps pagination for the FCOS catalog. The community repo has
	// 1,500+ releases (16+ pages at 100/page). We fetch up to 20 pages to avoid silently
	// missing extensions that only appear in older releases.
	maxFCOSCatalogPages = 20
)

// FCOSClient fetches the sysext catalog from fedora-sysexts/community.
// It filters releases by Fedora major version and architecture.
type FCOSClient struct {
	CatalogURL    string
	HTTP          *http.Client
	AuthToken     string
	FedoraVersion int
}

// NewFCOSClient creates a new FCOS sysext catalog client.
// fedoraVersion is obtained from fcos.FetchStreamFedoraVersion().
func NewFCOSClient(fedoraVersion int) *FCOSClient {
	return &FCOSClient{
		CatalogURL:    FCOSCatalogURL,
		HTTP:          &http.Client{Timeout: defaultTimeout},
		AuthToken:     githubTokenFromEnv(),
		FedoraVersion: fedoraVersion,
	}
}

// FetchCatalog delegates to FetchCatalogArch with amd64.
func (c *FCOSClient) FetchCatalog(ctx context.Context) ([]model.SysextEntry, error) {
	return c.FetchCatalogArch(ctx, "amd64")
}

// FetchCatalogArch fetches sysexts built for the given arch and the client's
// Fedora major version. Only releases whose tag matches both the Fedora version
// and architecture are included. Entries are deduplicated by name (newest wins).
//
// Asset download URLs point to GitHub Releases (not a CDN). The URL security
// constraints (maxSysextURLLen, HTTPS) apply unchanged.
func (c *FCOSClient) FetchCatalogArch(ctx context.Context, arch string) ([]model.SysextEntry, error) {
	var assetSuffix string
	switch arch {
	case "amd64":
		assetSuffix = "x86-64"
	case "arm64":
		assetSuffix = "arm64"
	default:
		return nil, fmt.Errorf("unsupported architecture %q: must be amd64 or arm64", arch)
	}

	fedoraStr := strconv.Itoa(c.FedoraVersion)

	const maxResponseSize = 5 << 20
	var allReleases []githubRelease
	nextURL := c.CatalogURL

	for page := 0; page < maxFCOSCatalogPages && nextURL != ""; page++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, nextURL, nil)
		if err != nil {
			return nil, fmt.Errorf("creating FCOS catalog request (page %d): %w", page+1, err)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "knuckle/1.0")
		if c.AuthToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.AuthToken)
		}

		resp, err := c.HTTP.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetching FCOS catalog (page %d): %w", page+1, err)
		}

		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("FCOS catalog returned status %d", resp.StatusCode)
		}

		body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
		_ = resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("reading FCOS catalog response (page %d): %w", page+1, err)
		}
		if int64(len(body)) >= maxResponseSize {
			return nil, fmt.Errorf("FCOS catalog response exceeds 5MB size limit")
		}

		var releases []githubRelease
		if err := json.Unmarshal(body, &releases); err != nil {
			return nil, fmt.Errorf("parsing FCOS catalog JSON (page %d): %w", page+1, err)
		}

		if len(releases) == 0 {
			break
		}
		allReleases = append(allReleases, releases...)

		nextURL, _ = parseLinkNext(resp.Header.Get("Link"))
	}

	seen := make(map[string]bool)
	sysexts := make([]model.SysextEntry, 0, len(allReleases))

	for _, rel := range allReleases {
		name, version, tagFedora, tagArch, err := ParseFCOSTagName(rel.TagName)
		if err != nil {
			continue
		}

		// Filter by Fedora version and architecture.
		if tagFedora != fedoraStr {
			continue
		}
		if tagArch != assetSuffix {
			continue
		}

		if seen[name] {
			continue
		}

		// Find the .raw asset matching our arch.
		var downloadURL, sha256sumsURL string
		for _, asset := range rel.Assets {
			switch {
			case asset.Name == "SHA256SUMS":
				sha256sumsURL = asset.BrowserDownloadURL
			case strings.HasSuffix(asset.Name, ".raw") && strings.Contains(asset.Name, assetSuffix):
				if downloadURL == "" {
					downloadURL = asset.BrowserDownloadURL
				}
			}
		}
		if downloadURL == "" {
			continue
		}
		if len(downloadURL) > maxSysextURLLen {
			continue
		}
		if sha256sumsURL != "" && len(sha256sumsURL) > maxSysextURLLen {
			sha256sumsURL = ""
		}
		if validate.SysextName(name) != nil {
			continue
		}

		seen[name] = true

		// Fetch SHA256 hash (best-effort).
		sha256Hash := ""
		if sha256sumsURL != "" {
			rawFilename := downloadURL[strings.LastIndex(downloadURL, "/")+1:]
			if h, fetchErr := c.fetchSHA256ForAsset(ctx, sha256sumsURL, rawFilename); fetchErr == nil {
				sha256Hash = h
			}
		}

		description := truncateDescription(rel.Body, 80)
		category := ""
		supportTier := ""

		if meta, ok := FCOSLookup(name); ok {
			if meta.Short != "" {
				description = meta.Short
			}
			category = meta.Category
			supportTier = meta.SupportTier
		}

		sysexts = append(sysexts, model.SysextEntry{
			Name:        name,
			Description: description,
			Version:     version,
			URL:         downloadURL,
			Sha256:      sha256Hash,
			Category:    category,
			SupportTier: supportTier,
			Selected:    false,
		})
	}

	return sysexts, nil
}

// fetchSHA256ForAsset reuses the same SHA256SUMS fetching logic as the Flatcar client.
func (c *FCOSClient) fetchSHA256ForAsset(ctx context.Context, sha256sumsURL, rawFilename string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sha256sumsURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating SHA256SUMS request: %w", err)
	}
	req.Header.Set("User-Agent", "knuckle/1.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching SHA256SUMS: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("SHA256SUMS returned HTTP %d", resp.StatusCode)
	}

	const maxSHA256Size = 64 << 10
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxSHA256Size))
	if err != nil {
		return "", fmt.Errorf("reading SHA256SUMS: %w", err)
	}

	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		baseName := fields[1]
		if idx := strings.LastIndex(baseName, "/"); idx >= 0 {
			baseName = baseName[idx+1:]
		}
		if baseName == rawFilename {
			hash := fields[0]
			if !reSHA256.MatchString(hash) {
				return "", fmt.Errorf("malformed SHA256 hash %q for %s", hash, rawFilename)
			}
			return hash, nil
		}
	}
	return "", nil
}
