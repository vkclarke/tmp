package pair

type Interface[A, B any] interface {
	A() A
	B() B
}

type impl[A, B any] struct {
	a A
	b B
}

func (self impl[A, B]) A() A { return self.a }
func (self impl[A, B]) B() B { return self.b }

func Of[A, B any](a A, b B) Interface[A, B] {
	return impl[A, B]{a, b}
}

func OfSame[T any](a, b T) Interface[T, T] {
	return impl[T, T]{a, b}
}
