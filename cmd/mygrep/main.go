package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

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

	ok, err := matchLine(string(line), pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}

func matchLine(line string, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) == 0 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	if pattern[0] == '^' {
		return matchPattern(line, pattern[1:], 0), nil
	} else if pattern[len(pattern)-1] == '$' {
		return matchPattern(line, pattern[:len(pattern)-1], len(line)-len(pattern)+1), nil
	}

	for pos := range len(line) {
		if matchPattern(string(line), pattern, pos) {
			fmt.Println("Matched!")
			return true, nil
		}
	}

	fmt.Println("Not Matched!")
	return false, nil
}

func matchPattern(line string, pattern string, pos int) bool {
	n, m := len(pattern), len(line)
	lI := pos
	for pI := 0; pI < n; pI++ {
		if lI >= m {
			if pI+1 < n && pattern[pI+1] == '?' {
				continue
			} else {
				return false
			}
		}
		if pattern[pI] == '\\' && pI+1 < n {
			if pattern[pI+1] == 'd' && !isDigit(rune(line[lI])) {
				return false
			} else if pattern[pI+1] == 'w' && !isAlphanumeric(rune(line[lI])) {
				return false
			} else {
				pI++
			}
		} else if pattern[pI] == '[' && pI+1 < n {
			cp := strings.Index(pattern[pI:], "]") + pI
			if cp == pI-1 {
				return false
			}
			if pattern[pI+1] == '^' {
				if isMatchAnyPattern(pattern[pI+2:cp], string(line[lI])) {
					return false
				}
			} else {
				if !isMatchAnyPattern(pattern[pI+1:cp], string(line[lI])) {
					return false
				}
			}
			pI = cp
		} else if pattern[pI] == '+' {
			if pI == 0 {
				return false
			}
			if pattern[pI-1] == '.' {
				if pI == n-1 {
					return true
				} else {
					for lI < m && pattern[pI+1] != line[lI] {
						lI++
					}
					lI--
				}
			} else {
				for lI < m && pattern[pI-1] == line[lI] {
					lI++
				}
				lI--
			}
		} else if pattern[pI] == '?' {
			continue
		} else if pattern[pI] != line[lI] && pattern[pI] != '.' {
			if pI+1 < n && pattern[pI+1] == '?' {
				lI--
			} else {
				return false
			}
		}
		lI++
	}
	return true
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func isMatchAnyPattern(pattern string, text string) bool {
	return strings.Contains(pattern, text)
}
