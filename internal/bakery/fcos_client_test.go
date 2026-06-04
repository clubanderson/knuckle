package bakery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFCOSClient_FetchCatalogArch_FiltersCorrectly(t *testing.T) {
	releases := []githubRelease{
		{
			TagName: "tailscale-0-1.98.4-1-44-x86-64",
			Body:    "Tailscale sysext",
			Assets: []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			}{
				{Name: "tailscale-0-1.98.4-1-44-x86-64.raw", BrowserDownloadURL: "https://github.com/fedora-sysexts/community/releases/download/tailscale-0-1.98.4-1-44-x86-64/tailscale-0-1.98.4-1-44-x86-64.raw"},
				{Name: "SHA256SUMS", BrowserDownloadURL: "https://github.com/fedora-sysexts/community/releases/download/tailscale-0-1.98.4-1-44-x86-64/SHA256SUMS"},
			},
		},
		{
			TagName: "tailscale-0-1.98.4-1-43-x86-64",
			Body:    "Tailscale for Fedora 43",
			Assets: []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			}{
				{Name: "tailscale-0-1.98.4-1-43-x86-64.raw", BrowserDownloadURL: "https://github.com/fedora-sysexts/community/releases/download/tailscale-0-1.98.4-1-43-x86-64/tailscale-0-1.98.4-1-43-x86-64.raw"},
			},
		},
		{
			TagName: "tailscale-0-1.98.4-1-44-arm64",
			Body:    "Tailscale arm64",
			Assets: []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			}{
				{Name: "tailscale-0-1.98.4-1-44-arm64.raw", BrowserDownloadURL: "https://github.com/fedora-sysexts/community/releases/download/tailscale-0-1.98.4-1-44-arm64/tailscale-0-1.98.4-1-44-arm64.raw"},
			},
		},
		{
			TagName: "docker-ce-3-29.5.2-1.fc44-44-x86-64",
			Body:    "Docker CE",
			Assets: []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			}{
				{Name: "docker-ce-3-29.5.2-1.fc44-44-x86-64.raw", BrowserDownloadURL: "https://github.com/fedora-sysexts/community/releases/download/docker-ce-3-29.5.2-1.fc44-44-x86-64/docker-ce-3-29.5.2-1.fc44-44-x86-64.raw"},
			},
		},
		// Bare name tag (no version info) — should be skipped
		{
			TagName: "tailscale",
			Body:    "Latest tailscale",
			Assets:  nil,
		},
	}

	body, _ := json.Marshal(releases)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	client := &FCOSClient{
		CatalogURL:    srv.URL,
		HTTP:          srv.Client(),
		FedoraVersion: 44,
	}

	entries, err := client.FetchCatalogArch(context.Background(), "amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should get tailscale (fedora 44, x86-64) and docker-ce (fedora 44, x86-64).
	// Tailscale fedora 43 and arm64 should be filtered out.
	if len(entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(entries))
	}

	names := map[string]bool{}
	for _, e := range entries {
		names[e.Name] = true
	}
	if !names["tailscale"] {
		t.Error("missing tailscale entry")
	}
	if !names["docker-ce"] {
		t.Error("missing docker-ce entry")
	}
}

func TestFCOSClient_FetchCatalogArch_Arm64(t *testing.T) {
	releases := []githubRelease{
		{
			TagName: "tailscale-0-1.98.4-1-44-arm64",
			Body:    "Tailscale arm64",
			Assets: []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			}{
				{Name: "tailscale-0-1.98.4-1-44-arm64.raw", BrowserDownloadURL: "https://github.com/fedora-sysexts/community/releases/download/tailscale-0-1.98.4-1-44-arm64/tailscale-0-1.98.4-1-44-arm64.raw"},
			},
		},
		{
			TagName: "tailscale-0-1.98.4-1-44-x86-64",
			Body:    "Tailscale x86",
			Assets: []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			}{
				{Name: "tailscale-0-1.98.4-1-44-x86-64.raw", BrowserDownloadURL: "https://github.com/fedora-sysexts/community/releases/download/tailscale-0-1.98.4-1-44-x86-64/tailscale-0-1.98.4-1-44-x86-64.raw"},
			},
		},
	}

	body, _ := json.Marshal(releases)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	client := &FCOSClient{
		CatalogURL:    srv.URL,
		HTTP:          srv.Client(),
		FedoraVersion: 44,
	}

	entries, err := client.FetchCatalogArch(context.Background(), "arm64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1 (arm64 only)", len(entries))
	}
	if entries[0].Name != "tailscale" {
		t.Errorf("got name %q, want tailscale", entries[0].Name)
	}
}

func TestFCOSClient_FetchCatalogArch_Deduplication(t *testing.T) {
	// Two releases for the same extension (different versions) — only newest kept.
	releases := []githubRelease{
		{
			TagName: "tailscale-0-1.98.4-1-44-x86-64",
			Body:    "Newer",
			Assets: []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			}{
				{Name: "tailscale-0-1.98.4-1-44-x86-64.raw", BrowserDownloadURL: "https://github.com/fedora-sysexts/community/releases/download/tailscale-0-1.98.4-1-44-x86-64/tailscale-0-1.98.4-1-44-x86-64.raw"},
			},
		},
		{
			TagName: "tailscale-0-1.98.3-1-44-x86-64",
			Body:    "Older",
			Assets: []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			}{
				{Name: "tailscale-0-1.98.3-1-44-x86-64.raw", BrowserDownloadURL: "https://github.com/fedora-sysexts/community/releases/download/tailscale-0-1.98.3-1-44-x86-64/tailscale-0-1.98.3-1-44-x86-64.raw"},
			},
		},
	}

	body, _ := json.Marshal(releases)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	client := &FCOSClient{
		CatalogURL:    srv.URL,
		HTTP:          srv.Client(),
		FedoraVersion: 44,
	}

	entries, err := client.FetchCatalogArch(context.Background(), "amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1 (deduped)", len(entries))
	}
	if entries[0].Version != "0-1.98.4-1" {
		t.Errorf("got version %q, want newest 0-1.98.4-1", entries[0].Version)
	}
}

func TestFCOSClient_FetchCatalogArch_CuratedDescription(t *testing.T) {
	releases := []githubRelease{
		{
			TagName: "docker-ce-3-29.5.2-1.fc44-44-x86-64",
			Body:    "Raw GitHub body text",
			Assets: []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			}{
				{Name: "docker-ce-3-29.5.2-1.fc44-44-x86-64.raw", BrowserDownloadURL: "https://github.com/fedora-sysexts/community/releases/download/test/docker-ce.raw"},
			},
		},
	}

	body, _ := json.Marshal(releases)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	client := &FCOSClient{
		CatalogURL:    srv.URL,
		HTTP:          srv.Client(),
		FedoraVersion: 44,
	}

	entries, err := client.FetchCatalogArch(context.Background(), "amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}
	// Should use curated description, not raw body.
	if entries[0].Description == "Raw GitHub body text" {
		t.Error("expected curated description, got raw body text")
	}
	if entries[0].Category != "Container Runtime" {
		t.Errorf("category = %q, want Container Runtime", entries[0].Category)
	}
}

func TestFCOSClient_FetchCatalogArch_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	client := &FCOSClient{
		CatalogURL:    srv.URL,
		HTTP:          srv.Client(),
		FedoraVersion: 44,
	}

	_, err := client.FetchCatalogArch(context.Background(), "amd64")
	if err == nil {
		t.Fatal("expected error for HTTP 403")
	}
}

func TestFCOSClient_FetchCatalogArch_UnsupportedArch(t *testing.T) {
	client := &FCOSClient{FedoraVersion: 44}
	_, err := client.FetchCatalogArch(context.Background(), "mips")
	if err == nil {
		t.Fatal("expected error for unsupported arch")
	}
}
