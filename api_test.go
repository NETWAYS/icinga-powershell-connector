package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func TestApiTimeout(t *testing.T) {
	api := RestAPI{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wait for the context timeout to kick in
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	api.URL = srv.URL

	args := make(map[string]interface{})

	_, err := api.ExecuteCheck("command", args, 1)

	if err == nil {
		t.Error("Expected error got nil")
	}

	actual := err.Error()
	expected := "timeout during HTTP request"

	if !strings.Contains(actual, expected) {
		t.Error("\nActual: ", actual, "\nExpected: ", expected)
	}
}
