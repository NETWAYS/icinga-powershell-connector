package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIPerfdataList_UnmarshalJSON(t *testing.T) {
	var pl APIPerfdataList

	err := json.Unmarshal([]byte("{}"), &pl)
	assert.NoError(t, err)
	assert.Equal(t, APIPerfdataList{}, pl)

	err = json.Unmarshal([]byte(`["a", "b", "c"]`), &pl)
	assert.NoError(t, err)
	assert.Equal(t, APIPerfdataList{"a", "b", "c"}, pl)
}

func TestAPIPerfdataList_String(t *testing.T) {

	testcases := []struct {
		result   APICheckResult
		expected string
	}{
		{
			result: APICheckResult{
				ExitCode:    1,
				CheckResult: "foo",
				Perfdata:    APIPerfdataList{"a", "b", "c"},
			},
			expected: "foo\n| a b c\n",
		},
		{
			result: APICheckResult{
				ExitCode:    1,
				CheckResult: "foo",
				Perfdata:    APIPerfdataList{},
			},
			expected: "foo\n",
		},
	}
	for _, test := range testcases {
		assert.Equal(t, test.expected, test.result.String())
	}
}
