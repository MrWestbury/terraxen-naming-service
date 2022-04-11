package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolvePattern(t *testing.T) {
	testData := map[string]string{
		"test":      "adam",
		"something": "or other",
		"foo":       "bar",
	}

	patterns := map[string]string{
		"{test}-string":        "adam-string",
		"{test}-{test}-string": "adam-adam-string",
		"Hello {foo}":          "Hello bar",
		"{test} {something}":   "adam or other",
	}

	for ptn, expected := range patterns {
		t.Run(ptn, func(tx *testing.T) {
			res, _ := ResolvePattern(ptn, testData)
			assert.Equal(tx, expected, res, ptn)
		})
	}

}
