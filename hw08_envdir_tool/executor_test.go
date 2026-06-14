package main

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func tc(t *testing.T, cmd []string, env Environment, rv int, out, err string) {
	t.Helper()

	// Сохраняем оригинальный stdout
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// Создаём буфер
	r, w, _ := os.Pipe()
	os.Stdout = w // Перенаправляем stdout в pipe

	os.Setenv("HELLO", "SHOULD_REPLACE")
	os.Setenv("FOO", "SHOULD_REPLACE")
	os.Setenv("UNSET", "SHOULD_REPLACE")
	os.Setenv("ADDED", "from original env")
	os.Setenv("EMPTY", "SHOULD_REPLACE")

	wg := sync.WaitGroup{}

	// Горутина для чтения из pipe в буфер
	var buf bytes.Buffer
	wg.Go(func() {
		_, _ = io.Copy(&buf, r)
	})

	er, ew, _ := os.Pipe()
	os.Stderr = ew // Перенаправляем stderr в pipe

	// Горутина для чтения из pipe в буфер
	var ebuf bytes.Buffer
	wg.Go(func() {
		_, _ = io.Copy(&ebuf, er)
	})

	res := RunCmd(cmd, env)

	w.Close()
	os.Stdout = originalStdout

	ew.Close()
	os.Stderr = originalStderr

	wg.Wait()

	es := ebuf.String()
	os := buf.String()

	require.Equal(t, rv, res)
	require.Equal(t, err, es)
	require.Equal(t, out, os)
}

func TestRunCmd(t *testing.T) {
	t.Run("TC1", func(t *testing.T) {
		out := "HELLO is (\"hello\")\nBAR is (bar)\nFOO is (   foo\nwith new line)\nUNSET is ()\n"
		out += "ADDED is (from original env)\nEMPTY is ()\narguments are arg1=1 arg2=2\n"
		err := ""

		env := make(Environment)
		env["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		env["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		env["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
		env["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
		env["UNSET"] = EnvValue{Value: "", NeedRemove: true}

		tc(t, []string{"/bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2"}, env, 0, out, err)
	})

	t.Run("TC2", func(t *testing.T) {
		out := "HELLO is (\"hello\")\nBAR is (bar)\nFOO is (   foo\nwith new line\nwith new line)\nUNSET is ()\n"
		out += "ADDED is (from original env)\nEMPTY is ()\narguments are arg1=1 arg2=2 arg3 arg4=4\n"
		err := ""

		env := make(Environment)
		env["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		env["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		env["FOO"] = EnvValue{Value: "   foo\nwith new line\nwith new line", NeedRemove: false}
		env["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
		env["UNSET"] = EnvValue{Value: "", NeedRemove: true}

		tc(t, []string{"/bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2", "arg3", "arg4=4"}, env, 0, out, err)
	})

	t.Run("TC ERROR 1", func(t *testing.T) {
		out := "Error executing command: fork/exec /bin/bash2: no such file or directory\n"
		err := ""

		env := make(Environment)
		env["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		env["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		env["FOO"] = EnvValue{Value: "   foo\nwith new line\nwith new line", NeedRemove: false}
		env["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: true}
		env["UNSET"] = EnvValue{Value: "", NeedRemove: false}

		tc(t, []string{"/bin/bash2", "./testdata/echo.sh", "arg1=1", "arg2=2", "arg3", "arg4=4"}, env, -1, out, err)
	})

	t.Run("TC ERROR 2", func(t *testing.T) {
		out := "start executing\n"
		err := "error executing\n"

		env := make(Environment)
		env["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		env["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		env["FOO"] = EnvValue{Value: "   foo\nwith new line\nwith new line", NeedRemove: false}
		env["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: true}
		env["UNSET"] = EnvValue{Value: "", NeedRemove: false}

		tc(t, []string{"/bin/bash", "./testdata/err.sh", "arg1=1", "arg2=2", "arg3", "arg4=4"}, env, 127, out, err)
	})

	t.Run("TC ERROR 3", func(t *testing.T) {
		out := "Error executing command: invalid env variable name (ERR=PAR)\n"
		err := ""

		env := make(Environment)
		env["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		env["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		env["FOO"] = EnvValue{Value: "   foo\nwith new line\nwith new line", NeedRemove: false}
		env["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
		env["UNSET"] = EnvValue{Value: "", NeedRemove: true}
		env["ERR=PAR"] = EnvValue{Value: "Error", NeedRemove: false}

		tc(t, []string{"/bin/bash", "./testdata/err.sh", "arg1=1", "arg2=2", "arg3", "arg4=4"}, env, -1, out, err)
	})
}
