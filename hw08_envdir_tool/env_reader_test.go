package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here

	t.Run("TC1", func(t *testing.T) {
		expEnv := make(Environment)

		env, err := ReadDir("./testdata/env")

		require.Empty(t, err)

		expEnv["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		expEnv["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		expEnv["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
		expEnv["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
		expEnv["UNSET"] = EnvValue{Value: "", NeedRemove: false}

		require.Equal(t, env, expEnv)
	})

	t.Run("TC2", func(t *testing.T) {
		expEnv := make(Environment)

		env, err := ReadDir("./testdata/env2")

		require.Empty(t, err)

		expEnv["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		expEnv["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		expEnv["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
		expEnv["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
		expEnv["UNSET"] = EnvValue{Value: "", NeedRemove: false}

		require.Equal(t, env, expEnv)
	})

	t.Run("TC3", func(t *testing.T) {
		expEnv := make(Environment)

		env, err := ReadDir("./testdata/env_err")

		require.Empty(t, err)

		expEnv["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		expEnv["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		expEnv["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
		expEnv["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
		expEnv["UNSET"] = EnvValue{Value: "", NeedRemove: false}
		expEnv["ENV=PARAM"] = EnvValue{Value: "Error", NeedRemove: false}

		require.Equal(t, env, expEnv)
	})

	t.Run("TC ERROR 1", func(t *testing.T) {
		expErr := "open ./testdata/no_exist: no such file or directory"

		_, err := ReadDir("./testdata/no_exist")

		require.Equal(t, err.Error(), expErr)
	})
}
