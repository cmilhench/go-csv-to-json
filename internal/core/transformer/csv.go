package transformer

import (
	"context"
	"strings"

	"github.com/cm-dev/template/internal/ports"
)

var _ ports.Transformer[[]string] = (*CsvTransformer)(nil)

type CsvTransformer struct{}

func (t *CsvTransformer) Transform(ctx context.Context, input <-chan []string) (<-chan []string, <-chan error) {
	out := make(chan []string)
	err := make(chan error)
	go func() {
		defer close(out)
		defer close(err)
		var head []string
		select {
		case <-ctx.Done():
			return
		case head = <-input:
			for i := 0; i < len(head); i++ {
				s := head[i]
				s = strings.ToUpper(string(s[0])) + s[1:] // capitalise
				s = strings.ReplaceAll(s, ".", " ")       // replace "."s with " "s
				head[i] = s
			}
			out <- head
		}
		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-input:
				if !ok {
					return
				}
				out <- line
			}
		}
	}()

	return out, err
}
