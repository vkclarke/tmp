package iter

func Each[V any](in <-chan V, action func(V) V) <-chan V {
	out := make(chan V)
	go func() {
		defer close(out)
		for value := range in {
			out <- action(value)
		}
	}()
	return out
}
