// The conformance skill's real, deterministic runtime: digest a supplied message.
// It has no tools. skil, rather than this process, evaluates and enforces its contract.
package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	var exchange struct {
		Version int `json:"version"`
		Request struct {
			Test struct {
				Input struct {
					Message string `json:"message"`
				} `json:"input"`
			} `json:"test"`
		} `json:"request"`
	}
	if err := json.NewDecoder(os.Stdin).Decode(&exchange); err != nil || exchange.Version != 1 {
		os.Exit(1)
	}
	sum := sha256.Sum256([]byte(exchange.Request.Test.Input.Message))
	if err := json.NewEncoder(os.Stdout).Encode(map[string]any{"type": "final", "final": map[string]any{"outputs": []string{fmt.Sprintf("%x", sum)}}}); err != nil {
		os.Exit(1)
	}
}
