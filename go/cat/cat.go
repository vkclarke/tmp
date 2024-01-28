package main

import (
	"io"
	"os"
	"unicode/utf8"
)

func main() {
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

	in := func() <-chan rune {
		in := make(chan rune)
		go func() {
			buf := make([]byte, 4)
			for {
				n, err := os.Stdin.Read(buf)
				if err == io.EOF {
					break
				}
				for _, char := range string(buf[:n]) {
					in <- char
				}
			}
			close(in)
		}()
		return in
	}()

	// main loop
	for char := range in {
		out <- char
	}
	close(out)
	<-wait
}
