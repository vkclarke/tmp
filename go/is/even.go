package is

import (
	constr "golang.org/x/exp/constraints"
)

func Even[T constr.Integer](n T) bool {
	return n%2 == 0
}