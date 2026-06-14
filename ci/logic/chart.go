package logic

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// ChartYAML holds the fields we care about from Chart.yaml.
type ChartYAML struct {
	Version string `yaml:"version"`
}

// CheckChartVersion parses Chart.yaml content and verifies the version field
// matches the supplied version. This is the pure, testable core of validation.
func CheckChartVersion(data []byte, version, path string) error {
	var chart ChartYAML
	if err := yaml.Unmarshal(data, &chart); err != nil {
		return fmt.Errorf("parsing Chart.yaml at %q: %w", path, err)
	}

	if chart.Version != version {
		return fmt.Errorf("Chart.yaml version %q does not match expected version %q in %q", chart.Version, version, path)
	}

	return nil
}
