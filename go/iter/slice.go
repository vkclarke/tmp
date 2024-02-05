package iter

func Slice[V any](values []V) <-chan V {
	out := make(chan V)
	go func() {
		defer close(out)
		for i := range values {
			out <- values[i]
		}
	}()
	return out
}
