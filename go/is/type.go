package is

func Type[T any](v any) (ok bool) {
	_, ok = v.(T)
	return
}
