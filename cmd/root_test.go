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
			command: fmt.Sprintf("--dirs %s %s",
				filepath.Join(testDir, "test1"),
				filepath.Join(testDir, "test1", "main.go"),
			),
			wantFilePath: filepath.Join(testDir, "test1", "want", "want.go.test"),
		},
		{
			command: fmt.Sprintf("--dirs %s %s",
				strings.Join([]string{
					filepath.Join(testDir, "test2"),
					filepath.Join(testDir, "test2", "lib"),
				}, ","),
				filepath.Join(testDir, "test2", "main.go"),
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
		fmt.Println(get)
		contents, err := ioutil.ReadFile(c.wantFilePath)
		if err != nil {
			t.Fail()
		}
		want := string(contents)
		if want != get {
			t.Errorf("unexpected response: want:%q, get:%q", want, get)
		}
	}
}
