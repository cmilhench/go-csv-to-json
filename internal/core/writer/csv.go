package writer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/cm-dev/template/internal/ports"
)

var _ ports.Writer[[]string] = (*CsvWriter)(nil)

type CsvWriter struct{}

func (t *CsvWriter) Write(ctx context.Context, writer io.Writer, records <-chan []string) error {
	if ctx == nil || writer == nil || records == nil {
		return fmt.Errorf("invalid parameters")
	}

	csvWriter := csv.NewWriter(writer)
	csvWriter.Comma = ','
	csvWriter.UseCRLF = true
	for {
		select {
		case <-ctx.Done():
			csvWriter.Flush()
			return nil
		case record, ok := <-records:
			if !ok {
				csvWriter.Flush()
				return nil
			}
			if e := csvWriter.Write(record); e != nil {
				csvWriter.Flush()
				return fmt.Errorf("failed to write record: %w", e)
			}
		}
	}
}
