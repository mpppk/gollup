package cmd_test

import (
	"bytes"
	"fmt"
	"github.com/spf13/afero"
	"strings"
	"testing"

	"github.com/mpppk/cli-template/cmd"
)

func TestSum(t *testing.T) {
	cases := []struct {
		command string
		want    string
	}{
		{command: "sum -- -1 2", want: "1\n"},
		{command: "sum --norm -- -1 2", want: "3\n"},
	}

	for _, c := range cases {
		buf := new(bytes.Buffer)
		rootCmd, err := cmd.NewRootCmd(afero.NewMemMapFs())
		if err != nil {
			t.Errorf("failed to create rootCmd: %s", err)
		}
		rootCmd.SetOut(buf)
		cmdArgs := strings.Split(c.command, " ")
		rootCmd.SetArgs(cmdArgs)
		if err := rootCmd.Execute(); err != nil {
			t.Errorf("failed to execute rootCmd: %s", err)
		}

		get := buf.String()
		if c.want != get {
			t.Errorf("unexpected response: want:%q, get:%q", c.want, get)
		}
	}
}

func TestSumWithOutFile(t *testing.T) {
	testFilePath := "test.txt"
	cases := []struct {
		command string
		want    string
	}{
		{command: fmt.Sprintf("sum --out %s -- -1 2", testFilePath), want: "1"},
	}

	for _, c := range cases {
		fs := afero.NewMemMapFs()
		rootCmd, err := cmd.NewRootCmd(fs)
		if err != nil {
			t.Errorf("failed to create rootCmd: %s", err)
		}
		cmdArgs := strings.Split(c.command, " ")
		rootCmd.SetArgs(cmdArgs)
		if err := rootCmd.Execute(); err != nil {
			t.Errorf("failed to execute rootCmd: %s", err)
		}
		byteContents, err := afero.ReadFile(fs, testFilePath)
		if err != nil {
			t.Fatal(err)
		}
		get := string(byteContents)
		if c.want != get {
			t.Errorf("unexpected response: want:%q, get:%q", c.want, get)
		}
	}
}


