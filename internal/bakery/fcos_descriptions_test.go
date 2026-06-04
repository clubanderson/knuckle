package bakery

import "testing"

func TestFCOSLookup_KnownExtensions(t *testing.T) {
	t.Parallel()

	required := []string{
		"docker-ce", "tailscale", "vscode",
		"cilium-cli", "1password-cli", "vscodium",
	}

	for _, name := range required {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			meta, ok := FCOSLookup(name)
			if !ok {
				t.Fatalf("FCOSLookup(%q) not found — required by acceptance criteria", name)
			}
			if meta.Short == "" {
				t.Errorf("FCOSLookup(%q).Short is empty", name)
			}
			if meta.Category == "" {
				t.Errorf("FCOSLookup(%q).Category is empty", name)
			}
			if meta.SupportTier == "" {
				t.Errorf("FCOSLookup(%q).SupportTier is empty", name)
			}
		})
	}
}

func TestFCOSLookup_Unknown(t *testing.T) {
	_, ok := FCOSLookup("definitely-not-in-catalog")
	if ok {
		t.Error("FCOSLookup should return false for unknown extensions")
	}
}

func TestFCOSCatalog_AllHaveRequiredFields(t *testing.T) {
	t.Parallel()
	for name, meta := range fcosCatalog {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if meta.Short == "" {
				t.Errorf("fcosCatalog[%q].Short is empty", name)
			}
			if meta.Category == "" {
				t.Errorf("fcosCatalog[%q].Category is empty", name)
			}
			if meta.SupportTier == "" {
				t.Errorf("fcosCatalog[%q].SupportTier is empty", name)
			}
			if meta.Long == "" {
				t.Errorf("fcosCatalog[%q].Long is empty", name)
			}
		})
	}
}
