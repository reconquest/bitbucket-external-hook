package main

import (
	"strings"
)

type HookSettings struct {
	Exe      string `json:"exe"`
	SafePath bool   `json:"safe_path"`
	Params   string `json:"params"`
}

type ResponseHooks struct {
	Values []Hook
}

type Hook struct {
	Details struct {
		Key     string `json:"key"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"details"`
	Configured bool `json:"configured"`
	Enabled    bool `json:"enabled"`
	Scope struct {
		Type string `json:"type"`
	} `json:"scope"`
}

type ResponseError struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (err ResponseError) String() string {
	var messages []string

	for _, nestedError := range err.Errors {
		messages = append(messages, nestedError.Message)
	}

	return strings.Join(messages, "\n")
}
