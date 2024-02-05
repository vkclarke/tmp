package iter

func While[V any](in <-chan V, cond func() bool, action func(V) V) <-chan V {
	out := make(chan V)
	go func() {
		defer close(out)
		for value := range in {
			if !cond() {
				break
			}
			out <- action(value)
		}
	}()
	return out
}
