package is

func Nil[T comparable](v T) bool {
	var zero T
	return v == zero
}
