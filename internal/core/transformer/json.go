package transformer

import (
	"context"

	"github.com/cm-dev/template/internal/core/json"
	"github.com/cm-dev/template/internal/ports"
)

var _ ports.Transformer[any] = (*JsonTransformer)(nil)

type JsonTransformer struct{}

func (t *JsonTransformer) Transform(ctx context.Context, input <-chan []string) (<-chan any, <-chan error) {
	out := make(chan any)
	err := make(chan error)
	go func() {
		defer close(out)
		defer close(err)
		var head []string
		select {
		case <-ctx.Done():
			return
		case head = <-input:
		}
		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-input:
				if !ok {
					return
				}
				o, e := json.Inflate(head, line)
				if e != nil {
					err <- e
				} else {
					out <- o
				}
			}
		}
	}()

	return out, err
}
