package main

import (
	flag "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func quietTest() func() {
	null, _ := os.Open(os.DevNull)

	stdOut := os.Stdout
	stdErr := os.Stderr

	os.Stdout = null
	os.Stderr = null

	return func() {
		defer func() {
			_ = null.Close()
		}()

		os.Stdout = stdOut
		os.Stderr = stdErr
	}
}

func TestParseConfigFromFlags(t *testing.T) {
	defer quietTest()()

	_, err := ParseConfigFromFlags([]string{"--help"})
	assert.ErrorIs(t, err, flag.ErrHelp)

	_, err = ParseConfigFromFlags([]string{"--version"})
	assert.ErrorIs(t, err, ErrVersionRequested)

	_, err = ParseConfigFromFlags([]string{})
	assert.ErrorIs(t, err, ErrNoCommand)

	// Just our flags
	config, err := ParseConfigFromFlags([]string{
		"--command", "Invoke-IcingaCheckUsedPartitionSpace",
		"--api", "https://localhost:8888"})
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, config.Command, "Invoke-IcingaCheckUsedPartitionSpace")

	// Using wrapper mode
	config, err = ParseConfigFromFlags([]string{"-Command", "Invoke-IcingaCheckUsedPartitionSpace", "-argWithH"})
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, config.Command, "Invoke-IcingaCheckUsedPartitionSpace")

	config, err = ParseConfigFromFlags([]string{
		"-C", "try { Use-Icinga -Minimal; ... -Command 'Invoke-IcingaCheckUsedPartitionSpace' "})
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, config.Command, "Invoke-IcingaCheckUsedPartitionSpace")

	config, err = ParseConfigFromFlags([]string{"-C", "Invoke-IcingaCheckUsedPartitionSpace", "-Warning", "80"})
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, config.Command, "Invoke-IcingaCheckUsedPartitionSpace")

	// Mixing both arguments
	config, err = ParseConfigFromFlags([]string{
		"--command", "Invoke-SomethingElse",
		"--api", "https://localhost:8888",
		"-C", "try { Use-Icinga -Minimal; ... -Command 'Invoke-IcingaCheckUsedPartitionSpace' ", "-argWithH"})
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, config.API, "https://localhost:8888")
	assert.Equal(t, config.Command, "Invoke-IcingaCheckUsedPartitionSpace")
	assert.Contains(t, config.Arguments, "-argWithH")
}

func TestSplitPowerShellArguments(t *testing.T) {
	flags, powerShellArgs := SplitPowerShellArguments([]string{
		"--command", "Invoke-SomethingElse",
		"--api", "https://localhost:8888",
		"-C", "try { Use-Icinga -Minimal; ... -Command 'Invoke-IcingaCheckUsedPartitionSpace' ", "-argWithH"})

	assert.Equal(t, []string{
		"--command", "Invoke-SomethingElse",
		"--api", "https://localhost:8888"},
		flags)
	assert.Equal(t, []string{
		"-C", "try { Use-Icinga -Minimal; ... -Command 'Invoke-IcingaCheckUsedPartitionSpace' ", "-argWithH"},
		powerShellArgs)
}
