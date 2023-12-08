package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/alexflint/go-arg"
)

// ActionInputs struct holds the required environment variables.
type ActionInputs struct {
	CurrentVersion string `arg:"--current-version,env:CURRENT_VERSION,required"` // The current semantic version
	BumpLevel      string `arg:"--bump-level,env:BUMP_LEVEL,required"`           // The level to bump the version (major, minor, patch, etc.)
	PreReleaseTag  string `arg:"--prerelease-tag,env:PRERELEASE_TAG"`            // Optional tag for prerelease versions
}

func main() {
	var args ActionInputs
	arg.MustParse(&args)

	// Bumping the semantic version
	newVersion, err := bumpSemver(args.CurrentVersion, args.BumpLevel, args.PreReleaseTag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error bumping semver: %v\n", err)
		os.Exit(1)
	}

	setActionOutput("new_version", newVersion)
}

func bumpSemver(currentVersion, bumpLevel, preReleaseTag string) (string, error) {
	// Handling version with 'v' prefix
	hasVPrefix := strings.HasPrefix(currentVersion, "v")
	if hasVPrefix {
		currentVersion = strings.TrimPrefix(currentVersion, "v")
	}

	// Parsing semantic version
	v, err := semver.NewVersion(currentVersion)
	if err != nil {
		return "", fmt.Errorf("invalid semver '%s': %w", currentVersion, err)
	}

	// Bumping version based on the level
	var newVersion semver.Version
	switch bumpLevel {
	case "major", "minor", "patch":
		newVersion, err = bumpStandardVersion(v, bumpLevel)
	case "premajor", "preminor", "prepatch", "prerelease":
		newVersion, err = bumpPreReleaseVersion(v, bumpLevel, preReleaseTag)
	default:
		return "", fmt.Errorf("unsupported bump level '%s'", bumpLevel)
	}

	if err != nil {
		return "", err
	}

	// Re-adding 'v' prefix if it was originally present
	result := newVersion.String()
	if hasVPrefix {
		result = "v" + result
	}
	return result, nil
}

func bumpStandardVersion(v *semver.Version, level string) (semver.Version, error) {
	switch level {
	case "major":
		return v.IncMajor(), nil
	case "minor":
		return v.IncMinor(), nil
	case "patch":
		return v.IncPatch(), nil
	}
	return semver.Version{}, errors.New("invalid standard version bump level")
}

func bumpPreReleaseVersion(v *semver.Version, level, preReleaseTag string) (semver.Version, error) {
	var newVersion semver.Version
	var err error

	switch level {
	case "premajor":
		newVersion = v.IncMajor()
	case "preminor":
		newVersion = v.IncMinor()
	case "prepatch":
		newVersion = v.IncPatch()
	case "prerelease":
		newVersion, err = incrementPrereleaseVersion(v, preReleaseTag)
		if err != nil {
			return semver.Version{}, err
		}
		return newVersion, nil
	default:
		return semver.Version{}, errors.New("invalid prerelease version bump level")
	}

	// Setting new prerelease version
	if preReleaseTag == "" {
		preReleaseTag = "alpha"
	}
	newVersion, err = newVersion.SetPrerelease(preReleaseTag + ".0")
	if err != nil {
		return semver.Version{}, err
	}

	return newVersion, nil
}

func incrementPrereleaseVersion(v *semver.Version, preReleaseTag string) (semver.Version, error) {
	prerelease := v.Prerelease()
	if preReleaseTag == "" {
		preReleaseTag = "alpha"
	}

	// If there's an existing prerelease version, increment it
	if strings.HasPrefix(prerelease, preReleaseTag) {
		parts := strings.SplitN(prerelease, ".", 2)
		if len(parts) == 2 {
			number, err := strconv.Atoi(parts[1])
			if err != nil {
				// Handle conversion error, which includes negative numbers
				return *v, fmt.Errorf("invalid prerelease format: %w", err)
			}
			if number < 0 {
				// Explicitly handle negative numbers
				return *v, fmt.Errorf("negative number in prerelease is not valid")
			}
			number++
			return v.SetPrerelease(fmt.Sprintf("%s.%d", preReleaseTag, number))
		}
	}

	// Start with .0 for the given prerelease tag
	return v.SetPrerelease(preReleaseTag + ".0")
}
