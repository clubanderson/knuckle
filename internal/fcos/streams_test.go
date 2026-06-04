package fcos

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseFedoraVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "stable 44", input: "44.20260510.3.1", want: 44},
		{name: "stable 41", input: "41.20250223.3.0", want: 41},
		{name: "fedora 43", input: "43.20250101.1.0", want: 43},
		{name: "fedora 39", input: "39.20231001.0.0", want: 39},
		{name: "bare major", input: "44", want: 44},
		{name: "boundary low", input: "30.0.0.0", want: 30},
		{name: "boundary high", input: "99.0.0.0", want: 99},
		{name: "empty", input: "", wantErr: true},
		{name: "not a number", input: "abc.123", wantErr: true},
		{name: "too low", input: "29.123.0.0", wantErr: true},
		{name: "too high", input: "100.123.0.0", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseFedoraVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFedoraVersion(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFedoraVersion(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func makeStreamJSON(t *testing.T, arch, release string) []byte {
	t.Helper()
	meta := map[string]any{
		"architectures": map[string]any{
			arch: map[string]any{
				"artifacts": map[string]any{
					"metal": map[string]any{
						"release": release,
					},
				},
			},
		},
	}
	b, err := json.Marshal(meta)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestFetchStreamFedoraVersion_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(makeStreamJSON(t, "x86_64", "44.20260510.3.1"))
	}))
	defer srv.Close()

	got, err := FetchStreamFedoraVersionFromURL(context.Background(), srv.URL, srv.Client())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 44 {
		t.Errorf("got %d, want 44", got)
	}
}

func TestFetchStreamFedoraVersion_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := FetchStreamFedoraVersionFromURL(context.Background(), srv.URL, srv.Client())
	if err == nil {
		t.Fatal("expected error for HTTP 404")
	}
}

func TestFetchStreamFedoraVersion_NoX86(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(makeStreamJSON(t, "aarch64", "44.20260510.3.1"))
	}))
	defer srv.Close()

	_, err := FetchStreamFedoraVersionFromURL(context.Background(), srv.URL, srv.Client())
	if err == nil {
		t.Fatal("expected error for missing x86_64")
	}
}

func TestFetchStreamFedoraVersion_NoMetal(t *testing.T) {
	meta := map[string]any{
		"architectures": map[string]any{
			"x86_64": map[string]any{
				"artifacts": map[string]any{
					"qemu": map[string]any{
						"release": "44.20260510.3.1",
					},
				},
			},
		},
	}
	body, _ := json.Marshal(meta)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	_, err := FetchStreamFedoraVersionFromURL(context.Background(), srv.URL, srv.Client())
	if err == nil {
		t.Fatal("expected error for missing metal artifact")
	}
}

func TestFetchStreamFedoraVersion_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}))
	defer srv.Close()

	_, err := FetchStreamFedoraVersionFromURL(context.Background(), srv.URL, srv.Client())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
