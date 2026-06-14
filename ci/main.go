package main

import (
	"context"
	"fmt"
	"strings"

	"ci/internal/dagger"
	"ci/logic"
)

// NetbirdHelm provides Dagger functions for releasing Helm charts to GHCR.
type NetbirdHelm struct{}

// Publish packages a Helm chart and pushes it to GHCR OCI, then creates a
// GitHub Release. Set dryRun=true to validate and log what would happen
// without actually pushing.
//
// chartTag format: charts/<name>/v<semver>, e.g. charts/netbird/v2.0.1
func (m NetbirdHelm) Publish(
	ctx context.Context,
	// Source directory of the repository root
	source *dagger.Directory,
	// Chart tag in format charts/<name>/v<semver>, e.g. charts/netbird/v2.0.1
	chartTag string,
	// GitHub token with packages:write and contents:write permissions
	registryToken *dagger.Secret,
	// GitHub repository owner (org or user), e.g. myorg
	repoOwner string,
	// If true, log what would be done without pushing or creating a release
	// +optional
	dryRun bool,
) (string, error) {
	chartName, version, err := logic.ParseTag(chartTag)
	if err != nil {
		return "", fmt.Errorf("invalid chart tag: %w", err)
	}

	if err := validateChartVersion(ctx, source, chartName, version); err != nil {
		return "", fmt.Errorf("chart version validation failed: %w", err)
	}

	ownerLower := strings.ToLower(repoOwner)
	pre := logic.IsPreRelease(version)

	if dryRun {
		return fmt.Sprintf(
			"[DRY-RUN] Would push: ghcr.io/%s/charts/%s:%s (pre-release: %v)\n[DRY-RUN] Would create GitHub Release: charts/%s/v%s",
			ownerLower, chartName, version, pre, chartName, version,
		), nil
	}

	if err := helmPublish(ctx, source, repoOwner, chartName, registryToken); err != nil {
		return "", err
	}

	if err := createGitHubRelease(ctx, repoOwner, "netbird-helm", chartName, version, registryToken); err != nil {
		return "", err
	}

	return fmt.Sprintf("Published ghcr.io/%s/charts/%s:%s and created GitHub Release charts/%s/v%s",
		ownerLower, chartName, version, chartName, version), nil
}
