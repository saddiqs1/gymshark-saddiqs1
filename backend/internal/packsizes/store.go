package packsizes

import (
	"context"
	"errors"
)

var (
	ErrAlreadyExists = errors.New("pack size already exists")
	ErrNotFound      = errors.New("pack size not found")
)

type Store interface {
	List(context.Context) ([]int, error)
	Add(context.Context, int) error
	Remove(context.Context, int) error
}
