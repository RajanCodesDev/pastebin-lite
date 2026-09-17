package validation

import (
	"errors"
	"strings"
	"testing"

	"github.com/pastebin-lite/backend/internal/model"
)

func TestValidateCreateSnippetInput(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr error
	}{
		{
			name:    "valid standard content",
			content: "package main\n\nfunc main() {}",
			wantErr: nil,
		},
		{
			name:    "valid content with leading/trailing whitespace",
			content: "  hello world  \n",
			wantErr: nil,
		},
		{
			name:    "empty content",
			content: "",
			wantErr: ErrEmptyContent,
		},
		{
			name:    "whitespace-only spaces",
			content: "     ",
			wantErr: ErrEmptyContent,
		},
		{
			name:    "whitespace-only tabs and newlines",
			content: "\n\t \r\n \t",
			wantErr: ErrEmptyContent,
		},
		{
			name:    "content exactly at maximum length",
			content: strings.Repeat("x", MaxSnippetLength),
			wantErr: nil,
		},
		{
			name:    "content exceeding maximum length by one rune",
			content: strings.Repeat("x", MaxSnippetLength+1),
			wantErr: ErrContentTooLong,
		},
		{
			name:    "multibyte unicode characters within limit",
			content: strings.Repeat("🚀", 100),
			wantErr: nil,
		},
		{
			name:    "multibyte unicode characters exceeding limit",
			content: strings.Repeat("🚀", MaxSnippetLength+1),
			wantErr: ErrContentTooLong,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := model.CreateSnippetInput{
				Content: tc.content,
			}
			err := ValidateCreateSnippetInput(input)

			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
		})
	}
}
