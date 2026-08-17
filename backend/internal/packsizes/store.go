package packsizes

import (
	"context"
	"errors"
)

var ErrAlreadyExists = errors.New("pack size already exists")

type Store interface {
	List(context.Context) ([]int, error)
	Add(context.Context, int) error
}
