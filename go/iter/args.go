package iter

func Args[V any](args ...V) <-chan V {
	n := len(args)
	out := make(chan V, n)
	go func() {
		defer close(out)
		for i := range n {
			out <- args[i]
		}
	}()
	return out
}
