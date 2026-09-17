package validation

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/pastebin-lite/backend/internal/model"
)

// MaxSnippetLength defines the maximum allowed length of a snippet in characters/runes.
// 65,536 runes (64 KiB at 1 byte/rune) is selected as a sensible limit for a pastebin service:
// it comfortably accommodates lengthy source code files, stack traces, and configurations,
// while preventing denial-of-service via unbounded memory allocation.
const MaxSnippetLength = 65536

// Sentinel domain validation errors.
var (
	ErrEmptyContent    = errors.New("snippet content cannot be empty or whitespace-only")
	ErrContentTooLong  = errors.New("snippet content exceeds maximum allowed length")
)

// ValidateCreateSnippetInput validates input for creating a snippet.
func ValidateCreateSnippetInput(input model.CreateSnippetInput) error {
	trimmed := strings.TrimSpace(input.Content)
	if trimmed == "" {
		return ErrEmptyContent
	}

	if utf8.RuneCountInString(input.Content) > MaxSnippetLength {
		return ErrContentTooLong
	}

	return nil
}
