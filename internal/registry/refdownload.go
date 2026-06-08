package registry

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DownloadRef downloads a skill directory at a specific git ref (branch, commit, tag)
// from a registry that supports archive downloads. The skill files are extracted
// to destDir. skillSubPath is the path of the skill within the repo (e.g. "skills/my-skill"
// or "" if the skill is at the repo root).
func DownloadRef(ctx context.Context, reg Registry, skillName, ref, skillSubPath, destDir string) error {
	switch r := reg.(type) {
	case *GitHubRegistry:
		return downloadGitHubRef(ctx, r, skillName, ref, skillSubPath, destDir)
	case *GitLabRegistry:
		return downloadGitLabRef(ctx, r, skillName, ref, skillSubPath, destDir)
	default:
		return fmt.Errorf("--ref is only supported for github and gitlab registries")
	}
}

// downloadGitHubRef downloads the repo archive at ref and extracts the skill subdirectory.
// GitHub API: GET /repos/{owner}/{repo}/zipball/{ref}
func downloadGitHubRef(ctx context.Context, r *GitHubRegistry, skillName, ref, skillSubPath, destDir string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/%s", r.org, r.repo, ref)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if r.token != "" {
		req.Header.Set("Authorization", "Bearer "+r.token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("github: download ref %s: %w", ref, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github: download ref %s: HTTP %d", ref, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("github: read archive: %w", err)
	}

	return extractSkillFromZIP(bytes.NewReader(data), int64(len(data)), skillName, skillSubPath, destDir)
}

// downloadGitLabRef uses GitLab's repository archive API with path filtering.
// GitLab API: GET /projects/{id}/repository/archive?sha={ref}&format=zip&path={subPath}
func downloadGitLabRef(ctx context.Context, r *GitLabRegistry, skillName, ref, skillSubPath, destDir string) error {
	encodedProject := strings.ReplaceAll(r.projectID, "/", "%2F")
	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/archive?sha=%s&format=zip",
		r.baseURL, encodedProject, ref)
	if skillSubPath != "" {
		url += "&path=" + skillSubPath
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if r.token != "" {
		req.Header.Set("PRIVATE-TOKEN", r.token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("gitlab: download ref %s: %w", ref, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gitlab: download ref %s: HTTP %d", ref, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("gitlab: read archive: %w", err)
	}

	return extractSkillFromZIP(bytes.NewReader(data), int64(len(data)), skillName, skillSubPath, destDir)
}

// extractSkillFromZIP finds the skill directory inside a ZIP archive and extracts it to destDir.
// For GitHub archives the root is a random directory like "owner-repo-abc123/".
// For GitLab archives the root is "projectname-ref-sha/".
// skillSubPath is the relative path to the skill within the repo (may be "").
func extractSkillFromZIP(r io.ReaderAt, size int64, skillName, skillSubPath, destDir string) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}

	// Find the archive root prefix (first directory component).
	var archiveRoot string
	for _, f := range zr.File {
		parts := strings.SplitN(f.Name, "/", 2)
		if len(parts) >= 1 && parts[0] != "" {
			archiveRoot = parts[0] + "/"
			break
		}
	}

	// Build the prefix to strip: archiveRoot + skillSubPath
	var skillPrefix string
	if skillSubPath != "" {
		skillPrefix = archiveRoot + strings.Trim(skillSubPath, "/") + "/"
	} else {
		// No subpath: look for a directory matching the skill name inside the archive root.
		for _, f := range zr.File {
			if strings.HasPrefix(f.Name, archiveRoot) {
				rest := strings.TrimPrefix(f.Name, archiveRoot)
				if strings.HasPrefix(rest, skillName+"/") {
					skillPrefix = archiveRoot + skillName + "/"
					break
				}
			}
		}
		if skillPrefix == "" {
			// Fall back: treat the whole archive root as the skill.
			skillPrefix = archiveRoot
		}
	}

	found := false
	for _, f := range zr.File {
		if !strings.HasPrefix(f.Name, skillPrefix) {
			continue
		}
		found = true
		rel := strings.TrimPrefix(f.Name, skillPrefix)
		if rel == "" {
			continue
		}

		outPath := filepath.Join(destDir, filepath.FromSlash(rel))
		if f.FileInfo().IsDir() {
			os.MkdirAll(outPath, 0o755)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		if err := extractZipFile(f, outPath); err != nil {
			return err
		}
	}

	if !found {
		return fmt.Errorf("skill %q not found in archive (tried prefix %q)", skillName, skillPrefix)
	}
	return nil
}

func extractZipFile(f *zip.File, dest string) error {
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}
