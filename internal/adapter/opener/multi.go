package opener

import (
	"errors"
	"io"
)

type multiReadCloser struct {
	readers []io.ReadCloser
}

func (mr *multiReadCloser) Close() error {
	var errs []error
	for _, v := range mr.readers {
		if err := v.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (mr *multiReadCloser) Read(p []byte) (n int, err error) {
	for len(mr.readers) > 0 {
		n, err = mr.readers[0].Read(p)
		if err == io.EOF {
			// Use eofReader instead of nil to avoid nil panic
			// after performing flatten (Issue 18232).
			mr.readers[0] = eofReader{} // permit earlier GC
			mr.readers = mr.readers[1:]
		}
		if n > 0 || err != io.EOF {
			if err == io.EOF && len(mr.readers) > 0 {
				// Don't return EOF yet. More readers remain.
				err = nil
			}
			return
		}
	}
	return 0, io.EOF
}

type eofReader struct{}

func (eofReader) Read([]byte) (int, error) {
	return 0, io.EOF
}
func (eofReader) Close() error {
	return nil
}

func MultiReadCloser(input ...io.ReadCloser) io.ReadCloser {
	c := make([]io.ReadCloser, len(input))
	copy(c, input)
	return &multiReadCloser{c}
}
