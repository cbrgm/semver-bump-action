# Semver Bump Action

**Automatically bump the given semver version up.**

[![GitHub release](https://img.shields.io/github/release/cbrgm/semver-bump-action.svg)](https://github.com/cbrgm/semver-bump-action)
[![Go Report Card](https://goreportcard.com/badge/github.com/cbrgm/semver-bump-action)](https://goreportcard.com/report/github.com/cbrgm/semver-bump-action)
[![go-lint-test](https://github.com/cbrgm/semver-bump-action/actions/workflows/go-lint-test.yml/badge.svg)](https://github.com/cbrgm/semver-bump-action/actions/workflows/go-lint-test.yml)
[![go-binaries](https://github.com/cbrgm/semver-bump-action/actions/workflows/go-binaries.yml/badge.svg)](https://github.com/cbrgm/semver-bump-action/actions/workflows/go-binaries.yml)
[![container](https://github.com/cbrgm/semver-bump-action/actions/workflows/container.yml/badge.svg)](https://github.com/cbrgm/semver-bump-action/actions/workflows/container.yml)

## Inputs

### `current-version`
**Required** - The semantic version (semver) that needs to be bumped. For example, `1.2.3`.

### `bump-level`
**Required** - Specifies the semver bump level. Allowed values:
- `major` - Increments the major version (e.g., 1.2.3 to 2.0.0).
- `minor` - Increments the minor version (e.g., 1.2.3 to 1.3.0).
- `patch` - Increments the patch version (e.g., 1.2.3 to 1.2.4).
- `premajor` - Creates a premajor prerelease (e.g., 1.2.3 to 2.0.0-alpha.0).
- `preminor` - Creates a preminor prerelease (e.g., 1.2.3 to 1.3.0-alpha.0).
- `prepatch` - Creates a prepatch prerelease (e.g., 1.2.3 to 1.2.4-alpha.0).
- `prerelease` - Increments an existing prerelease or creates a new one (e.g., 1.2.3-alpha.0 to 1.2.3-alpha.1).

### `prerelease-tag`
**Required** - The tag to use for prereleases (e.g., `alpha`, `beta`). For example, specifying `alpha` in combination with `bump-level` as `prerelease` will result in versions like `1.2.3-alpha.0`.

## Outputs

### `new_version`
The bumped semantic version. For example, if `current-version` is `1.2.3` and `bump-level` is `minor`, `new_version` will be `1.3.0`.

## Workflow Usage

Add the following step to your GitHub Actions Workflow:

```yaml
name: Demo Workflow

on:
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Get Latest Tag
        id: current
        run: |
          latest_tag="v1.0.0"
          echo "latest_tag=$latest_tag" >> $GITHUB_ENV

      - name: Bump Minor Version
        id: bump
        uses: cbrgm/semver-bump-action@main
        with:
          current-version: ${{ env.latest_tag }}
          bump-level: minor

      - name: Output
        run: |
          new_tag=${{ steps.bump.outputs.new_version }}
          echo $new_tag
          echo $new_tag
```

Example: Bump + Publish Minor Version Tag

```yaml
name: Publish Tag

on:
  workflow_dispatch:
    inputs:
      bump-level:
        required: true
        type: choice
        description: 'The semver bump level'
        options:
          - 'major'
          - 'minor'
          - 'patch'
          - 'premajor'
          - 'preminor'
          - 'prepatch'
          - 'prerelease'
        default: 'patch'

      prerelease-tag:
        required: false
        description: 'The tag to use for prereleases'
        default: 'alpha'

jobs:
  publish-tag:
    runs-on: ubuntu-latest
    steps:

      - name: Checkout
        uses: actions/checkout@v4
        with:
          ssh-key: "${{ secrets.COMMIT_KEY }}"

      - name: Get Latest Tag
        id: current
        run: |
          git fetch --tags
          latest_tag=$(git tag --sort=taggerdate | tail -n 1)
          echo "current version is: $latest_tag"
          echo "latest_tag=$latest_tag" >> $GITHUB_ENV

      - name: Bump Version
        id: bump
        uses: cbrgm/semver-bump-action@main
        with:
          current-version: ${{ env.latest_tag }}
          bump-level: ${{ github.event.inputs.bump-level }}
          prerelease-tag: ${{ github.event.inputs.prerelease-tag }}

      - name: Publish Tag
        run: |
          git fetch --tags
          latest_tag=$(git tag --sort=taggerdate | tail -n 1)
          new_tag=${{ steps.bump.outputs.new_version }}
          if [[ $(git rev-list $latest_tag..HEAD --count) -gt 0 ]]; then
            git config user.name "GitHub Actions"
            git config user.email "github-actions@users.noreply.github.com"
            git tag $new_tag
            git push origin $new_tag
          else
            echo "No new commits since last tag. Skipping tag push."
          fi
```

### Local Development

You can build this action from source using `Go`:

```bash
make build
```

## Contributing & License

We welcome and value your contributions to this project! 👍 If you're interested in making improvements or adding features, please refer to our [Contributing Guide](https://github.com/cbrgm/semver-bump-action/blob/main/CONTRIBUTING.md). This guide provides comprehensive instructions on how to submit changes, set up your development environment, and more.

Please note that this project is developed in my spare time and is available for free 🕒💻. As an open-source initiative, it is governed by the [Apache 2.0 License](https://github.com/cbrgm/semver-bump-action/blob/main/LICENSE). This license outlines your rights and obligations when using, modifying, and distributing this software.

Your involvement, whether it's through code contributions, suggestions, or feedback, is crucial for the ongoing improvement and success of this project. Together, we can ensure it remains a useful and well-maintained resource for everyone 🌍.
