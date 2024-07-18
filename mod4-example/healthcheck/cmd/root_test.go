package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/marianina8/gocodecli/mod4/healthcheck/logger"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func setupTestLogger() {
	l = logger.New(logFile, false, true, output)
}

func executeCommandC(cmd *cobra.Command, args ...string) (string, error) {
	b := bytes.NewBufferString("")
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	outC := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	cmd.SetOut(b)
	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()

	w.Close()
	os.Stdout = oldStdout
	cmdOutput := b.String() + <-outC
	return cmdOutput, err
}

func TestRootCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "unknown command",
			args:           []string{"nonexistant"},
			expectedOutput: "unknown command",
			expectedError:  true,
		},
		{
			name:           "unexpected flag",
			args:           []string{"--unexpected-flag"},
			expectedOutput: "unknown flag: --unexpected-flag",
			expectedError:  true,
		},
		{
			name:          "missing mandatory url",
			args:          []string{"check"},
			expectedError: true,
		},
		{
			name:           "combined flags unexpected order",
			args:           []string{"check", "--output", "json", "http://example.com"},
			expectedOutput: "successful check",
			expectedError:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := executeCommandC(rootCmd, tc.args...)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if err != nil {
				assert.Contains(t, err.Error(), tc.expectedOutput)
			}
		})
	}
}
