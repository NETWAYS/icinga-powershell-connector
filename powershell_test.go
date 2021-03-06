package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPowershellArgs(t *testing.T) {
	command, args := GetPowershellArgs([]string{"-C", "Invoke-IcingaCheckUsedPartitionSpace", "-Warning", "80"})
	assert.Equal(t, "Invoke-IcingaCheckUsedPartitionSpace", command)
	assert.Equal(t, map[string]interface{}{"-Warning": "80"}, args)

	command, args = GetPowershellArgs([]string{"-Switch", "-Warning", "80"})
	assert.Equal(t, "", command)
	assert.Equal(t, map[string]interface{}{"-Switch": true, "-Warning": "80"}, args)

	command, args = GetPowershellArgs([]string{"-Switch"})
	assert.Equal(t, "", command)
	assert.Equal(t, map[string]interface{}{"-Switch": true}, args)

	command, args = GetPowershellArgs([]string{"--powershell-insecure"})
	assert.Equal(t, "", command)
	assert.Equal(t, map[string]interface{}{}, args)

	command, args = GetPowershellArgs([]string{
		"--powershell-api",
		"https://battlestation:5668",
		"--powershell-insecure",
		"-C",
		"try { Use-Icinga -Minimal; } catch { <# some error #>; exit 3; }; "+
			"Exit-IcingaExecutePlugin -Command 'Invoke-IcingaCheckUsedPartitionSpace' ",
		"-Warning",
		"80",
		"-Critical",
		"95",
		"-Include",
		"@()",
		"-Exclude",
		"@('abc')",
		"-Verbosity",
		"2",
	})
	assert.Equal(t, "Invoke-IcingaCheckUsedPartitionSpace", command)
	assert.Equal(t, map[string]interface{}{
		"-Critical":"95", "-Verbosity":"2", "-Warning":"80", "-Exclude":[]string{"abc"}, "-Include":[]string{},
	}, args)
}

func TestParsePowershellTryCatch(t *testing.T) {
	command := ParsePowershellTryCatch(
		"try { Use-Icinga -Minimal; } catch { <# something #> exit 3; }; "+
			"Exit-IcingaExecutePlugin -Command 'Invoke-IcingaCheckUsedPartitionSpace' ")
	assert.Equal(t, "Invoke-IcingaCheckUsedPartitionSpace", command)

	command = ParsePowershellTryCatch(
		"try { Use-Icinga } catch { <# something #> exit 3; }; Invoke-IcingaCheckUsedPartitionSpace ")
	assert.Equal(t, "Invoke-IcingaCheckUsedPartitionSpace", command)

	command = ParsePowershellTryCatch("Invoke-IcingaCheckUsedPartitionSpace")
	assert.Equal(t, "Invoke-IcingaCheckUsedPartitionSpace", command)
}

