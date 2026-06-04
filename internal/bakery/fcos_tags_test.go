package bakery

import "testing"

func TestParseFCOSTagName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		tag         string
		wantName    string
		wantVersion string
		wantFedora  string
		wantArch    string
		wantErr     bool
	}{
		{
			tag:         "tailscale-0-1.98.4-1-44-x86-64",
			wantName:    "tailscale",
			wantVersion: "0-1.98.4-1",
			wantFedora:  "44",
			wantArch:    "x86-64",
		},
		{
			tag:         "tailscale-0-1.98.4-1-44-arm64",
			wantName:    "tailscale",
			wantVersion: "0-1.98.4-1",
			wantFedora:  "44",
			wantArch:    "arm64",
		},
		{
			tag:         "docker-ce-3-29.5.2-1.fc44-44-x86-64",
			wantName:    "docker-ce",
			wantVersion: "3-29.5.2-1.fc44",
			wantFedora:  "44",
			wantArch:    "x86-64",
		},
		{
			tag:         "docker-ce-3-29.5.2-1.fc44-44-arm64",
			wantName:    "docker-ce",
			wantVersion: "3-29.5.2-1.fc44",
			wantFedora:  "44",
			wantArch:    "arm64",
		},
		{
			tag:         "vscodium-1.121.03429-el8-44-x86-64",
			wantName:    "vscodium",
			wantVersion: "1.121.03429-el8",
			wantFedora:  "44",
			wantArch:    "x86-64",
		},
		{
			tag:         "vscode-1.122.1-1780040915.el8-43-arm64",
			wantName:    "vscode",
			wantVersion: "1.122.1-1780040915.el8",
			wantFedora:  "43",
			wantArch:    "arm64",
		},
		{
			tag:         "virtctl-1.5.0-44-x86-64",
			wantName:    "virtctl",
			wantVersion: "1.5.0",
			wantFedora:  "44",
			wantArch:    "x86-64",
		},
		{
			tag:         "google-chrome-149.0.7827.53-1-44-x86-64",
			wantName:    "google-chrome",
			wantVersion: "149.0.7827.53-1",
			wantFedora:  "44",
			wantArch:    "x86-64",
		},
		{
			tag:         "openconnect-9.12.git.270.549fd2d-0.fc44-44-x86-64",
			wantName:    "openconnect",
			wantVersion: "9.12.git.270.549fd2d-0.fc44",
			wantFedora:  "44",
			wantArch:    "x86-64",
		},
		{
			tag:         "cilium-cli-0.19.4-44-x86-64",
			wantName:    "cilium-cli",
			wantVersion: "0.19.4",
			wantFedora:  "44",
			wantArch:    "x86-64",
		},
		{
			tag:         "cloud-hypervisor-51.0.0-39.41-43-arm64",
			wantName:    "cloud-hypervisor",
			wantVersion: "51.0.0-39.41",
			wantFedora:  "43",
			wantArch:    "arm64",
		},
		{
			tag:         "nordvpn-gui-5.0.0-1-44-x86-64",
			wantName:    "nordvpn-gui",
			wantVersion: "5.0.0-1",
			wantFedora:  "44",
			wantArch:    "x86-64",
		},
		{
			tag:         "netbird-ui-0.71.4-1-43-x86-64",
			wantName:    "netbird-ui",
			wantVersion: "0.71.4-1",
			wantFedora:  "43",
			wantArch:    "x86-64",
		},
		{
			tag:         "1password-gui-8.12.22-1-44-x86-64",
			wantName:    "1password-gui",
			wantVersion: "8.12.22-1",
			wantFedora:  "44",
			wantArch:    "x86-64",
		},
		{
			tag:         "bitwarden-2026.5.0-1-44-x86-64",
			wantName:    "bitwarden",
			wantVersion: "2026.5.0-1",
			wantFedora:  "44",
			wantArch:    "x86-64",
		},
		{
			tag:         "microsoft-edge-148.0.3967.96-1-43-x86-64",
			wantName:    "microsoft-edge",
			wantVersion: "148.0.3967.96-1",
			wantFedora:  "43",
			wantArch:    "x86-64",
		},
		{
			tag:         "littlesnitch-1.0.9-1-44-arm64",
			wantName:    "littlesnitch",
			wantVersion: "1.0.9-1",
			wantFedora:  "44",
			wantArch:    "arm64",
		},
		{
			tag:         "glab-1.101.0-1-44-arm64",
			wantName:    "glab",
			wantVersion: "1.101.0-1",
			wantFedora:  "44",
			wantArch:    "arm64",
		},
		// Error cases
		{tag: "", wantErr: true},
		{tag: "tailscale", wantErr: true},
		{tag: "vscode", wantErr: true},
		{tag: "latest", wantErr: true},
		{tag: "no-arch-suffix-here", wantErr: true},
		{tag: "bad-fedver-xx-x86-64", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			t.Parallel()
			name, version, fedora, arch, err := ParseFCOSTagName(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFCOSTagName(%q) error = %v, wantErr %v", tt.tag, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if name != tt.wantName {
				t.Errorf("name = %q, want %q", name, tt.wantName)
			}
			if version != tt.wantVersion {
				t.Errorf("version = %q, want %q", version, tt.wantVersion)
			}
			if fedora != tt.wantFedora {
				t.Errorf("fedora = %q, want %q", fedora, tt.wantFedora)
			}
			if arch != tt.wantArch {
				t.Errorf("arch = %q, want %q", arch, tt.wantArch)
			}
		})
	}
}

func TestSplitNameVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input       string
		wantName    string
		wantVersion string
	}{
		{"tailscale-0-1.98.4-1", "tailscale", "0-1.98.4-1"},
		{"docker-ce-3-29.5.2-1.fc44", "docker-ce", "3-29.5.2-1.fc44"},
		{"vscodium-1.121.03429-el8", "vscodium", "1.121.03429-el8"},
		{"cilium-cli-0.19.4", "cilium-cli", "0.19.4"},
		{"virtctl-1.5.0", "virtctl", "1.5.0"},
		{"google-chrome-149.0.7827.53-1", "google-chrome", "149.0.7827.53-1"},
		{"no-version-here", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			name, ver := splitNameVersion(tt.input)
			if name != tt.wantName {
				t.Errorf("name = %q, want %q", name, tt.wantName)
			}
			if ver != tt.wantVersion {
				t.Errorf("version = %q, want %q", ver, tt.wantVersion)
			}
		})
	}
}
