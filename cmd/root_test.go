package cmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mpppk/gollup/cmd"

	"github.com/spf13/afero"
)

const testDir = "../testdata"

func TestRoot(t *testing.T) {
	cases := []struct {
		command      string
		wantFilePath string
	}{
		{
			command: fmt.Sprintf("%s",
				filepath.Join(testDir, "test1"),
			),
			wantFilePath: filepath.Join(testDir, "test1", "want", "want.go.test"),
		},
		{
			// execute with entry point
			command: fmt.Sprintf("--entrypoint main.main %s",
				filepath.Join(testDir, "test1"),
			),
			wantFilePath: filepath.Join(testDir, "test1", "want", "want.go.test"),
		},
		{
			command: fmt.Sprintf("%s %s",
				filepath.Join(testDir, "test2"),
				filepath.Join(testDir, "test2", "lib"),
			),
			wantFilePath: filepath.Join(testDir, "test2", "want", "want.go.test"),
		},
		{
			// execute with entry point
			command: fmt.Sprintf("--entrypoint main.main %s %s",
				filepath.Join(testDir, "test2"),
				filepath.Join(testDir, "test2", "lib"),
			),
			wantFilePath: filepath.Join(testDir, "test2", "want", "want.go.test"),
		},
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		rootCmd, err := cmd.NewRootCmd(afero.NewMemMapFs())
		if err != nil {
			t.Errorf("failed to create rootCmd: %s", err)
		}
		rootCmd.SetOut(buf)
		rootCmd.SetErr(buf)
		cmdArgs := strings.Split(c.command, " ")
		rootCmd.SetArgs(cmdArgs)
		if err := rootCmd.Execute(); err != nil {
			t.Errorf("failed to execute rootCmd: %s", err)
		}

		get := buf.String()
		trimmedGet := removeCarriageReturn(get)
		contents, err := ioutil.ReadFile(c.wantFilePath)
		if err != nil {
			t.Fail()
		}
		want := string(contents)
		trimmedWant := removeCarriageReturn(want)
		if want != get {
			t.Errorf("unexpected response: want:\n%s\nget:\n%s", trimmedWant, trimmedGet)
		}
	}
}

func removeCarriageReturn(s string) string {
	return strings.Replace(s, "\r", "", -1)
}
