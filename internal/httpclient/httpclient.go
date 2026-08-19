// Package httpclient is the single place skpm constructs *http.Client and
// bounds response bodies for every HTTP call to a registry, git host, or
// any other URL a lockfile/config points at — none of which skpm should
// treat as fully trusted. Before this package existed, every registry
// backend built its own bare &http.Client{} (or used http.DefaultClient
// directly) with no request timeout, and read response bodies with
// io.Copy/io.ReadAll with no size cap: a slow or hostile server could hang
// a request indefinitely or exhaust memory/disk on an unbounded response.
package httpclient

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Timeout bounds an entire HTTP round trip (connect + redirects + headers +
// body). Generous enough for a large artifact download over a slow link,
// bounded enough that a hung/slow-loris server can't stall a command
// forever.
const Timeout = 5 * time.Minute

// MaxRedirects caps how many redirects a single request will follow.
// Matches net/http's own built-in default (it errors past 10 redirects
// even with no CheckRedirect set) — set explicitly here so the behavior is
// documented and doesn't silently change if that default ever does.
const MaxRedirects = 10

// MaxDownloadBytes bounds how much response body skpm will read for a
// single artifact/archive download. 2GiB is far larger than any legitimate
// skill package; it exists to stop a malicious or misconfigured server
// (or a redirect to one) from exhausting disk/memory via an effectively
// unbounded response.
const MaxDownloadBytes = 2 * 1024 * 1024 * 1024

// ErrResponseTooLarge is returned by CopyLimited/ReadAllLimited when a
// response body exceeds MaxDownloadBytes.
var ErrResponseTooLarge = fmt.Errorf("httpclient: response body exceeds the %d byte download limit", MaxDownloadBytes)

// New returns an *http.Client configured with skpm's default timeout and
// redirect policy. Every registry backend and the installer's downloader
// should use a client from here rather than http.DefaultClient or a bare
// &http.Client{}.
func New() *http.Client {
	return &http.Client{
		Timeout:       Timeout,
		CheckRedirect: checkRedirect,
	}
}

func checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= MaxRedirects {
		return fmt.Errorf("httpclient: stopped after %d redirects", MaxRedirects)
	}
	if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
		return fmt.Errorf("httpclient: refusing redirect to non-HTTP(S) scheme %q", req.URL.Scheme)
	}
	return nil
}

// CopyLimited copies src to dst, stopping with ErrResponseTooLarge if more
// than MaxDownloadBytes would be written. Use in place of a bare io.Copy
// when src is an HTTP response body from a download.
func CopyLimited(dst io.Writer, src io.Reader) (int64, error) {
	return CopyLimitedN(dst, src, MaxDownloadBytes)
}

// CopyLimitedN is CopyLimited with an explicit cap, for a call site that
// knows a resource should be much smaller than MaxDownloadBytes (or, in
// tests, needs a cap small enough to exercise the "too large" path without
// actually transferring gigabytes of data).
func CopyLimitedN(dst io.Writer, src io.Reader, max int64) (int64, error) {
	limited := io.LimitReader(src, max+1)
	n, err := io.Copy(dst, limited)
	if err != nil {
		return n, err
	}
	if n > max {
		return n, ErrResponseTooLarge
	}
	return n, nil
}

// ReadAllLimited is io.ReadAll with the same MaxDownloadBytes cap as
// CopyLimited. Use in place of a bare io.ReadAll when reading an HTTP
// response body fully into memory (e.g. a source archive to extract from).
func ReadAllLimited(src io.Reader) ([]byte, error) {
	return ReadAllLimitedN(src, MaxDownloadBytes)
}

// ReadAllLimitedN is ReadAllLimited with an explicit cap; see CopyLimitedN.
func ReadAllLimitedN(src io.Reader, max int64) ([]byte, error) {
	limited := io.LimitReader(src, max+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return data, err
	}
	if int64(len(data)) > max {
		return nil, ErrResponseTooLarge
	}
	return data, nil
}

// IsResponseTooLarge reports whether err is (or wraps) ErrResponseTooLarge.
func IsResponseTooLarge(err error) bool {
	return errors.Is(err, ErrResponseTooLarge)
}
