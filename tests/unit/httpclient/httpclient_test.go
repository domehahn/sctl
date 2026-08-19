package httpclient_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/domehahn/skpm/v2/internal/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientHasTimeout(t *testing.T) {
	c := httpclient.New()
	assert.Equal(t, httpclient.Timeout, c.Timeout, "New() must not return a client with an unbounded (zero) timeout")
	assert.NotZero(t, c.Timeout)
}

func TestNewClientHasRedirectPolicy(t *testing.T) {
	c := httpclient.New()
	require.NotNil(t, c.CheckRedirect, "New() must set a redirect policy rather than relying on net/http's implicit default")
}

func TestClientStopsAfterMaxRedirects(t *testing.T) {
	var redirectCount int
	mux := http.NewServeMux()
	mux.HandleFunc("/loop", func(w http.ResponseWriter, r *http.Request) {
		redirectCount++
		http.Redirect(w, r, "/loop", http.StatusFound)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := httpclient.New()
	_, err := c.Get(srv.URL + "/loop")
	require.Error(t, err, "an infinite-redirect server must not be followed forever")
	assert.LessOrEqual(t, redirectCount, httpclient.MaxRedirects+2,
		"client followed noticeably more than MaxRedirects redirects before giving up")
}

func TestClientRejectsRedirectToNonHTTPScheme(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "file:///etc/passwd")
		w.WriteHeader(http.StatusFound)
	}))
	defer srv.Close()

	c := httpclient.New()
	resp, err := c.Get(srv.URL)
	if err == nil {
		defer resp.Body.Close()
		// Some http.Client versions surface a non-HTTP Location as a
		// completed (non-redirected) response rather than an error from
		// CheckRedirect; either way the client must not have followed it.
		assert.NotEqual(t, http.StatusOK, resp.StatusCode)
		return
	}
	assert.Contains(t, err.Error(), "non-HTTP(S) scheme")
}

func TestClientTimesOutOnSlowServer(t *testing.T) {
	unblock := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-unblock // blocks until explicitly released below
	}))
	defer srv.Close()
	// httptest.Server.Close() blocks until all in-flight handlers return, so
	// the handler above must be released (close(unblock)) before Close()
	// runs — done explicitly right after the client call rather than via a
	// second defer, to avoid depending on defer LIFO ordering between the
	// two cleanups.

	c := httpclient.New()
	c.Timeout = 100 * time.Millisecond // don't actually wait 5 real minutes for this test

	start := time.Now()
	_, err := c.Get(srv.URL)
	elapsed := time.Since(start)
	close(unblock) // release the handler goroutine before srv.Close() runs

	require.Error(t, err, "a request to a server that never responds must eventually fail, not hang forever")
	assert.Less(t, elapsed, 5*time.Second, "client took far longer than its configured Timeout to give up")
}

// ── CopyLimited / ReadAllLimited ────────────────────────────────────────

func TestCopyLimitedAllowsWithinCap(t *testing.T) {
	var buf bytes.Buffer
	n, err := httpclient.CopyLimitedN(&buf, strings.NewReader("hello"), 10)
	require.NoError(t, err)
	assert.Equal(t, int64(5), n)
	assert.Equal(t, "hello", buf.String())
}

func TestCopyLimitedAllowsExactlyAtCap(t *testing.T) {
	var buf bytes.Buffer
	n, err := httpclient.CopyLimitedN(&buf, strings.NewReader("12345"), 5)
	require.NoError(t, err)
	assert.Equal(t, int64(5), n)
}

func TestCopyLimitedRejectsOversized(t *testing.T) {
	var buf bytes.Buffer
	_, err := httpclient.CopyLimitedN(&buf, strings.NewReader("this is definitely too long"), 5)
	require.Error(t, err)
	assert.True(t, httpclient.IsResponseTooLarge(err))
}

func TestReadAllLimitedAllowsWithinCap(t *testing.T) {
	data, err := httpclient.ReadAllLimitedN(strings.NewReader("hello"), 10)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))
}

func TestReadAllLimitedRejectsOversized(t *testing.T) {
	_, err := httpclient.ReadAllLimitedN(strings.NewReader("this is definitely too long"), 5)
	require.Error(t, err)
	assert.True(t, httpclient.IsResponseTooLarge(err))
}

// TestCopyLimitedRejectsOversizedRealResponse exercises the same guard
// against an actual HTTP response body (not just a bare string reader) —
// the shape every real registry Download() call uses it in.
func TestCopyLimitedRejectsOversizedRealResponse(t *testing.T) {
	const oversizedBody = "0123456789"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(len(oversizedBody)))
		_, _ = io.WriteString(w, oversizedBody)
	}))
	defer srv.Close()

	c := httpclient.New()
	resp, err := c.Get(srv.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = httpclient.CopyLimitedN(&buf, resp.Body, 4)
	require.Error(t, err)
	assert.True(t, httpclient.IsResponseTooLarge(err))
	assert.LessOrEqual(t, buf.Len(), 5, "should not have buffered substantially more than the cap before rejecting")
}
