package ports

import (
	"context"
	"io"
)

type Opener interface {
	Open(ctx context.Context, name string) (io.ReadCloser, error)
}
