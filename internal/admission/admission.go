package admission

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Decision string

const (
	DecisionAllow    Decision = "ALLOW"
	DecisionDeny     Decision = "DENY"
	DecisionReview   Decision = "REVIEW"
	DecisionReassess Decision = "REASSESS"
)

type AdmissionRequest struct {
	Name           string            `json:"name"`
	Version        string            `json:"version"`
	PackageDigest  string            `json:"package_digest,omitempty"`
	ArtifactDigest string            `json:"artifact_digest,omitempty"`
	Source         string            `json:"source,omitempty"`
	Registry       string            `json:"registry,omitempty"`
	Action         string            `json:"action"` // "publish" or "install"
	Environment    string            `json:"environment,omitempty"`
	LockfileDigest string            `json:"lockfile_digest,omitempty"`
	Attestations   []string          `json:"attestations,omitempty"`
	Timestamp      string            `json:"timestamp,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type AdmissionDecision struct {
	ID            string   `json:"id,omitempty"`
	Decision      Decision `json:"decision"`
	Reason        string   `json:"reason,omitempty"`
	IssuedAt      string   `json:"issued_at,omitempty"`
	ExpiresAt     string   `json:"expires_at,omitempty"`
	PackageDigest string   `json:"package_digest,omitempty"`
	Signature     string   `json:"signature,omitempty"`
}

type Client interface {
	Evaluate(ctx context.Context, req AdmissionRequest) (*AdmissionDecision, error)
}

type HTTPClient struct {
	URL        string
	Enforce    bool
	HTTPClient *http.Client
}

func NewClientFromEnv() *HTTPClient {
	url := os.Getenv("SKPM_ADMISSION_URL")
	enforce := strings.EqualFold(os.Getenv("SKPM_ADMISSION_ENFORCE"), "true") || os.Getenv("SKPM_ADMISSION_ENFORCE") == "1"
	if url == "" {
		return nil
	}
	return &HTTPClient{
		URL:     url,
		Enforce: enforce,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HTTPClient) Evaluate(ctx context.Context, req AdmissionRequest) (*AdmissionDecision, error) {
	if c == nil || c.URL == "" {
		return &AdmissionDecision{Decision: DecisionAllow}, nil
	}

	if req.Timestamp == "" {
		req.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal admission request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create admission request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	hc := c.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}

	resp, err := hc.Do(httpReq)
	if err != nil {
		if c.Enforce {
			return nil, fmt.Errorf("admission service unavailable (enforced): %w", err)
		}
		return &AdmissionDecision{Decision: DecisionAllow, Reason: fmt.Sprintf("admission unavailable: %v", err)}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		if c.Enforce {
			return nil, fmt.Errorf("admission service error (HTTP %d): %s", resp.StatusCode, string(respBody))
		}
		return &AdmissionDecision{Decision: DecisionAllow, Reason: fmt.Sprintf("admission service returned HTTP %d", resp.StatusCode)}, nil
	}

	var dec AdmissionDecision
	if err := json.NewDecoder(resp.Body).Decode(&dec); err != nil {
		if c.Enforce {
			return nil, fmt.Errorf("admission decision decode: %w", err)
		}
		return &AdmissionDecision{Decision: DecisionAllow}, nil
	}

	// Validate decision expiration / freshness
	if dec.ExpiresAt != "" {
		expTime, err := time.Parse(time.RFC3339, dec.ExpiresAt)
		if err == nil && time.Now().After(expTime) {
			if c.Enforce {
				return &dec, fmt.Errorf("admission decision expired at %s", dec.ExpiresAt)
			}
		}
	}

	// Validate digest binding if returned in decision
	if dec.PackageDigest != "" && req.PackageDigest != "" && dec.PackageDigest != req.PackageDigest {
		if c.Enforce {
			return &dec, fmt.Errorf("admission decision digest mismatch: expected %s, got %s", req.PackageDigest, dec.PackageDigest)
		}
	}

	if dec.Decision == DecisionDeny || dec.Decision == DecisionReview || dec.Decision == DecisionReassess {
		if c.Enforce {
			return &dec, fmt.Errorf("admission check rejected action %s for %s@%s: %s (%s)", req.Action, req.Name, req.Version, dec.Decision, dec.Reason)
		}
	}

	return &dec, nil
}
