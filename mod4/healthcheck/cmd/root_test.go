package cmd

import (
	"bytes"
	"io"
	"os"

	"github.com/marianina8/gocodecli/mod4/healthcheck/logger"
	"github.com/spf13/cobra"
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
