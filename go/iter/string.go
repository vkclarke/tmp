package iter

func String(s string) <-chan rune {
	out := make(chan rune)
	go func() {
		defer close(out)
		for _, char := range s {
			out <- char
		}
	}()
	return out
}
