// usage: <prog> [target]
// (will read from stdin if target is absent)
package main

import (
	"io"
	"log"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	args := os.Args[1:]
	stderr := log.New(os.Stderr, "[errors] ", 0)

	txt := func() string {
		if len(args) < 1 {
			input, err := io.ReadAll(os.Stdin)
			if err != nil {
				stderr.Fatal(err)
			}
			return string(input)
		}
		contents, err := os.ReadFile(args[0])
		if err != nil {
			stderr.Fatal(err)
		}
		return string(contents)
	}()

	wait := make(chan struct{})
	stdout := func() chan<- rune {
		stdout := make(chan rune)
		go func() {
			buf := make([]byte, 4)
			for char := range stdout {
				n := utf8.EncodeRune(buf, char)
				os.Stdout.Write(buf[:n])
			}
			close(wait)
		}()
		return stdout
	}()

	for char := range shrinkwrap(iterString(txt)) {
		stdout <- char
	}

	close(stdout)
	<-wait
}

func shrinkwrap(in <-chan rune) <-chan rune {
	out := make(chan rune)
	go func() {
		var quote rune
		for char := range in {
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
			out <- char
		}
		close(out)
	}()
	return out
}

func iterString(txt string) <-chan rune {
	chars := make(chan rune)
	go func() {
		for _, char := range txt {
			chars <- char
		}
		close(chars)
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
