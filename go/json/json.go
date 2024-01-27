// usage: <prog> [target]
// (will read from stdin if target is absent)
package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"unicode"
)

func main() {
	args := os.Args[1:]
	stderr := log.New(os.Stderr, "[errors] ", 0)

	input := bytes.NewBuffer(make([]byte, 0, os.Getpagesize()))
	input.ReadFrom(func() *os.File {
		if len(args) < 1 {
			return os.Stdin
		}
		f, err := os.Open(args[0])
		if err != nil {
			stderr.Fatalln(err)
		}
		return f
	}())

	os.Stdout.WriteString(shrinkwrap(input.String()))
}

func shrinkwrap(src string) string {
	var wrapped strings.Builder
	var quote rune
	for _, char := range src {
		switch {
		case isSpace(char) && quote == 0:
			continue
		case char == '"':
			switch {
			case quote == char:
				quote = 0
			case quote == 0:
				quote = char
			}
		}
		wrapped.WriteRune(char)
	}
	return wrapped.String()
}

func iterString(txt string, size int) <-chan rune {
	chars := make(chan rune, size)
	go func() {
		defer close(chars)
		for _, char := range txt {
			chars <- char
		}
	}()
	return chars
}

func isSpace(char rune) bool  { return unicode.IsSpace(char) }
func isNumber(char rune) bool { return unicode.IsNumber(char) }
func isLetter(char rune) bool { return unicode.IsLetter(char) }

var bounds = map[rune]rune{
	'(':  ')',
	'[':  ']',
	'{':  '}',
	'<':  '>',
	'\'': '\'',
	'"':  '"',
	'`':  '`',
}
