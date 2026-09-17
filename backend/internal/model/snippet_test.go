package model

import (
	"testing"
	"time"
)

func TestSnippetModel(t *testing.T) {
	now := time.Now().UTC()
	snippet := Snippet{
		ID:        42,
		Slug:      "a8f31c",
		Content:   "fmt.Println(\"hello world\")",
		CreatedAt: now,
	}

	if snippet.ID != 42 {
		t.Errorf("expected ID 42, got %d", snippet.ID)
	}
	if snippet.Slug != "a8f31c" {
		t.Errorf("expected Slug 'a8f31c', got %q", snippet.Slug)
	}
	if snippet.Content != "fmt.Println(\"hello world\")" {
		t.Errorf("expected Content 'fmt.Println(\"hello world\")', got %q", snippet.Content)
	}
	if !snippet.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, snippet.CreatedAt)
	}
}

func TestCreateSnippetInput(t *testing.T) {
	input := CreateSnippetInput{
		Content: "test snippet content",
	}

	if input.Content != "test snippet content" {
		t.Errorf("expected Content 'test snippet content', got %q", input.Content)
	}
}
