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
		ok = matchDigits(line)
	} else if pattern == "\\w" {
		ok = matchAlphanumeric(line)
	} else if isPositiveCharGroups(pattern) {
		ok = matchLiteralChar(line, pattern[1:len(pattern)-1])
	} else if isNegativeCharGroups(pattern) {
		ok = !matchLiteralChar(line, pattern[2:len(pattern)-1])
	} else {
		ok = matchLiteralChar(line, pattern)
	}

	return ok, nil
}

func matchLiteralChar(line []byte, pattern string) bool {
	return bytes.ContainsAny(line, pattern)
}

func matchDigits(line []byte) bool {
	return bytes.ContainsAny(line, "0123456789")
}

func matchAlphanumeric(line []byte) bool {
	return bytes.ContainsAny(line, "abcdefghijklmnopqrstvuwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
}

func isPositiveCharGroups(pattern string) bool {
	return len(pattern) >= 3 &&
		pattern[0] == '[' &&
		pattern[1] != '^' &&
		pattern[len(pattern)-1] == ']'
}

func isNegativeCharGroups(pattern string) bool {
	return len(pattern) >= 4 &&
		pattern[0] == '[' &&
		pattern[1] == '^' &&
		pattern[len(pattern)-1] == ']'
}
