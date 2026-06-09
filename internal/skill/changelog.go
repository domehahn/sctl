package skill

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var changelogEntryPattern = regexp.MustCompile(`(?m)^##\s+v?(\S+)`)

// ReadVersion reads the VERSION file from dir and returns the trimmed version string.
func ReadVersion(dir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(dir, "VERSION"))
	if err != nil {
		return "", err
	}
	v := strings.TrimSpace(string(data))
	v = strings.TrimPrefix(v, "v")
	if v == "" {
		return "", fmt.Errorf("VERSION file is empty")
	}
	return v, nil
}

// EnsureChangelogEntry checks whether CHANGELOG.md contains an entry for version.
// If not, it prepends a placeholder entry and writes the file back.
// Returns true if a new entry was added, false if the entry already existed.
func EnsureChangelogEntry(dir, version string) (bool, error) {
	path := filepath.Join(dir, "CHANGELOG.md")

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("read CHANGELOG.md: %w", err)
	}

	// Check if entry already exists.
	for _, m := range changelogEntryPattern.FindAllStringSubmatch(string(existing), -1) {
		if m[1] == version || m[1] == "v"+version {
			return false, nil
		}
	}

	newEntry := fmt.Sprintf("## %s\n\n- \n\n", version)

	var content string
	if len(existing) == 0 {
		content = "# Changelog\n\n" + newEntry
	} else {
		// Insert after the top-level # heading if present, otherwise prepend.
		text := string(existing)
		if idx := strings.Index(text, "\n"); idx >= 0 && strings.HasPrefix(text, "# ") {
			content = text[:idx+1] + "\n" + newEntry + strings.TrimLeft(text[idx+1:], "\n")
		} else {
			content = newEntry + text
		}
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return false, fmt.Errorf("write CHANGELOG.md: %w", err)
	}
	return true, nil
}
