package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func tc(t *testing.T, fromPath, resFile string, offset, limit int64) {
	t.Helper()
	out, err := os.CreateTemp("/tmp", "out")
	require.Empty(t, err)
	out.Close()
	toPath := out.Name()

	defer os.Remove(toPath)

	fmt.Println(toPath)

	err = Copy(fromPath, toPath, offset, limit)
	require.Empty(t, err)
	require.FileExists(t, toPath)

	outData, err1 := os.ReadFile(toPath)
	require.Empty(t, err1)

	resData, err2 := os.ReadFile(resFile)
	require.Empty(t, err2)

	require.Equal(t, outData, resData)
}

func etc(t *testing.T, fromPath string, resError string, offset, limit int64) {
	t.Helper()
	out, err := os.CreateTemp("/tmp", "out")
	require.Empty(t, err)
	out.Close()
	toPath := out.Name()

	defer os.Remove(toPath)

	fmt.Println(toPath)

	err = Copy(fromPath, toPath, offset, limit)
	require.NotEmpty(t, err)
	require.Equal(t, resError, err.Error())
}

func TestCopy(t *testing.T) {
	t.Run("TC1 out_offset0_limit0", func(t *testing.T) {
		tc(t, "testdata/input.txt", "testdata/out_offset0_limit0.txt", 0, 0)
	})

	t.Run("TC2 out_offset0_limit10", func(t *testing.T) {
		tc(t, "testdata/input.txt", "testdata/out_offset0_limit10.txt", 0, 10)
	})

	t.Run("TC3 out_offset0_limit1000", func(t *testing.T) {
		tc(t, "testdata/input.txt", "testdata/out_offset0_limit1000.txt", 0, 1000)
	})

	t.Run("TC4 out_offset0_limit10000", func(t *testing.T) {
		tc(t, "testdata/input.txt", "testdata/out_offset0_limit10000.txt", 0, 10000)
	})

	t.Run("TC5 out_offset100_limit1000", func(t *testing.T) {
		tc(t, "testdata/input.txt", "testdata/out_offset100_limit1000.txt", 100, 1000)
	})

	t.Run("TC5 out_offset6000_limit1000", func(t *testing.T) {
		tc(t, "testdata/input.txt", "testdata/out_offset6000_limit1000.txt", 6000, 10000)
	})

	t.Run("TC6 tolstoy_voyna-i-mir.txt", func(t *testing.T) {
		tc(t, "testdata/tolstoy_voyna-i-mir.txt", "testdata/tolstoy_voyna-i-mir.txt", 0, 0)
	})

	t.Run("TC7 tolstoy_voyna-i-mir.txt -offset 5350525 -limit 1000 ", func(t *testing.T) {
		tc(t, "testdata/tolstoy_voyna-i-mir.txt", "testdata/empty_out.txt", 5350525, 1000)
	})

	t.Run("ETC1 no_exist_input ", func(t *testing.T) {
		etc(t, "testdata/no_exist_input.txt", "open testdata/no_exist_input.txt: no such file or directory", 0, 0)
	})

	t.Run("ETC2 unsupported file", func(t *testing.T) {
		etc(t, "/dev/urandom", "unsupported file", 0, 0)
	})

	t.Run("ETC3 offset exceeds file size", func(t *testing.T) {
		etc(t, "testdata/input.txt", "offset exceeds file size", 100000, 1000)
	})
}
