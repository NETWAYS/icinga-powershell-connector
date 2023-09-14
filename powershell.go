package main

import (
	"strings"
)

// GetPowershellArgs returns remaining args and parse them as if we were powershell.exe.
// nolint:gocognit
func GetPowershellArgs(args []string) (command string, arguments map[string]interface{}) {
	arguments = map[string]interface{}{}
	l := len(args)

	for i := 0; i < l; i++ {
		arg := args[i]

		// ignore any flag with double dash (our flags)
		if strings.HasPrefix(arg, "--") {
			if i+1 < l && args[i+1][0] != '-' {
				i++ // skip the next one also
			}

			continue
		}

		// retrieve command
		if strings.EqualFold(arg, "-Command") || strings.EqualFold(arg, "-C") {
			if i+1 < l && args[i+1][0] != '-' {
				command = args[i+1]
				i++
			}

			continue
		}

		// warning and critical might be ranges
		if strings.ToLower(arg) == "-warning" || strings.ToLower(arg) == "-critical" {
			// No matter what happens, the next argument is a threshold value
			// This is a dirty hack to allow the usage of range expressions which might
			// begin with '-' (which then might be interpreted as options
			arguments[arg] = BuildPowershellType(args[i+1])
			i++

			continue
		}

		// all other flags
		if i+1 >= l || args[i+1][0] == '-' {
			// next argument is also a flag, so this is a switch
			arguments[arg] = true
		} else {
			arguments[arg] = BuildPowershellType(args[i+1])
			i++
		}
	}

	if command != "" {
		command = ParsePowershellTryCatch(command)
	}

	return
}

// nolint: gocritic
func BuildPowershellType(value string) interface{} {
	if strings.EqualFold(value, `$null`) {
		return nil
	} else if strings.EqualFold(value, `$true`) {
		return true
	} else if strings.EqualFold(value, `$false`) {
		return false
	} else if IsPowershellArray(value) {
		return ConvertPowershellArray(value)
	}

	value = strings.Trim(value, "\"")

	if len(value) >= 6 && value[0:4] == "'\\''" && value[len(value)-4:] == "'\\''" {
		return value[4 : len(value)-4]
	} else if value[0] == '\'' && value[len(value)-1] == '\'' {
		return value[1 : len(value)-1]
	}

	return value
}

// ConvertPowershellArray to a golang type.
//
// Examples:
//
//	@() -> []string{}
//	@('abc') -> []string{"abc"}
//	@('abc','def') -> []string{"abc","def"}
//
// nolint: gocognit
func ConvertPowershellArray(value string) []string {
	if value == "@()" || len(value) == 0 {
		return []string{}
	}

	if string(value[0]) == `@` {
		// strip array beginning and end
		value = value[2 : len(value)-1]
	}

	// Am I inside of a string
	insideString := false

	// Should the current character be escaped
	escapingMode := false

	// Which kind of quotes are we using right now? (Could be " or ')
	quoteMode := ``
	result := []string{}

	// Remember position
	marker := 0

	for i := range value {
		if value[i] == '"' && !escapingMode {
			if insideString && quoteMode == `"` {
				insideString = false
				quoteMode = ``
			} else {
				quoteMode = `"`
				insideString = true
			}

			continue
		}

		if value[i] == '\'' && !escapingMode {
			if insideString && quoteMode == `'` {
				insideString = false
				quoteMode = ``
			} else {
				quoteMode = `'`
				insideString = true
			}

			continue
		}

		if value[i] == ',' && !insideString {
			if value[i-1] == ',' {
				// Two consecutive commas
				result = append(result, "")
			} else {
				result = append(result, unquoteString(value[marker:i]))
			}
			// Point to char after comma
			marker = i + 1

			continue
		}

		if value[i] == '\\' && !escapingMode {
			escapingMode = true
			continue
		}

		if escapingMode {
			escapingMode = false
			continue
		}
	}

	// There might be a rest remaining
	if marker != len(value) {
		result = append(result, unquoteString(value[marker:]))
	}

	return result
}

func unquoteString(s string) string {
	if len(s) <= 1 {
		return s
	}

	if ((s[0]) == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}

	return s
}

// ParsePowershellTryCatch parses the actual command from a try/catch PowerShell code snippet.
//
// Examples:
//
//	 try { Use-Icinga -Minimal; } catch { <# something #> exit 3; };
//			Exit-IcingaExecutePlugin -Command 'Invoke-IcingaCheckUsedPartitionSpace'
//	 try { Use-Icinga -Minimal; } catch { <# something #> exit 3; }; Invoke-IcingaCheckUsedPartitionSpace
func ParsePowershellTryCatch(command string) string {
	command = strings.TrimSpace(command)
	// For now just parse the last word, dequote it and use it as command
	parts := strings.Split(command, " ")
	// Trim ' " and whitespaces
	return strings.Trim(parts[len(parts)-1], "'\" \t")
}

// nolint:gocognit
func IsPowershellArray(s string) bool {
	l := len(s)
	if l <= 2 {
		return false
	}

	if len(s) >= 3 && s[0] == '@' && s[1] == '(' && s[l-1] == ')' {
		return true
	}

	if !strings.Contains(s, ",") {
		return false
	}

	insideString := false
	escapingMode := false
	quoteMode := ``
	foundArraySeparator := false

	for i := range s {
		if string(s[i]) == `"` && !escapingMode {
			if insideString && quoteMode == `"` {
				insideString = false
				quoteMode = ``
			} else {
				quoteMode = `"`
				insideString = true
			}

			continue
		}

		if string(s[i]) == `'` && !escapingMode {
			if insideString && quoteMode == `'` {
				insideString = false
				quoteMode = ``
			} else {
				quoteMode = `'`
				insideString = true
			}

			continue
		}

		if string(s[i]) == `,` && !insideString {
			foundArraySeparator = true
		}

		if escapingMode {
			escapingMode = false
			continue
		}
	}

	return foundArraySeparator
}
