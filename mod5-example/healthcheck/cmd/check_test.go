package cmd

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		expectedError error
	}{
		{
			"Empty URL",
			"",
			&URLValidationError{
				URL:    "",
				Detail: "missing scheme",
			},
		},
		{
			"Invalid URL",
			"http//googl",
			&URLValidationError{
				URL:    "http//googl",
				Detail: "missing scheme",
			},
		},
		{
			"Valid http URL",
			"http://www.google.com",
			nil,
		},
		{
			"Valid https URL",
			"https://www.google.com",
			nil,
		},
		{
			"URL with no host",
			"http://",
			&URLValidationError{
				URL:    "http://",
				Detail: "missing host",
			},
		},
		{
			"URL with IP",
			"http://192.168.0.1",
			nil,
		},
		{
			"URL with port",
			"http://localhost:8080",
			nil,
		},
		{
			"URL with special characters",
			"http://example.com/path?name=val#anchor",
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualError := isValidURL(tc.url)
			if actualError != nil && tc.expectedError != nil {
				if actualError.Error() != tc.expectedError.Error() {
					t.Errorf("isValidURL(%q)=%v, expected error %v", tc.url, actualError, tc.expectedError)
				}
			}
			if actualError != nil && tc.expectedError == nil {
				t.Errorf("isValidURL(%q)=%v, expected no error", tc.url, actualError)
			}
			if actualError == nil && tc.expectedError != nil {
				t.Errorf("isValidURL(%q)=nil, expected error: %v", tc.url, tc.expectedError)
			}
		})
	}
}

func TestCheckURL(t *testing.T) {
	setupTestLogger()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Override the ExitFunction
	originalExit := ExitFunction
	defer func() { ExitFunction = originalExit }()
	var exitCode int
	ExitFunction = func(code int) {
		exitCode = code
	}
	tests := []struct {
		name              string
		url               string
		threshold         float64
		retries           int
		mockResp          httpmock.Responder
		expected          bool
		expectedExitCode  int
		delayBetweenCalls time.Duration
	}{
		{
			name:      "Successful Request - 200 OK",
			url:       "http://www.google.com",
			threshold: 2.0,
			retries:   1,
			mockResp:  httpmock.NewStringResponder(200, "OK"),
			expected:  true,
		},
		{
			name:      "Server Error - 500 Server Error",
			url:       "http://www.tripadvisor.com",
			threshold: 2.0,
			retries:   1,
			mockResp:  httpmock.NewStringResponder(500, "Server Error"),
			expected:  false,
		},
		{
			name:      "Timeout Exceeded",
			url:       "http://www.tripadvisor.com",
			threshold: 0.01,
			retries:   1,
			mockResp: func(req *http.Request) (*http.Response, error) {
				time.Sleep(50 * time.Millisecond)
				return httpmock.NewStringResponse(200, "OK"), nil
			},
			expected: true,
		},
		{
			name:             "Failed request",
			url:              "http://www.example.com",
			threshold:        2.0,
			retries:          1,
			mockResp:         httpmock.NewErrorResponder(context.DeadlineExceeded),
			expected:         false,
			expectedExitCode: 1, //expected code when context is canceled
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, tc.url, tc.mockResp)
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			actual := checkURL(ctx, tc.url, tc.threshold, tc.retries)
			assert.Equal(t, tc.expected, actual)
			if tc.expectedExitCode != 0 {
				assert.Equal(t, tc.expectedExitCode, exitCode, "expected exit code does not match")
			}
		})
	}

}

func TestRun_OutputTable(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "http://example.com", httpmock.NewStringResponder(200, "OK"))
	var buf bytes.Buffer
	outputWriter = &buf
	_, err := executeCommandC(rootCmd, "check", "--output", "table", "http://example.com")
	output := buf.String()

	assert.NoError(t, err)
	assert.Contains(t, output, "Up", "Expected table output to contain 'Up'")
	outputWriter = os.Stdout
}

func TestRun_MultipleURLs(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodGet, "http://example1.com", httpmock.NewStringResponder(200, "OK"))
	httpmock.RegisterResponder(http.MethodGet, "http://example2.com", httpmock.NewStringResponder(200, "OK"))
	output, err := executeCommandC(rootCmd, "check", "--output", "text", "http://example1.com", "http://example2.com")
	assert.NoError(t, err)
	assert.Contains(t, output, "successful check")
}
