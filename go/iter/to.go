package iter

func To[V any](in <-chan V, out chan<- V) {
	for value := range in {
		out <- value
	}
}
