package attestation

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBindSubject(t *testing.T) {
	digest := strings.Repeat("a", 64)
	for _, tc := range []struct {
		name, payload, digest string
		valid                 bool
	}{
		{"exact archive", `{"version":1,"subject":{"sha256":"` + digest + `"}}`, digest, true},
		{"other archive", `{"version":1,"subject":{"sha256":"` + strings.Repeat("b", 64) + `"}}`, digest, false},
		{"missing subject", `{"version":1}`, digest, false},
		{"future schema", `{"version":2,"subject":{"sha256":"` + digest + `"}}`, digest, false},
		{"short digest", `{"version":1,"subject":{"sha256":"abc"}}`, "abc", false},
		{"malformed", `null`, digest, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := BindSubject(json.RawMessage(tc.payload), tc.digest)
			if (err == nil) != tc.valid {
				t.Fatalf("BindSubject: %v", err)
			}
		})
	}
}
