package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/rs/zerolog/log"
)

type OutputFormat string

const (
	OutputText OutputFormat = "text"
	OutputJSON OutputFormat = "json"
)

// CommandResult is the structured envelope for --output json responses.
type CommandResult struct {
	Success  bool        `json:"success"`
	Command  string      `json:"command"`
	Data     interface{} `json:"data,omitempty"`
	Errors   []string    `json:"errors,omitempty"`
	Warnings []string    `json:"warnings,omitempty"`
}

// UserError represents a user-facing error (exit 1).
type UserError struct{ Message string }

func (e *UserError) Error() string { return e.Message }

// InternalError represents an infrastructure error (exit 2).
type InternalError struct {
	Message string
	Cause   error
}

func (e *InternalError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// AdmissionError represents an admission check failure or service rejection (exit 3).
type AdmissionError struct{ Message string }

func (e *AdmissionError) Error() string { return e.Message }

// FrozenLockfileError represents a lockfile drift / inconsistency in frozen mode (exit 4).
type FrozenLockfileError struct{ Message string }

func (e *FrozenLockfileError) Error() string { return e.Message }

// IntegrityError represents a checksum mismatch or corrupt artifact (exit 5).
type IntegrityError struct{ Message string }

func (e *IntegrityError) Error() string { return e.Message }

var secretRegex = regexp.MustCompile(`(?i)(bearer\s+|token[=:]\s*|password[=:]\s*|secret[=:]\s*)([A-Za-z0-9_\-\.]{8,})`)

func redactSecrets(msg string) string {
	return secretRegex.ReplaceAllString(msg, "$1[REDACTED]")
}

// HandleError prints a clean error message and exits with the appropriate code.
// It never prints raw stack traces to the user.
func HandleError(err error, format OutputFormat) {
	if err == nil {
		return
	}

	var userErr *UserError
	var internalErr *InternalError
	var admissionErr *AdmissionError
	var frozenErr *FrozenLockfileError
	var integrityErr *IntegrityError

	cleanMsg := redactSecrets(err.Error())

	switch {
	case errors.As(err, &admissionErr):
		printErrorAndExit(cleanMsg, format, 3)
	case errors.As(err, &frozenErr):
		printErrorAndExit(cleanMsg, format, 4)
	case errors.As(err, &integrityErr):
		printErrorAndExit(cleanMsg, format, 5)
	case isInternalError(err, &internalErr):
		msg := redactSecrets(internalErr.Message)
		printErrorAndExit(msg, format, 2)
	case isUserError(err, &userErr):
		printErrorAndExit(cleanMsg, format, 1)
	default:
		printErrorAndExit(cleanMsg, format, 1)
	}
}

func printErrorAndExit(msg string, format OutputFormat, code int) {
	if format == OutputJSON {
		printJSONError("", msg)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	}
	log.Debug().Int("exit_code", code).Msg("command error exit")
	os.Exit(code)
}

func isUserError(err error, target **UserError) bool {
	if e, ok := err.(*UserError); ok {
		*target = e
		return true
	}
	return false
}

func isInternalError(err error, target **InternalError) bool {
	if e, ok := err.(*InternalError); ok {
		*target = e
		return true
	}
	return false
}

func printJSONError(command, msg string) {
	res := CommandResult{Success: false, Command: command, Errors: []string{msg}}
	data, _ := json.MarshalIndent(res, "", "  ")
	fmt.Fprintln(os.Stdout, string(data))
}

func PrintResult(format OutputFormat, result CommandResult) {
	if format == OutputJSON {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Fprintln(os.Stdout, string(data))
		return
	}
}
