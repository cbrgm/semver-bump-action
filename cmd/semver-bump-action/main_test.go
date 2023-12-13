package main

import (
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestBumpSemver(t *testing.T) {
	testCases := []struct {
		name           string
		currentVersion string
		bumpLevel      string
		preReleaseTag  string
		expected       string
		expectError    bool
	}{
		{"Major Bump", "1.0.0", "major", "", "2.0.0", false},
		{"Minor Bump", "1.2.0", "minor", "", "1.3.0", false},
		{"Patch Bump", "1.2.3", "patch", "", "1.2.4", false},
		{"PreMajor Bump", "1.0.0", "premajor", "beta", "2.0.0-beta.0", false},
		{"PreMinor Bump", "1.2.0", "preminor", "alpha", "1.3.0-alpha.0", false},
		{"PrePatch Bump", "1.2.3", "prepatch", "rc", "1.2.4-rc.0", false},
		{"Prerelease Bump", "1.2.3-alpha.0", "prerelease", "alpha", "1.2.3-alpha.1", false},
		{"Invalid Version", "invalid", "major", "", "", true},
		{"Invalid Bump Level", "1.2.3", "invalid", "", "", true},
		{"With 'v' Prefix", "v1.2.3", "patch", "", "v1.2.4", false},
		{"Negative Prerelease Number", "1.2.3-alpha.-1", "prerelease", "alpha", "", true},
		{"Zero Patch Version", "1.2.0", "prepatch", "rc", "1.2.1-rc.0", false},
		// Version with leading zeros
		{"Leading Zero in Minor", "1.05.0", "minor", "", "1.6.0", false},
		{"Leading Zero in Patch", "1.0.09", "patch", "", "1.0.10", false},
		{"Leading Zero in Prerelease", "1.0.0-alpha.01", "prerelease", "alpha", "1.0.0", true},

		// Prerelease to Major/Minor/Patch
		{"Prerelease to Major", "1.0.0-alpha.1", "major", "", "2.0.0", false},
		{"Prerelease to Minor", "1.2.3-beta.2", "minor", "", "1.3.0", false},
		{"Prerelease to Patch", "1.2.3-rc.3", "patch", "", "1.2.3", false},

		// Bumping to specific prerelease versions
		{"Specific Prerelease Tag", "1.2.3", "prerelease", "beta", "1.2.3-beta.0", false},
		{"Increment Specific Prerelease", "1.2.3-beta.0", "prerelease", "beta", "1.2.3-beta.1", false},

		// Error cases
		{"Empty Current Version", "", "major", "", "", true},
		{"Non-numeric Version", "1.x.3", "major", "", "", true},
		{"Unsupported Prerelease", "1.0.0", "premajor", "xyz123", "2.0.0-xyz123.0", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := bumpSemver(tc.currentVersion, tc.bumpLevel, tc.preReleaseTag)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected an error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestBumpStandardVersion(t *testing.T) {
	testCases := []struct {
		name        string
		version     string
		bumpLevel   string
		expected    string
		expectError bool
	}{
		{"Increment Major", "1.2.3", "major", "2.0.0", false},
		{"Increment Minor", "1.2.3", "minor", "1.3.0", false},
		{"Increment Patch", "1.2.3", "patch", "1.2.4", false},
		{"Invalid Bump Level", "1.2.3", "invalid", "", true},
		{"Non-Semantic Version", "1.x.3", "major", "", true},
		// Additional test cases can be added here...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := semver.StrictNewVersion(tc.version)
			if err != nil {
				if !tc.expectError {
					t.Fatalf("unexpected error: %v", err) // Error not expected, fail the test
				}
				// If an error is expected and occurs, no further action needed for this test case
				return
			}

			result, err := bumpStandardVersion(v, tc.bumpLevel)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected an error, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.String() != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestBumpPreReleaseVersion(t *testing.T) {
	testCases := []struct {
		name          string
		version       string
		bumpLevel     string
		preReleaseTag string
		expected      string
		expectError   bool
	}{
		{"Major Version", "1.0.0", "premajor", "alpha", "2.0.0-alpha.0", false},
		{"Minor Version", "1.2.0", "preminor", "beta", "1.3.0-beta.0", false},
		{"Patch Version", "1.2.3", "prepatch", "rc", "1.2.4-rc.0", false},
		{"Existing Prerelease Increment", "1.2.3-alpha.1", "prerelease", "alpha", "1.2.3-alpha.2", false},
		{"New Prerelease", "1.2.3", "prerelease", "alpha", "1.2.3-alpha.0", false},
		{"Invalid Version", "invalid", "premajor", "alpha", "", true},
		{"Invalid Level", "1.2.3", "invalid", "alpha", "", true},
		// Different prerelease tags
		{"Different Prerelease Tag", "1.2.3-beta.1", "prerelease", "rc", "1.2.3-rc.0", false},
		{"No Prerelease Tag Provided", "1.2.3", "prerelease", "", "1.2.3-alpha.0", false},

		// Edge cases
		{"Prerelease with High Number", "1.2.3-alpha.10", "prerelease", "alpha", "1.2.3-alpha.11", false},
		{"Preminor without Tag", "1.2.0", "preminor", "", "1.3.0-alpha.0", false},

		// Invalid versions and levels
		{"Non-Semantic Version", "1.2", "prerelease", "alpha", "", true},
		{"Negative Prerelease Number", "1.2.3-alpha.-1", "prerelease", "alpha", "", true},

		// Specific version scenarios
		{"Zero Patch Version", "1.2.0", "prepatch", "rc", "1.2.1-rc.0", false},
		{"Zero Minor and Patch Version", "1.0.0", "preminor", "beta", "1.1.0-beta.0", false},
		{"Prerelease on Zero Patch", "1.2.0-alpha.0", "prerelease", "alpha", "1.2.0-alpha.1", false},

		// Special cases
		{"Empty Version String", "", "prerelease", "alpha", "", true},
		{"Invalid Bump Level", "1.2.3", "xyz", "alpha", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := semver.StrictNewVersion(tc.version)
			if err != nil {
				if !tc.expectError {
					t.Fatalf("unexpected error: %v", err) // Error not expected, fail the test
				}
				// If an error is expected and occurs, no further action needed for this test case
				return
			}

			// Skip calling bumpPreReleaseVersion if v is nil
			if v == nil {
				t.Fatalf("version is nil for non-error case")
			}

			result, err := bumpPreReleaseVersion(v, tc.bumpLevel, tc.preReleaseTag)
			if err != nil {
				if !tc.expectError {
					t.Fatalf("unexpected error: %v", err)
				}
				// Expected error occurred
				return
			}

			if result.String() != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}
