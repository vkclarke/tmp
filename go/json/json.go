// usage: <prog> [target]
// (will read from stdin if target is absent)
package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"
)

func main() {
	name := filepath.Base(os.Args[0])
	args := os.Args[1:]
	errlog := log.New(os.Stderr, name+": ", 0)

	if len(args) < 1 {
		errlog.Fatal("no arguments")
	}

	var cmd func(<-chan rune) <-chan rune
	switch args[0] {
	case "shrink":
		cmd = shrinkwrap
	case "expand":
		cmd = expand
	default:
		errlog.Fatalf("invalid command: %q", args[0])
	}

	in := func() <-chan rune {
		in := make(chan rune)
		go func() {
			b := func() *bufio.Reader {
				if len(args) < 2 {
					return bufio.NewReader(os.Stdin)
				}
				file, _ := os.Open(args[1])
				return bufio.NewReader(file)
			}()
			for {
				char, _, err := b.ReadRune()
				if err != nil {
					break
				}
				in <- char
			}
			close(in)
		}()
		return in
	}()

	wait := make(chan struct{})
	out := func() chan<- rune {
		out := make(chan rune)
		go func() {
			buf := make([]byte, 4)
			for char := range out {
				n := utf8.EncodeRune(buf, char)
				os.Stdout.Write(buf[:n])
			}
			close(wait)
		}()
		return out
	}()

	// Main loop
	for char := range cmd(in) {
		out <- char
	}
	close(out)
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

func expand(in <-chan rune) <-chan rune {
	out := make(chan rune)
	go func() {
		var tabs int
		indent := func() {
			for i := 0; i < tabs; i++ {
				out <- '\t'
			}
		}
		for {
			char := <-in
			switch char {
			case 0:
				out <- '\n'
				close(out)
				return
			case '{', '[':
				out <- char
				out <- '\n'
				tabs += 1
				indent()
			case '}', ']':
				tabs -= 1
				out <- '\n'
				indent()
				out <- char
			case ',':
				out <- char
				out <- '\n'
				indent()
			case ':':
				out <- char
				out <- ' '
			default:
				out <- char
			}
		}
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
