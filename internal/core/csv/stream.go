package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
)

func Read(ctx context.Context, reader io.Reader) (<-chan []string, <-chan error) {
	out := make(chan []string, 1)
	err := make(chan error, 1)
	if ctx == nil || reader == nil {
		err <- fmt.Errorf("invalid parameters")
		close(out)
		close(err)
		return out, err
	}

	go func() {
		defer close(out)
		defer close(err)

		csvReader := csv.NewReader(reader)
		csvReader.ReuseRecord = false
		csvReader.FieldsPerRecord = -1
		csvReader.TrimLeadingSpace = true

		// Read the headers
		headers, e := csvReader.Read()
		if e != nil {
			err <- fmt.Errorf("failed to read CSV headers: %w", e)
			return
		}
		select {
		case out <- headers:
		case <-ctx.Done():
			err <- ctx.Err()
			return
		}

		// Read the data
		for {
			record, e := csvReader.Read()
			if e != nil {
				if e == io.EOF {
					return
				}
				select {
				case err <- fmt.Errorf("failed to read record: %w", e):
				case <-ctx.Done():
					err <- ctx.Err()
				}
				continue
			}
			select {
			case out <- record:
			case <-ctx.Done():
				err <- ctx.Err()
				return
			}
		}
	}()
	return out, err
}
