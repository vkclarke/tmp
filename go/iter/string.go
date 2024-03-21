package iter

import "unicode/utf8"

func String(txt string) <-chan rune {
	out := make(chan rune, utf8.RuneCountInString(txt))
	go func() {
		defer close(out)
		for _, char := range txt {
			out <- char
		}
	}()
	return out
}

func StringBytes(txt string) <-chan byte {
	n := len(txt)
	out := make(chan byte, n)
	go func() {
		defer close(out)
		for i := range n {
			out <- txt[i]
		}
	}()
	return out
}
