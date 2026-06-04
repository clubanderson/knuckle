package bakery

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ParseFCOSTagName extracts the extension name, version, Fedora major version,
// and architecture from a fedora-sysexts/community release tag.
//
// Tag format: <name>-<version-parts>-<fedoraVersion>-<arch>
// where arch is "x86-64" (two segments) or "arm64" (one segment).
//
// Examples:
//
//	"tailscale-0-1.98.4-1-44-x86-64"              → ("tailscale", "0-1.98.4-1", "44", "x86-64", nil)
//	"docker-ce-3-29.5.2-1.fc44-44-arm64"          → ("docker-ce", "3-29.5.2-1.fc44", "44", "arm64", nil)
//	"vscodium-1.121.03429-el8-44-x86-64"          → ("vscodium", "1.121.03429-el8", "44", "x86-64", nil)
//	"vscode"                                       → ("", "", "", "", error)
func ParseFCOSTagName(tag string) (name, version, fedoraVersion, arch string, err error) {
	if tag == "" {
		return "", "", "", "", fmt.Errorf("empty tag")
	}

	// Strip architecture suffix. "x86-64" is two dash-separated segments.
	var rest string
	switch {
	case strings.HasSuffix(tag, "-x86-64"):
		arch = "x86-64"
		rest = tag[:len(tag)-len("-x86-64")]
	case strings.HasSuffix(tag, "-arm64"):
		arch = "arm64"
		rest = tag[:len(tag)-len("-arm64")]
	default:
		return "", "", "", "", fmt.Errorf("unrecognized arch suffix in tag %q", tag)
	}

	// The last dash-separated segment of rest is the Fedora major version.
	lastDash := strings.LastIndexByte(rest, '-')
	if lastDash < 0 {
		return "", "", "", "", fmt.Errorf("no fedora version segment in tag %q", tag)
	}
	fedoraVersion = rest[lastDash+1:]
	nameVersion := rest[:lastDash]

	// Validate Fedora version is a two-digit number in expected range.
	fv, parseErr := strconv.Atoi(fedoraVersion)
	if parseErr != nil || fv < 30 || fv > 99 {
		return "", "", "", "", fmt.Errorf("invalid Fedora version %q in tag %q", fedoraVersion, tag)
	}

	// Split name from version using the same heuristic as Flatcar ParseTagName:
	// find the first dash followed by a digit.
	name, version = splitNameVersion(nameVersion)
	if name == "" {
		return "", "", "", "", fmt.Errorf("could not extract name from tag %q", tag)
	}

	return name, version, fedoraVersion, arch, nil
}

// splitNameVersion splits "<name>-<version>" where version starts with a digit.
func splitNameVersion(s string) (name, version string) {
	for i := 1; i < len(s); i++ {
		if s[i-1] == '-' && unicode.IsDigit(rune(s[i])) {
			return s[:i-1], s[i:]
		}
	}
	return "", ""
}
