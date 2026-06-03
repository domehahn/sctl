package cli

import (
	"encoding/json"
	"fmt"
	"os"

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

// HandleError prints a clean error message and exits with the appropriate code.
// It never prints raw stack traces to the user.
func HandleError(err error, format OutputFormat) {
	if err == nil {
		return
	}

	var userErr *UserError
	var internalErr *InternalError

	switch {
	case isUserError(err, &userErr):
		if format == OutputJSON {
			printJSONError("", err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
		log.Debug().Err(err).Msg("user error")
		os.Exit(1)
	case isInternalError(err, &internalErr):
		if format == OutputJSON {
			printJSONError("", err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", internalErr.Message)
		}
		log.Debug().Err(internalErr.Cause).Msg("internal error")
		os.Exit(2)
	default:
		if format == OutputJSON {
			printJSONError("", err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
		log.Debug().Err(err).Msg("error")
		os.Exit(1)
	}
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
