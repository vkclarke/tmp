package iter

func Slice[V any](slice []V) <-chan V {
	n := len(slice)
	out := make(chan V, n)
	go func() {
		defer close(out)
		for i := range n {
			out <- slice[i]
		}
	}()
	return out
}
