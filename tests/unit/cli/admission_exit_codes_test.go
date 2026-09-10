package cli_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/domehahn/skpm/v2/internal/admission"
	"github.com/domehahn/skpm/v2/internal/cli"
	"github.com/domehahn/skpm/v2/internal/lockfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdmissionIntegrationInInstall(t *testing.T) {
	// Create mock admission server returning ALLOW
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		var req admission.AdmissionRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "install", req.Action)

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(admission.AdmissionDecision{
			Decision: admission.DecisionAllow,
			Reason:   "passed policy",
		})
	}))
	defer ts.Close()

	t.Setenv("SKPM_ADMISSION_URL", ts.URL)
	t.Setenv("SKPM_ADMISSION_ENFORCE", "true")

	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "test-skill",
		Version: "1.0.0",
		Source:  "local",
		SHA256:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	})
	require.NoError(t, lf.Write(filepath.Join(dir, "agent-skills.lock")))

	cmd := cli.NewRootCmd()
	cmd.SetArgs([]string{"install", "--dry-run"})
	require.NoError(t, cmd.Execute())
}

func TestAdmissionIntegrationDeniedEnforced(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(admission.AdmissionDecision{
			Decision: admission.DecisionDeny,
			Reason:   "forbidden skill",
		})
	}))
	defer ts.Close()

	t.Setenv("SKPM_ADMISSION_URL", ts.URL)
	t.Setenv("SKPM_ADMISSION_ENFORCE", "true")

	dir := t.TempDir()
	t.Chdir(dir)

	lf := lockfile.New()
	lf.Upsert(lockfile.SkillLock{
		Name:    "bad-skill",
		Version: "1.0.0",
		Source:  "local",
		SHA256:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	})
	require.NoError(t, lf.Write(filepath.Join(dir, "agent-skills.lock")))

	cmd := cli.NewRootCmd()
	cmd.SetArgs([]string{"install", "--dry-run"})
	err := cmd.Execute()
	require.Error(t, err)
	assert.IsType(t, &cli.AdmissionError{}, err)
	assert.Contains(t, err.Error(), "admission check rejected")
}

func TestErrorTypesMapping(t *testing.T) {
	usrErr := &cli.UserError{Message: "user issue"}
	assert.Equal(t, "user issue", usrErr.Error())

	intErr := &cli.InternalError{Message: "internal issue"}
	assert.Equal(t, "internal issue", intErr.Error())

	admErr := &cli.AdmissionError{Message: "admission issue"}
	assert.Equal(t, "admission issue", admErr.Error())

	frzErr := &cli.FrozenLockfileError{Message: "frozen issue"}
	assert.Equal(t, "frozen issue", frzErr.Error())

	ingErr := &cli.IntegrityError{Message: "integrity issue"}
	assert.Equal(t, "integrity issue", ingErr.Error())
}
