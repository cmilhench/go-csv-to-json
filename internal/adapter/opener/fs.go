package opener

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cm-dev/template/internal/ports"
)

type OpenerFS struct{}

var _ ports.Opener = (*OpenerFS)(nil)

func (o *OpenerFS) Open(ctx context.Context, name string) (io.ReadCloser, error) {
	var readers []io.ReadCloser

	entries, err := os.ReadDir(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	for _, entry := range entries {
		filePath := filepath.Join(name, entry.Name())
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %v", err)
		}
		readers = append(readers, f)
	}
	return MultiReadCloser(readers...), nil
}
