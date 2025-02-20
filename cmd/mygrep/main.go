package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	var ok bool

	if pattern == "\\d" {
		ok = bytes.ContainsAny(line, "0123456789")
	} else if pattern == "\\w" {
		ok = bytes.ContainsAny(line, "abcdefghijklmnopqrstvuwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
	} else if len(line) > 0 && line[0] == '[' && line[len(line)-1] == ']' {
		ok = bytes.ContainsAny(line, pattern[1:len(pattern)-1])
	} else {
		ok = bytes.ContainsAny(line, pattern)
	}

	return ok, nil
}
