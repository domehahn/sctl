package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/domehahn/skpm/v2/internal/attestation"
	"github.com/domehahn/skpm/v2/internal/config"
	"github.com/domehahn/skpm/v2/internal/registry"
	"github.com/spf13/cobra"
)

// skilAttestationSubject is the subset of a skil Attestation
// (https://github.com/domehahn/skil, `skil attest --output`) skpm reads
// to auto-fill --digest and --type when they aren't passed explicitly.
// skpm treats the rest of the file as an opaque predicate — it does not
// depend on skil's full schema.
type skilAttestationSubject struct {
	Subject struct {
		SHA256 string `json:"sha256"`
	} `json:"subject"`
}

// defaultSkilAttestationType is "scan" — a skil scan/eval attestation is
// evidence that the artifact was scanned, which is the closest fit among
// the fixed set some registries accept. SkillForge, for example, only
// accepts type in {signature, scan, provenance, sbom} — there is no
// free-form predicate-type URI convention here, unlike DSSE/in-toto.
// --type lets a user override this for a registry with different rules.
const defaultSkilAttestationType = "scan"

func newAttestCmd() *cobra.Command {
	var (
		source        string
		file          string
		predicateType string
		digest        string
	)

	cmd := &cobra.Command{
		Use:   "attest <skill>@<version>",
		Short: "Attach a third-party attestation (e.g. a skil scan/eval result) to a published skill version",
		Long: `Uploads an attestation — evidence produced outside skpm, most commonly the
output of "skil attest --output attestation.json" — and attaches it to an
already-published skill version as first-class registry metadata,
independent of the artifact bytes.

The predicate is stored unchanged. Scan evidence must use attestation version 1
and bind subject.sha256 to the exact published package digest; --digest cannot
override a different signed subject. Attest the final archive with skil. Requires a registry that implements attestation storage (see
"skpm registry capabilities" — SkillForge registries support this).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, version, err := parseSkillAtVersion(args[0])
			if err != nil {
				return err
			}
			if file == "" {
				return &UserError{Message: "--file is required (path to an attestation JSON document)"}
			}

			payload, err := os.ReadFile(file)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("read attestation file: %v", err)}
			}
			if !json.Valid(payload) {
				return &UserError{Message: fmt.Sprintf("%s is not valid JSON", file)}
			}

			if digest == "" {
				var subj skilAttestationSubject
				if err := json.Unmarshal(payload, &subj); err == nil && subj.Subject.SHA256 != "" {
					digest = subj.Subject.SHA256
				}
			}
			if digest == "" {
				return &UserError{Message: "--digest is required (could not infer subject.sha256 from the attestation file — pass it explicitly)"}
			}
			if predicateType == "" {
				predicateType = defaultSkilAttestationType
			}
			if predicateType == "scan" {
				if err := attestation.BindSubject(payload, digest); err != nil {
					return &UserError{Message: err.Error()}
				}
			}

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			src := sourceOrDefault(source, cfg)
			if src == "" {
				return &UserError{Message: "no registry specified and no default_registry configured"}
			}
			reg, err := registry.New(src, cfg)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("registry: %v", err)}
			}
			att, ok := reg.(registry.AttestationRegistry)
			if !ok {
				return &UserError{Message: fmt.Sprintf("registry %q does not support attestations", src)}
			}

			ref := registry.SkillVersionRef{Namespace: "default", Name: name, Version: version}

			if globalDryRun {
				fmt.Fprintf(cmd.OutOrStdout(), "Dry run: would attest %s@%s with %s (type=%s, digest=%s)\n", name, version, file, predicateType, digest)
				return nil
			}

			rec, err := att.Attest(cmd.Context(), ref, registry.AttestationRequest{
				Type:      predicateType,
				Digest:    digest,
				Predicate: json.RawMessage(payload),
			})
			if err != nil {
				return &InternalError{Message: "attest", Cause: err}
			}

			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: true, Command: "attest", Data: rec})
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Attested %s@%s\n", name, version)
			fmt.Fprintf(cmd.OutOrStdout(), "  Type:   %s\n", rec.Type)
			fmt.Fprintf(cmd.OutOrStdout(), "  Digest: %s\n", rec.Digest)
			if rec.ID != 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  ID:     %d\n", rec.ID)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "Registry source (uses default_registry if not set)")
	cmd.Flags().StringVar(&file, "file", "", "Path to the attestation JSON document (required)")
	cmd.Flags().StringVar(&predicateType, "type", "", "Predicate type (default: "+defaultSkilAttestationType+")")
	cmd.Flags().StringVar(&digest, "digest", "", "Subject sha256 (default: read from the attestation file's subject.sha256)")
	return cmd
}

// attestationWithVerification augments a stored AttestationRecord with the
// outcome of independently checking its signature (see internal/attestation
// and config.Config.TrustedSigners), when --verify was passed.
type attestationWithVerification struct {
	registry.AttestationRecord
	Verified      *bool  `json:"verified,omitempty"`
	VerifiedKeyID string `json:"verified_key_id,omitempty"`
	VerifyError   string `json:"verify_error,omitempty"`
}

func newAttestationsCmd() *cobra.Command {
	var source string
	var verify bool

	cmd := &cobra.Command{
		Use:   "attestations <skill>@<version>",
		Short: "List attestations attached to a published skill version",
		Long: `Lists attestations attached to a published skill version.

With --verify, each attestation's Ed25519 signature (e.g. one produced by
"skil attest --signing-key") is independently checked against
config.trusted_signers — a map from key_id to base64 Ed25519 public key —
without skpm depending on skil as a library. An attestation with no
"signature" field, or one signed by a key not in trusted_signers, is
reported as not verified and causes a nonzero exit. Missing evidence also
fails verification. Listing without --verify remains an inspection operation.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, version, err := parseSkillAtVersion(args[0])
			if err != nil {
				return err
			}

			cfg, err := config.Load()
			if err != nil {
				return &InternalError{Message: "load config", Cause: err}
			}
			src := sourceOrDefault(source, cfg)
			if src == "" {
				return &UserError{Message: "no registry specified and no default_registry configured"}
			}
			reg, err := registry.New(src, cfg)
			if err != nil {
				return &UserError{Message: fmt.Sprintf("registry: %v", err)}
			}
			att, ok := reg.(registry.AttestationRegistry)
			if !ok {
				return &UserError{Message: fmt.Sprintf("registry %q does not support attestations", src)}
			}

			ref := registry.SkillVersionRef{Namespace: "default", Name: name, Version: version}
			records, err := att.ListAttestations(cmd.Context(), ref)
			if err != nil {
				return &InternalError{Message: "list attestations", Cause: err}
			}

			if verify && len(cfg.TrustedSigners) == 0 {
				return &UserError{Message: "--verify requires at least one entry under trusted_signers in the skpm config"}
			}

			if verify && len(records) == 0 {
				return &UserError{Message: "ATTESTATION_MISSING: no evidence to verify"}
			}
			var verificationError error
			results := make([]attestationWithVerification, len(records))
			for i, rec := range records {
				results[i] = attestationWithVerification{AttestationRecord: rec}
				if !verify {
					continue
				}
				ok := false
				sig, err := attestation.Verify(rec.Predicate, cfg.TrustedSigners)
				if err == nil && rec.Type == "scan" {
					err = attestation.BindSubject(rec.Predicate, rec.Digest)
				}
				if err != nil {
					results[i].Verified = &ok
					results[i].VerifyError = err.Error()
					verificationError = &UserError{Message: "ATTESTATION_VERIFICATION_FAILED: " + err.Error()}
					continue
				}
				ok = true
				results[i].Verified = &ok
				results[i].VerifiedKeyID = sig.KeyID
			}

			if outputFormat() == OutputJSON {
				PrintResult(OutputJSON, CommandResult{Success: verificationError == nil, Command: "attestations", Data: results, Errors: verificationErrors(verificationError)})
				return verificationError
			}
			if len(results) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No attestations for %s@%s\n", name, version)
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Attestations for %s@%s:\n", name, version)
			for _, r := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "  - type=%s digest=%s created_by=%s created_at=%s",
					r.Type, r.Digest, r.CreatedBy, r.CreatedAt)
				if r.Verified != nil {
					if *r.Verified {
						fmt.Fprintf(cmd.OutOrStdout(), " signature=verified(key=%s)", r.VerifiedKeyID)
					} else {
						fmt.Fprintf(cmd.OutOrStdout(), " signature=NOT-VERIFIED(%s)", r.VerifyError)
					}
				}
				fmt.Fprintln(cmd.OutOrStdout())
			}
			return verificationError
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "Registry source (uses default_registry if not set)")
	cmd.Flags().BoolVar(&verify, "verify", false, "Independently verify each attestation's Ed25519 signature against config.trusted_signers")
	return cmd
}

func verificationErrors(err error) []string {
	if err == nil {
		return nil
	}
	return []string{err.Error()}
}
