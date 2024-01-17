package is

import (
	constr "golang.org/x/exp/constraints"
)

func Odd[T constr.Integer](n T) bool {
	return n%2 != 0
}