package registry

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

type ConstraintOptions struct {
	IncludePrerelease bool
	AllowDeprecated   bool
	AllowYanked       bool
}

func SelectVersion(versions []VersionInfo, constraint string, opts ConstraintOptions) (string, error) {
	constraint = strings.TrimSpace(constraint)
	if constraint == "" {
		constraint = "latest"
	}
	var candidates []VersionInfo
	for _, v := range versions {
		if !semver.IsValid("v" + v.Version) {
			continue
		}
		if !opts.IncludePrerelease && isPrerelease(v.Version) {
			continue
		}
		if !opts.AllowYanked && v.Yanked {
			continue
		}
		if !matchesConstraint(v.Version, constraint) {
			continue
		}
		candidates = append(candidates, v)
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("no version satisfies constraint %q", constraint)
	}
	sort.Slice(candidates, func(i, j int) bool {
		return semver.Compare("v"+candidates[i].Version, "v"+candidates[j].Version) > 0
	})
	if !opts.AllowDeprecated {
		for _, v := range candidates {
			if !v.Deprecated {
				return v.Version, nil
			}
		}
	}
	return candidates[0].Version, nil
}

func matchesConstraint(version, constraint string) bool {
	switch {
	case constraint == "" || constraint == "latest":
		return true
	case strings.HasPrefix(constraint, "^"):
		return matchesCaret(version, strings.TrimPrefix(constraint, "^"))
	case strings.HasPrefix(constraint, "~"):
		return matchesTilde(version, strings.TrimPrefix(constraint, "~"))
	case strings.Contains(constraint, " "):
		return matchesRange(version, constraint)
	default:
		return normalizeVersion(version) == normalizeVersion(strings.TrimPrefix(constraint, "="))
	}
}

func matchesCaret(version, base string) bool {
	major, minor, patch, ok := parseStable(base)
	if !ok {
		return false
	}
	vMajor, vMinor, vPatch, ok := parseStable(version)
	if !ok {
		return false
	}
	if compareTriplet(vMajor, vMinor, vPatch, major, minor, patch) < 0 {
		return false
	}
	if major > 0 {
		return vMajor == major
	}
	if minor > 0 {
		return vMajor == 0 && vMinor == minor
	}
	return vMajor == 0 && vMinor == 0 && vPatch == patch
}

func matchesTilde(version, base string) bool {
	major, minor, patch, ok := parseStable(base)
	if !ok {
		return false
	}
	vMajor, vMinor, vPatch, ok := parseStable(version)
	if !ok {
		return false
	}
	return vMajor == major && vMinor == minor && compareTriplet(vMajor, vMinor, vPatch, major, minor, patch) >= 0
}

func matchesRange(version, constraint string) bool {
	for _, part := range strings.Fields(constraint) {
		if !matchesComparator(version, part) {
			return false
		}
	}
	return true
}

func matchesComparator(version, comparator string) bool {
	ops := []string{">=", "<=", ">", "<", "="}
	for _, op := range ops {
		if strings.HasPrefix(comparator, op) {
			target := strings.TrimSpace(strings.TrimPrefix(comparator, op))
			cmp := semver.Compare("v"+normalizeVersion(version), "v"+normalizeVersion(target))
			switch op {
			case ">=":
				return cmp >= 0
			case "<=":
				return cmp <= 0
			case ">":
				return cmp > 0
			case "<":
				return cmp < 0
			case "=":
				return cmp == 0
			}
		}
	}
	return normalizeVersion(version) == normalizeVersion(comparator)
}

func isPrerelease(version string) bool {
	return semver.Prerelease("v"+version) != ""
}

func normalizeVersion(version string) string {
	return strings.TrimPrefix(strings.TrimSpace(version), "v")
}

func parseStable(version string) (int, int, int, bool) {
	version = normalizeVersion(version)
	if semver.Prerelease("v"+version) != "" || !semver.IsValid("v"+version) {
		return 0, 0, 0, false
	}
	var major, minor, patch int
	if _, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch); err != nil {
		return 0, 0, 0, false
	}
	return major, minor, patch, true
}

func compareTriplet(aMajor, aMinor, aPatch, bMajor, bMinor, bPatch int) int {
	switch {
	case aMajor != bMajor:
		return aMajor - bMajor
	case aMinor != bMinor:
		return aMinor - bMinor
	default:
		return aPatch - bPatch
	}
}
