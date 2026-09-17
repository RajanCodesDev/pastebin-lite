package slug

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// DefaultSlugLength is the default character length of generated slugs.
const DefaultSlugLength = 6

// Charset defines URL-safe characters for slug generation (lowercase alphanumeric, base36).
const Charset = "abcdefghijklmnopqrstuvwxyz0123456789"

// Generator defines the interface for generating snippet slugs.
type Generator interface {
	Generate() (string, error)
}

// RandomGenerator implements Generator using cryptographically secure random bytes.
type RandomGenerator struct {
	length int
}

// NewGenerator returns a new RandomGenerator with the default length of 6 characters.
func NewGenerator() *RandomGenerator {
	return &RandomGenerator{length: DefaultSlugLength}
}

// NewCustomGenerator returns a RandomGenerator with a specified slug length.
func NewCustomGenerator(length int) (*RandomGenerator, error) {
	if length <= 0 {
		return nil, fmt.Errorf("slug length must be greater than zero, got %d", length)
	}
	return &RandomGenerator{length: length}, nil
}

// Generate creates a cryptographically secure random slug.
func (g *RandomGenerator) Generate() (string, error) {
	return GenerateCustom(g.length)
}

// Generate is a package-level convenience function generating a default 6-character slug.
func Generate() (string, error) {
	return GenerateCustom(DefaultSlugLength)
}

// GenerateCustom generates a slug of specified length using crypto/rand.
// It uses uniform random selection from the Charset without modulo bias.
func GenerateCustom(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("slug length must be greater than zero, got %d", length)
	}

	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(Charset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate random slug: %w", err)
		}
		result[i] = Charset[num.Int64()]
	}

	return string(result), nil
}
