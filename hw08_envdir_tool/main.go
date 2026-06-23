package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Not all arguments are set.\nUsage: go-envdir <path to env dir> <command> (<argbment>)*")
		os.Exit(-1)
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(-1)
	}
	rc := RunCmd(os.Args[2:], env)

	os.Exit(rc)
}
