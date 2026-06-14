package main

import (
	"context"
	"fmt"
	"strings"

	"ci/internal/dagger"
	"ci/logic"
)

// helmPublish packages a Helm chart and pushes it to GHCR OCI using the Dagger helm module.
func helmPublish(ctx context.Context, source *dagger.Directory, owner, chartName string, token *dagger.Secret) error {
	ownerLower := strings.ToLower(owner)
	registry := fmt.Sprintf("oci://ghcr.io/%s/charts", ownerLower)

	return dag.Helm().
		WithRegistryAuth("ghcr.io", ownerLower, token).
		Chart(source.Directory(fmt.Sprintf("charts/%s", chartName))).
		Package().
		Publish(ctx, registry)
}

// createGitHubRelease creates a GitHub Release using the Dagger gh module.
func createGitHubRelease(ctx context.Context, owner, repo, chartName, version string, token *dagger.Secret) error {
	tag := fmt.Sprintf("charts/%s/v%s", chartName, version)

	return dag.Gh().Release().Create(ctx, tag, tag, dagger.GhReleaseCreateOpts{
		Repo:          owner + "/" + repo,
		Token:         token,
		GenerateNotes: true,
		PreRelease:    logic.IsPreRelease(version),
	})
}
