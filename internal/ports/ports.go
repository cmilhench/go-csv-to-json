package ports

import (
	"context"
	"io"
)

type Transformer[T any] interface {
	Transform(context.Context, <-chan []string) (<-chan T, <-chan error)
}

type Writer[T any] interface {
	Write(context.Context, io.Writer, <-chan T) error
}
