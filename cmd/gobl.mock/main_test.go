package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoot_Help(t *testing.T) {
	cmd := root()
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestRoot_Version(t *testing.T) {
	cmd := root()
	cmd.SetArgs([]string{"--version"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestRoot_DefaultRegime(t *testing.T) {
	cmd := root()
	cmd.SetArgs([]string{"--seed", "42"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestRoot_WithRegime(t *testing.T) {
	regimes := []string{"ES", "DE", "MX", "FR", "IT", "US"}
	for _, r := range regimes {
		t.Run(r, func(t *testing.T) {
			cmd := root()
			cmd.SetArgs([]string{"--regime", r, "--seed", "42"})
			err := cmd.Execute()
			assert.NoError(t, err)
		})
	}
}

func TestRoot_WithAddon(t *testing.T) {
	cmd := root()
	cmd.SetArgs([]string{"--regime", "ES", "--addon", "es-facturae-v3", "--seed", "42"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestRoot_WithType(t *testing.T) {
	types := []string{"credit-note", "corrective", "debit-note", "proforma"}
	for _, invType := range types {
		t.Run(invType, func(t *testing.T) {
			cmd := root()
			cmd.SetArgs([]string{"--regime", "ES", "--type", invType, "--seed", "42"})
			err := cmd.Execute()
			assert.NoError(t, err)
		})
	}
}

func TestRoot_Simplified(t *testing.T) {
	cmd := root()
	cmd.SetArgs([]string{"--regime", "ES", "--simplified", "--seed", "42"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestRoot_Lines(t *testing.T) {
	cmd := root()
	cmd.SetArgs([]string{"--regime", "ES", "--lines", "10", "--seed", "42"})
	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestRoot_Output(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")

	cmd := root()
	cmd.SetArgs([]string{"--regime", "ES", "--seed", "42", "-o", path})
	err := cmd.Execute()
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"$schema"`)
}

func TestRoot_InvalidRegime(t *testing.T) {
	cmd := root()
	cmd.SetArgs([]string{"--regime", "ZZ", "--seed", "42"})
	err := cmd.Execute()
	assert.Error(t, err)
}

func TestVersionString(t *testing.T) {
	assert.Equal(t, "dev", versionString())
}
