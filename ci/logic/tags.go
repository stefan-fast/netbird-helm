// Package logic contains pure functions with no Dagger dependencies,
// making them easily unit-testable with plain go test.
package logic

import (
	"fmt"
	"strings"
)

// ParseTag parses a GHA ref or short-form tag into chart name and version.
// Accepts:
//   - full form:  refs/tags/charts/<name>/v<semver>
//   - short form: charts/<name>/v<semver>
//
// Returns the chart name and version without the leading "v".
func ParseTag(ref string) (chartName, version string, err error) {
	// Strip optional refs/tags/ prefix.
	trimmed := strings.TrimPrefix(ref, "refs/tags/")

	parts := strings.Split(trimmed, "/")
	if len(parts) != 3 {
		return "", "", fmt.Errorf("invalid tag %q: expected 3 slash-separated parts (charts/<name>/v<ver>), got %d", ref, len(parts))
	}

	if parts[0] != "charts" {
		return "", "", fmt.Errorf("invalid tag %q: first segment must be \"charts\", got %q", ref, parts[0])
	}

	name := parts[1]
	if name == "" {
		return "", "", fmt.Errorf("invalid tag %q: chart name is empty", ref)
	}

	ver := parts[2]
	if !strings.HasPrefix(ver, "v") {
		return "", "", fmt.Errorf("invalid tag %q: version segment %q must start with \"v\"", ref, ver)
	}

	ver = strings.TrimPrefix(ver, "v")
	if ver == "" {
		return "", "", fmt.Errorf("invalid tag %q: version is empty after stripping \"v\"", ref)
	}

	return name, ver, nil
}

// IsPreRelease returns true if version contains a pre-release segment (i.e. contains "-").
func IsPreRelease(version string) bool {
	return strings.Contains(version, "-")
}
