package admission_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/domehahn/skpm/v2/internal/admission"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdmissionAllow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(admission.AdmissionDecision{
			Decision: admission.DecisionAllow,
			Reason:   "passed checks",
		})
	}))
	defer ts.Close()

	client := &admission.HTTPClient{
		URL:     ts.URL,
		Enforce: true,
	}

	dec, err := client.Evaluate(context.Background(), admission.AdmissionRequest{
		Name:    "my-skill",
		Version: "1.0.0",
		Action:  "install",
	})
	require.NoError(t, err)
	assert.Equal(t, admission.DecisionAllow, dec.Decision)
}

func TestAdmissionDenyEnforced(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(admission.AdmissionDecision{
			Decision: admission.DecisionDeny,
			Reason:   "policy restriction",
		})
	}))
	defer ts.Close()

	client := &admission.HTTPClient{
		URL:     ts.URL,
		Enforce: true,
	}

	dec, err := client.Evaluate(context.Background(), admission.AdmissionRequest{
		Name:    "my-skill",
		Version: "1.0.0",
		Action:  "install",
	})
	require.Error(t, err)
	assert.Equal(t, admission.DecisionDeny, dec.Decision)
	assert.Contains(t, err.Error(), "admission check rejected")
}

func TestAdmissionUnavailableEnforcedFailsClosed(t *testing.T) {
	client := &admission.HTTPClient{
		URL:     "http://127.0.0.1:1", // closed port
		Enforce: true,
	}

	_, err := client.Evaluate(context.Background(), admission.AdmissionRequest{
		Name:    "my-skill",
		Version: "1.0.0",
		Action:  "install",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "admission service unavailable (enforced)")
}
