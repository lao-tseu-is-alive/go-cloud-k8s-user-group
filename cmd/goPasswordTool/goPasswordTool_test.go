package main

import (
	"bytes"
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMainExecution(t *testing.T) {
	// We manipulate the Args to set them up for the testcases
	// after this test we restore the initial args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	cases := []struct {
		Name           string
		Args           []string
		ExpectedExit   int
		ExpectedOutput string
	}{
		{"it should return error if password flag is not set", []string{}, 1, "💥💥 ERROR: 'missing parameter"},
		{"it should output valid sha256", []string{"-password", "testString"}, 0, "sha256: 4acf0b39d9c4766709a3689f553ac01ab550545ffa4544dfc0b2cea82fba02a3"},
	}
	for _, tc := range cases {
		// this call is required because otherwise flags panics, if args are set between flag.Parse calls
		flag.CommandLine = flag.NewFlagSet(tc.Name, flag.ExitOnError)
		// we need a value to set Args[0] to, cause flag begins parsing at Args[1]
		os.Args = append([]string{"goPasswordTool"}, tc.Args...)
		var buf bytes.Buffer
		actualExit := realMain(&buf)
		if tc.ExpectedExit != actualExit {
			t.Errorf("Wrong exit code for args: %v, expected: %v, got: %v",
				tc.Args, tc.ExpectedExit, actualExit)
		}
		actualOutput := buf.String()
		assert.Contains(t, actualOutput, tc.ExpectedOutput)
	}
}
