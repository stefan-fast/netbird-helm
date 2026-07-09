// Tests provides integration tests for the netbird-helm Dagger module.
// Run all tests with: dagger call -m ./ci/tests all --source=.
package main

import (
	"context"
	"errors"
	"strings"

	"dagger/tests/internal/dagger"
)

type Tests struct{}

// All runs all tests. Pass the repository root as source.
func (m *Tests) All(ctx context.Context, source *dagger.Directory) error {
	if err := m.PublishDryRun(ctx, source); err != nil {
		return err
	}
	if err := m.PublishDryRunPreRelease(ctx, source); err != nil {
		return err
	}
	if err := m.InvalidTagFails(ctx, source); err != nil {
		return err
	}
	return nil
}

// PublishDryRun verifies that a dry-run publish returns the expected output.
func (m *Tests) PublishDryRun(ctx context.Context, source *dagger.Directory) error {
	fakeToken := dag.SetSecret("registry-token", "fake-token")

	// Build a synthetic source with Chart.yaml version matching the tag to avoid
	// coupling the test to the real chart version in the repository.
	chartYAML := "apiVersion: v2\nname: netbird\nversion: 2.0.0\n"
	syntheticSource := source.WithNewFile("charts/netbird/Chart.yaml", chartYAML)

	result, err := dag.NetbirdHelm().Publish(
		ctx,
		syntheticSource,
		"charts/netbird/v2.0.0",
		fakeToken,
		"testorg",
		dagger.NetbirdHelmPublishOpts{DryRun: true},
	)
	if err != nil {
		return err
	}

	if !strings.Contains(result, "[DRY-RUN]") {
		return errors.New("expected [DRY-RUN] in output, got: " + result)
	}
	if !strings.Contains(result, "ghcr.io/testorg/charts/netbird:2.0.0") {
		return errors.New("expected registry path in output, got: " + result)
	}
	if !strings.Contains(result, "pre-release: false") {
		return errors.New("expected 'pre-release: false' in output, got: " + result)
	}
	return nil
}

// PublishDryRunPreRelease verifies that a pre-release version is flagged correctly.
// It uses a synthetic source directory with Chart.yaml matching the pre-release version.
func (m *Tests) PublishDryRunPreRelease(ctx context.Context, source *dagger.Directory) error {
	fakeToken := dag.SetSecret("registry-token", "fake-token")

	// Build a synthetic source with Chart.yaml version matching the pre-release tag.
	chartYAML := "apiVersion: v2\nname: netbird\nversion: 2.0.0-beta.1\n"
	syntheticSource := source.WithNewFile("charts/netbird/Chart.yaml", chartYAML)

	result, err := dag.NetbirdHelm().Publish(
		ctx,
		syntheticSource,
		"charts/netbird/v2.0.0-beta.1",
		fakeToken,
		"testorg",
		dagger.NetbirdHelmPublishOpts{DryRun: true},
	)
	if err != nil {
		return err
	}

	if !strings.Contains(result, "pre-release: true") {
		return errors.New("expected 'pre-release: true' in output, got: " + result)
	}
	return nil
}

// InvalidTagFails verifies that an invalid chart tag returns an error.
func (m *Tests) InvalidTagFails(ctx context.Context, source *dagger.Directory) error {
	fakeToken := dag.SetSecret("registry-token", "fake-token")

	_, err := dag.NetbirdHelm().Publish(
		ctx,
		source,
		"invalid-tag-format",
		fakeToken,
		"testorg",
		dagger.NetbirdHelmPublishOpts{DryRun: true},
	)
	if err == nil {
		return errors.New("expected error for invalid tag, got nil")
	}
	return nil
}
