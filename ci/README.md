# CI — Helm Chart Release Pipeline

This directory contains a [Dagger](https://dagger.io) module that packages and publishes Helm charts to GHCR (GitHub Container Registry) as OCI artifacts and creates GitHub Releases.

## Overview

Pushing a git tag of the form `charts/<name>/v<semver>` triggers the release pipeline. The pipeline:

1. Validates that the tag version matches the `version` field in the chart's `Chart.yaml`
2. Packages the chart with `helm package`
3. Pushes the OCI artifact to `ghcr.io/<owner>/charts/<name>:<version>`
4. Creates a GitHub Release with auto-generated release notes

Pre-release versions (e.g. `v2.0.1-beta.1`) are published to GHCR and marked as pre-release on GitHub.

## Supported charts

| Tag prefix | Chart directory |
|---|---|
| `charts/netbird/v*` | `charts/netbird/` |
| `charts/kubernetes-operator/v*` | `charts/kubernetes-operator/` |
| `charts/netbird-operator-config/v*` | `charts/netbird-operator-config/` |

## GitHub Actions

The workflow in `.github/workflows/release.yml` triggers automatically on any tag matching `charts/*/v*`. It uses [`dagger/dagger-for-github`](https://github.com/dagger/dagger-for-github), which handles Dagger engine setup — no manual CLI installation or Go setup is needed.

The `GITHUB_TOKEN` provided by GitHub Actions is used for both pushing to GHCR (`packages: write`) and creating the GitHub Release (`contents: write`). No additional secrets are required.

To release a chart, push a tag:

```sh
git tag charts/netbird/v2.1.0
git push origin charts/netbird/v2.1.0
```

The tag version must match the `version` field in `charts/netbird/Chart.yaml` exactly. The pipeline fails early if they differ.

## Local usage

Requires the [Dagger CLI](https://docs.dagger.io/getting-started/installation).

### Dry run

Validates tag format and `Chart.yaml` version, then prints what would be pushed — without touching GHCR or GitHub:

```sh
GITHUB_TOKEN=fake dagger call -m ./ci publish \
  --source=. \
  --chart-tag=charts/netbird/v2.0.0 \
  --registry-token=env://GITHUB_TOKEN \
  --repo-owner=<your-github-username> \
  --dry-run=true
```

Expected output:

```
[DRY-RUN] Would push: ghcr.io/<owner>/charts/netbird:2.0.0 (pre-release: false)
[DRY-RUN] Would create GitHub Release: charts/netbird/v2.0.0
```

### Publish

Requires a real token with `packages:write` and `contents:write` scopes:

```sh
GITHUB_TOKEN=ghp_... dagger call -m ./ci publish \
  --source=. \
  --chart-tag=charts/netbird/v2.0.0 \
  --registry-token=env://GITHUB_TOKEN \
  --repo-owner=<your-github-username>
```

### Discover available functions

```sh
dagger functions -m ./ci
```

## Installing a published chart

```sh
helm install netbird oci://ghcr.io/<owner>/charts/netbird --version 2.0.0
```

Pre-release versions are not shown as "latest" and must be specified explicitly:

```sh
helm install netbird oci://ghcr.io/<owner>/charts/netbird --version 2.0.0-beta.1
```

## Project structure

```
ci/
  main.go         # Dagger module: NetbirdHelm struct with Publish() function
  git.go          # validateChartVersion() — reads Chart.yaml via Dagger Directory API
  helm.go         # helmPublish() — helm package/login/push via alpine/helm:3 container
                  # createGitHubRelease() — GitHub Releases REST API
  logic/
    tags.go       # ParseTag(), IsPreRelease() — pure functions, no Dagger dependency
    chart.go      # CheckChartVersion() — pure YAML validation
    tags_test.go  # Unit tests for ParseTag and IsPreRelease
    chart_test.go # Unit tests for CheckChartVersion
  tests/
    main.go       # Dagger test module: integration tests against the ci/ module
  dagger.json     # Dagger module metadata
  go.mod          # Go module for ci/
```

## Testing

### Unit tests

Pure logic tests in `ci/logic/` run with standard `go test`. No Dagger engine required:

```sh
go test ./ci/logic/...
```

### Integration tests

Integration tests are a [Dagger test module](https://docs.dagger.io/reference/best-practices/modules#test-module) in `ci/tests/`. They call the `ci/` module directly through its public API and require a running Dagger engine (Docker must be running):

```sh
# Run all integration tests
dagger call -m ./ci/tests all --source=.

# Run individual tests
dagger call -m ./ci/tests publish-dry-run --source=.
dagger call -m ./ci/tests publish-dry-run-pre-release --source=.
dagger call -m ./ci/tests invalid-tag-fails --source=.
```

The integration tests cover:
- Dry-run output for a stable release (correct registry path, `pre-release: false`)
- Dry-run output for a pre-release version (`pre-release: true`)
- Invalid tag format returns an error
