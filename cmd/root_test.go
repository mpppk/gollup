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
		name         string
		command      string
		wantFilePath string
	}{
		{
			name: "single_pkg",
			command: fmt.Sprintf("%s",
				filepath.Join(testDir, "single_pkg"),
			),
			wantFilePath: filepath.Join(testDir, "single_pkg", "want", "want.go.test"),
		},
		{
			// execute with entry point
			name: "single_pkg with entry point",
			command: fmt.Sprintf("--entrypoint main.main %s",
				filepath.Join(testDir, "single_pkg"),
			),
			wantFilePath: filepath.Join(testDir, "single_pkg", "want", "want.go.test"),
		},
		{
			name: "multi_pkg",
			command: fmt.Sprintf("%s %s",
				filepath.Join(testDir, "multi_pkg"),
				filepath.Join(testDir, "multi_pkg", "lib"),
			),
			wantFilePath: filepath.Join(testDir, "multi_pkg", "want", "want.go.test"),
		},
		{
			// execute with entry point
			name: "multi_pkg with entry point",
			command: fmt.Sprintf("--entrypoint main.main %s %s",
				filepath.Join(testDir, "multi_pkg"),
				filepath.Join(testDir, "multi_pkg", "lib"),
			),
			wantFilePath: filepath.Join(testDir, "multi_pkg", "want", "want.go.test"),
		},
		{
			name: "comments",
			command: fmt.Sprintf("%s %s",
				filepath.Join(testDir, "comments"),
				filepath.Join(testDir, "comments", "lib"),
			),
			wantFilePath: filepath.Join(testDir, "comments", "want", "want.go.test"),
		},
		{
			name: "struct",
			command: fmt.Sprintf("%s %s",
				filepath.Join(testDir, "struct"),
				filepath.Join(testDir, "struct", "lib"),
			),
			wantFilePath: filepath.Join(testDir, "struct", "want", "want.go.test"),
		},
		{
			name: "const",
			command: fmt.Sprintf("%s %s",
				filepath.Join(testDir, "const"),
				filepath.Join(testDir, "const", "lib"),
			),
			wantFilePath: filepath.Join(testDir, "const", "want", "want.go.test"),
		},
		{
			name: "compro",
			command: fmt.Sprintf("%s %s",
				filepath.Join(testDir, "compro"),
				filepath.Join(testDir, "compro", "lib"),
			),
			wantFilePath: filepath.Join(testDir, "compro", "want", "want.go.test"),
		},
		{
			name: "params_and_results",
			command: fmt.Sprintf("%s %s",
				filepath.Join(testDir, "params_and_results"),
				filepath.Join(testDir, "params_and_results", "lib"),
			),
			wantFilePath: filepath.Join(testDir, "params_and_results", "want", "want.go.test"),
		},
		{
			command: fmt.Sprintf("%s %s",
				filepath.Join(testDir, "abc007C"),
				filepath.Join(testDir, "abc007C", "lib"),
			),
			wantFilePath: filepath.Join(testDir, "abc007C", "want", "want.go.test"),
		},
		{
			command: fmt.Sprintf("%s %s",
				filepath.Join(testDir, "atc001B"),
				filepath.Join(testDir, "atc001B", "lib"),
			),
			wantFilePath: filepath.Join(testDir, "atc001B", "want", "want.go.test"),
		},
		// duplicated name const is not supported yet
		//{
		//	command: fmt.Sprintf("%s %s",
		//		filepath.Join(testDir, "dup_const"),
		//		filepath.Join(testDir, "dup_const", "lib"),
		//	),
		//	wantFilePath: filepath.Join(testDir, "dup_const", "want", "want.go.test"),
		//},
		// duplicated name struct is not supported yet
		//{
		//	command: fmt.Sprintf("%s %s",
		//		filepath.Join(testDir, "dup_struct"),
		//		filepath.Join(testDir, "dup_struct", "lib"),
		//	),
		//	wantFilePath: filepath.Join(testDir, "dup_struct", "want", "want.go.test"),
		//},
		// duplicated name struct is not supported yet
		//{
		//	command: fmt.Sprintf("%s %s",
		//		filepath.Join(testDir, "nested_struct"),
		//		filepath.Join(testDir, "nested_struct", "lib"),
		//	),
		//	wantFilePath: filepath.Join(testDir, "nested_struct", "want", "want.go.test"),
		//},
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
		get = removeCarriageReturn(get)
		contents, err := ioutil.ReadFile(c.wantFilePath)
		if err != nil {
			t.Fail()
		}
		want := string(contents)
		want = removeCarriageReturn(want)
		if want != get {
			t.Errorf(":%s unexpected response: want:\n%s\nget:\n%s", c.name, want, get)
		}
	}
}

func removeCarriageReturn(s string) string {
	return strings.Replace(s, "\r", "", -1)
}
