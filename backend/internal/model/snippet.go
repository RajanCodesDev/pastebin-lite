package model

import "time"

// Snippet represents the core domain entity for a code/text snippet.
type Snippet struct {
	ID        int64     `json:"id"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateSnippetInput contains the data needed to create a new Snippet.
type CreateSnippetInput struct {
	Content string `json:"content"`
}
