package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSkDir_IdempotentOnRepeatCalls(t *testing.T) {
	tmpDir := t.TempDir()
	origSkDir := skDir
	skDir = filepath.Join(tmpDir, ".sk")
	t.Cleanup(func() { skDir = origSkDir })

	// First call creates the dir; second call must not return an error.
	assert.NoError(t, createSkDir())
	assert.NoError(t, createSkDir())
}

func TestResolveSelection_ExactMatch(t *testing.T) {
	selection, err := resolveSelection([]string{"dev-cluster", "prod-cluster"}, "prod-cluster")

	assert.NoError(t, err)
	assert.Equal(t, "prod-cluster", selection)
}

func TestResolveSelection_FuzzyMatchSelectsOnlyMatch(t *testing.T) {
	selection, err := resolveSelection([]string{"dev-cluster", "prod-cluster"}, "pc")

	assert.NoError(t, err)
	assert.Equal(t, "prod-cluster", selection)
}

func TestResolveSelection_EmptyInputSelectsFirstOption(t *testing.T) {
	selection, err := resolveSelection([]string{"current-cluster", "other-cluster"}, "")

	assert.NoError(t, err)
	assert.Equal(t, "current-cluster", selection)
}

func TestResolveSelection_RejectsNoMatch(t *testing.T) {
	selection, err := resolveSelection([]string{"dev-cluster", "prod-cluster"}, "xyz")

	assert.Error(t, err)
	assert.Empty(t, selection)
}

func TestResolveSelection_RejectsAmbiguousFuzzyMatch(t *testing.T) {
	selection, err := resolveSelection([]string{"prod-eu", "prod-us"}, "pr")

	assert.EqualError(t, err, `"pr" matches multiple selections; type a more specific selection`)
	assert.Empty(t, selection)
}
