package main

import (
	"context"
	"fmt"

	"ci/internal/dagger"
	"ci/logic"
)

// validateChartVersion reads charts/<chartName>/Chart.yaml from the Dagger
// source directory and verifies its version field matches the supplied version.
func validateChartVersion(ctx context.Context, source *dagger.Directory, chartName, version string) error {
	path := fmt.Sprintf("charts/%s/Chart.yaml", chartName)

	data, err := source.File(path).Contents(ctx)
	if err != nil {
		return fmt.Errorf("reading Chart.yaml at %q: %w", path, err)
	}

	return logic.CheckChartVersion([]byte(data), version, path)
}
