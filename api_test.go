package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type ApiTest struct {
	name     string
	server   *httptest.Server
	api      RestAPI
	expected APICheckResult
}

func TestApiCmd(t *testing.T) {
	tests := []ApiTest{
		{
			name: "api-simple-test",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"Invoke-IcingaCheckFoo": {"exitcode": 0, "checkresult": "[OK] \"foo\"", "perfdata": ["'foo'=1.00%;;;0;100 "]}}`))
			})),
			api: RestAPI{},
			expected: APICheckResult{
				ExitCode:    0,
				CheckResult: "[OK] \"foo\"",
				Perfdata:    APIPerfdataList{},
			},
		},
	}

	var (
		actual *APICheckResult
		err    error
	)

	args := make(map[string]interface{})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.server.Close()
			test.api.URL = test.server.URL

			actual, err = test.api.ExecuteCheck("command", args, 10)

			if err != nil {
				t.Error(err)
			}

			if actual.CheckResult != test.expected.CheckResult {
				t.Error("\nActual: ", actual.CheckResult, "\nExpected: ", test.expected.CheckResult)
			}
		})
	}
}
