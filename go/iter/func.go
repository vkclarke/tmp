package iter

func Func[V any](next func() (V, bool)) <-chan V {
	out := make(chan V)
	go func() {
		defer close(out)
		for v, ok := next(); ok; v, ok = next() {
			out <- v
		}
	}()
	return out
}
