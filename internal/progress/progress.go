package progress

import (
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

// isCI returns true when running in a CI environment where progress bars
// should be suppressed.
func isCI() bool {
	for _, v := range []string{"CI", "GITHUB_ACTIONS", "GITLAB_CI", "NO_COLOR"} {
		if os.Getenv(v) != "" {
			return true
		}
	}
	return false
}

// NewDownloadBar returns a progressbar for download progress.
// In CI environments it returns a silent bar that does not write to stdout.
func NewDownloadBar(total int64, description string) *progressbar.ProgressBar {
	opts := []progressbar.Option{
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(65),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {}),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	}
	if isCI() {
		opts = append(opts, progressbar.OptionSetVisibility(false))
	}
	return progressbar.NewOptions64(total, opts...)
}

// NewReader wraps r with a progress bar that advances as bytes are read.
// If total is -1, the bar is displayed without a known total.
func NewReader(r io.Reader, total int64, description string) io.Reader {
	bar := NewDownloadBar(total, description)
	pr := progressbar.NewReader(r, bar)
	return &pr
}
