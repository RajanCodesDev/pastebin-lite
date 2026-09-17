package slug

import (
	"strings"
	"sync"
	"testing"
)

func TestGenerate_DefaultLength(t *testing.T) {
	slug, err := Generate()
	if err != nil {
		t.Fatalf("unexpected error generating slug: %v", err)
	}

	if len(slug) != DefaultSlugLength {
		t.Errorf("expected slug length %d, got %d (slug: %q)", DefaultSlugLength, len(slug), slug)
	}
}

func TestGenerate_AllowedCharacters(t *testing.T) {
	for i := 0; i < 50; i++ {
		slug, err := Generate()
		if err != nil {
			t.Fatalf("unexpected error generating slug: %v", err)
		}

		if slug == "" {
			t.Fatal("generator returned an empty string")
		}

		for _, ch := range slug {
			if !strings.ContainsRune(Charset, ch) {
				t.Fatalf("slug %q contains invalid character %q; allowed charset is %s", slug, ch, Charset)
			}
		}
	}
}

func TestGenerateCustom(t *testing.T) {
	lengths := []int{1, 4, 6, 8, 12, 32}
	for _, l := range lengths {
		slug, err := GenerateCustom(l)
		if err != nil {
			t.Fatalf("unexpected error generating slug of length %d: %v", l, err)
		}
		if len(slug) != l {
			t.Errorf("expected length %d, got %d", l, len(slug))
		}
	}

	// Invalid lengths
	invalidLengths := []int{0, -1, -10}
	for _, l := range invalidLengths {
		_, err := GenerateCustom(l)
		if err == nil {
			t.Errorf("expected error for non-positive length %d, got nil", l)
		}
	}
}

func TestRandomGenerator_Interface(t *testing.T) {
	var gen Generator = NewGenerator()
	slug, err := gen.Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(slug) != DefaultSlugLength {
		t.Errorf("expected length %d, got %d", DefaultSlugLength, len(slug))
	}
}

// TestSlugUniqueness_Sample checks for duplicate slugs within a small sample (1,000 items).
// Statistical note:
// The total keyspace for 6 base36 characters is 36^6 = 2,176,782,336.
// By the Birthday Paradox, the collision probability for n = 1,000 is approximately:
// P(collision) ≈ 1 - exp(-n^2 / (2 * N)) ≈ 1 - exp(-1,000,000 / 4,353,564,672) ≈ 0.023% (~1 in 4,350 runs).
// This test provides a reasonable sanity check that the generator is not stuck on a constant or narrow cycle,
// but the true guarantee of uniqueness in production MUST be enforced by PostgreSQL's UNIQUE constraint.
func TestSlugUniqueness_Sample(t *testing.T) {
	const sampleSize = 1000
	seen := make(map[string]struct{}, sampleSize)

	for i := 0; i < sampleSize; i++ {
		slug, err := Generate()
		if err != nil {
			t.Fatalf("generation failed at iteration %d: %v", i, err)
		}
		if _, exists := seen[slug]; exists {
			t.Logf("Notice: generated duplicate slug %q within %d iterations (statistically possible with random values)", slug, sampleSize)
		}
		seen[slug] = struct{}{}
	}

	// Sanity check: the set of unique slugs should be nearly identical to the sample size
	if len(seen) < sampleSize-5 {
		t.Fatalf("excessive collisions detected: generated %d unique values out of %d", len(seen), sampleSize)
	}
}

func TestGenerate_Concurrent(t *testing.T) {
	const numGoroutines = 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			slug, err := Generate()
			if err != nil {
				t.Errorf("concurrent generation failed: %v", err)
			}
			if len(slug) != DefaultSlugLength {
				t.Errorf("expected length %d, got %d", DefaultSlugLength, len(slug))
			}
		}()
	}

	wg.Wait()
}
