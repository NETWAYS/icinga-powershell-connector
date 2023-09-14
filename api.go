package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

type RestAPI struct {
	URL    string
	Client *http.Client
}

func (a RestAPI) ExecuteCheck(command string, arguments map[string]interface{}, timeout uint32) (*APICheckResult, error) { //nolint:lll
	// Build body
	body, err := json.Marshal(arguments)
	if err != nil {
		return nil, fmt.Errorf("could not build JSON body: %w", err)
	}

	// With timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// Build request
	requestURL := a.URL + "/v1/checker?command=" + url.QueryEscape(command)

	log.WithFields(log.Fields{
		"body": string(body),
		"url":  requestURL,
	}).Debug("sending request")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := a.getClient().Do(req)

	if err != nil {
		// We want to override the context error message to be more expressive
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("timeout during HTTP request: %w", err)
		}

		return nil, fmt.Errorf("executing API request failed: %w", err)
	}

	defer resp.Body.Close()

	// Read response
	resultBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read result: %w", err)
	}

	log.WithField("body", string(resultBody)).Debug("received response")

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request not successful code=%d: %s", resp.StatusCode, string(resultBody))
	}

	// Parse result
	var result APICheckResults

	err = json.Unmarshal(resultBody, &result)
	if err != nil {
		return nil, fmt.Errorf("could not parse result JSON: %w", err)
	}

	// return first check result
	for _, r := range result {
		return &r, nil
	}

	return nil, fmt.Errorf("no check result in API response")
}

func (a *RestAPI) getClient() *http.Client {
	if a.Client == nil {
		a.Client = http.DefaultClient
	}

	return a.Client
}
