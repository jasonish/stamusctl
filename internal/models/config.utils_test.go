package models

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNestMap(t *testing.T) {
	result := nestMap(map[string]interface{}{
		"Release.truc":   "test",
		"Release.machin": "test",
	})
	assert.Equal(t, map[string]interface{}{
		"Release": map[string]interface{}{
			"truc":   "test",
			"machin": "test",
		},
	}, result)
}

func TestNestMapWithMoreNesting(t *testing.T) {
	result := nestMap(map[string]interface{}{
		"Release.deep.truc":   "test",
		"Release.deep.machin": "test",
	})
	assert.Equal(t, map[string]interface{}{
		"Release": map[string]interface{}{
			"deep": map[string]interface{}{
				"truc":   "test",
				"machin": "test",
			},
		},
	}, result)
}
