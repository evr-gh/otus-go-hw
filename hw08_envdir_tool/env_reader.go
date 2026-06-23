package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	cleanPath := filepath.Clean(dir)
	cleanPath = strings.TrimRight(cleanPath, "/\\")

	dirCont, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)
	for _, f := range dirCont {
		if f.Type().IsRegular() {
			fn := filepath.Clean(filepath.Join(cleanPath, f.Name()))

			cont, err := os.ReadFile(fn)
			if cont == nil {
				var v EnvValue
				v.NeedRemove = true
				env[f.Name()] = v
			} else {
				lines := bytes.Split(cont, []byte{'\n'})
				line := lines[0]
				value := bytes.ReplaceAll(line, []byte{0x00}, []byte{'\n'})
				if err == nil {
					var v EnvValue
					v.NeedRemove = false
					v.Value = strings.TrimRight(string(value), " \t")
					env[f.Name()] = v
				}
			}
		}
	}

	return env, nil
}
