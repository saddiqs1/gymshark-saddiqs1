package packsizes

import "context"

type Store interface {
	List(context.Context) ([]int, error)
}
