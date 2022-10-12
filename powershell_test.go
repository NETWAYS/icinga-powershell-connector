package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		"try { Use-Icinga -Minimal; } catch { <# some error #>; exit 3; }; " +
			"Exit-IcingaExecutePlugin -Command 'Invoke-IcingaCheckUsedPartitionSpace' ",
		"-Warning",
		"80",
		"-Critical",
		"95",
		"-Include",
		"@()",
		"-Exclude",
		"@('abc','def')",
		"-Verbosity",
		"2",
	})
	assert.Equal(t, "Invoke-IcingaCheckUsedPartitionSpace", command)
	assert.Equal(t, map[string]interface{}{
		"-Critical": "95", "-Verbosity": "2", "-Warning": "80", "-Exclude": []string{"abc", "def"}, "-Include": []string{},
	}, args)
}

func TestPowershellArrayConversionEmpty(t *testing.T) {
	assert.Equal(t, []string{}, ConvertPowershellArray("@()"))
	assert.Equal(t, []string{}, ConvertPowershellArray(""))
}

func TestPowershellArrayTest(t *testing.T) {
	assert.Equal(t, true, IsPowershellArray(`@('abc',"de\"f',15)`))
	assert.Equal(t, true, IsPowershellArray(`'abc',"de\"f',15`))
	assert.Equal(t, true, IsPowershellArray("'',''"))

	assert.Equal(t, false, IsPowershellArray(`'';"de\"f';1`))
	assert.Equal(t, false, IsPowershellArray(""))
}

func TestPowershellArrayConversion(t *testing.T) {
	assert.Equal(t, []string{"abc", `de\"f`, "15"}, ConvertPowershellArray(`@('abc',"de\"f",15)`))
	assert.Equal(t, []string{"abc", `de\"f`, "15"}, ConvertPowershellArray(`'abc',"de\"f",15`))

	assert.Equal(t,
		[]string{"ASFASDFASF[]}", "", "", "1523423", "1"},
		ConvertPowershellArray(`'ASFASDFASF[]}',,"",1523423,1`))
}

func TestParsePowershellTryCatch(t *testing.T) {
	command := ParsePowershellTryCatch(
		"try { Use-Icinga -Minimal; } catch { <# something #> exit 3; }; " +
			"Exit-IcingaExecutePlugin -Command 'Invoke-IcingaCheckUsedPartitionSpace' ")
	assert.Equal(t, "Invoke-IcingaCheckUsedPartitionSpace", command)

	command = ParsePowershellTryCatch(
		"try { Use-Icinga } catch { <# something #> exit 3; }; Invoke-IcingaCheckUsedPartitionSpace ")
	assert.Equal(t, "Invoke-IcingaCheckUsedPartitionSpace", command)

	command = ParsePowershellTryCatch("Invoke-IcingaCheckUsedPartitionSpace")
	assert.Equal(t, "Invoke-IcingaCheckUsedPartitionSpace", command)
}

func TestPowershellQuotes(t *testing.T) {
	command, arguments := GetPowershellArgs([]string{"-Command", "Invoke-IcingaCheckPerfCounter", "-PerfCounter", `"'\TCPv4\Connections Established'"`})

	assert.Equal(t, "Invoke-IcingaCheckPerfCounter", command)
	expected_args := map[string]interface{}{
		"-PerfCounter": "\\TCPv4\\Connections Established",
	}
	assert.Equal(t, expected_args, arguments)
}
