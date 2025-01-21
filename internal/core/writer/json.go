package writer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

func WriteJson(ctx context.Context, writer io.Writer, records <-chan any) error {
	if ctx == nil || writer == nil || records == nil {
		return fmt.Errorf("invalid parameters")
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case record, ok := <-records:
			if !ok {
				return nil
			}
			out, err := json.MarshalIndent(record, "", "  ")
			if err != nil {
				log.Fatalf("Error marshaling record: %v", err)
			}
			if _, e := writer.Write(out); e != nil {
				return fmt.Errorf("failed to write record: %w", e)
			}
		}
	}
}
