package logic

import "testing"

func TestParseTag(t *testing.T) {
	tests := []struct {
		name        string
		ref         string
		wantChart   string
		wantVersion string
		wantErr     bool
	}{
		{
			name:        "full ref with semver",
			ref:         "refs/tags/charts/netbird/v2.0.1",
			wantChart:   "netbird",
			wantVersion: "2.0.1",
		},
		{
			name:        "short form",
			ref:         "charts/kubernetes-operator/v0.3.2",
			wantChart:   "kubernetes-operator",
			wantVersion: "0.3.2",
		},
		{
			name:        "full ref with pre-release version",
			ref:         "refs/tags/charts/netbird/v2.0.1-beta.1",
			wantChart:   "netbird",
			wantVersion: "2.0.1-beta.1",
		},
		{name: "non-charts prefix", ref: "refs/tags/helm-v1.9.0", wantErr: true},
		{name: "missing version segment", ref: "refs/tags/charts/netbird", wantErr: true},
		{name: "version missing v prefix", ref: "refs/tags/charts/netbird/2.0.1", wantErr: true},
		{name: "empty string", ref: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chart, version, err := ParseTag(tt.ref)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseTag(%q) expected error, got nil (chart=%q, version=%q)", tt.ref, chart, version)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseTag(%q) unexpected error: %v", tt.ref, err)
			}
			if chart != tt.wantChart {
				t.Errorf("ParseTag(%q) chart = %q, want %q", tt.ref, chart, tt.wantChart)
			}
			if version != tt.wantVersion {
				t.Errorf("ParseTag(%q) version = %q, want %q", tt.ref, version, tt.wantVersion)
			}
		})
	}
}

func TestIsPreRelease(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"2.0.1", false},
		{"2.0.1-beta.1", true},
		{"2.0.1-rc.1", true},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got := IsPreRelease(tt.version)
			if got != tt.want {
				t.Errorf("IsPreRelease(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}
