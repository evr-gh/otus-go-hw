package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	paramIsOK := false
	switch {
	case from == "":
		fmt.Println("from file isn't set")
		fallthrough
	case to == "":
		fmt.Println("out file isn't set")
	case from == to:
		fmt.Println("from file and out file are the same")
	default:
		paramIsOK = true
	}

	if !paramIsOK {
		return
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Printf("Fail to copy %s to %s: %s\n", from, to, err)
	}
}
