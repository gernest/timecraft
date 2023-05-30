// Package stream is a library of generic types designed to work on streams of
// values.
package stream

import "io"

// Reader is an interface implemented by types that read a stream of values of
// type T.
type Reader[T any] interface {
	// Reads values from the stream, returning the number of values read and any
	// error that occurred.
	//
	// The error is io.EOF when the end of the stream has been reached.
	Read(values []T) (int, error)
}

// NewReader constructs a Reader from a sequence of values.
func NewReader[T any](values ...T) Reader[T] {
	return &reader[T]{values: append([]T{}, values...)}
}

type reader[T any] struct{ values []T }

func (r *reader[T]) Read(values []T) (n int, err error) {
	n = copy(values, r.values)
	r.values = r.values[n:]
	if len(r.values) == 0 {
		err = io.EOF
	}
	return n, err
}

// ReadCloser represents a closable stream of values of T.
//
// ReadClosers is like io.ReadCloser for values of any type.
type ReadCloser[T any] interface {
	Reader[T]
	io.Closer
}

// NewReadCloser constructs a ReadCloser from the pair of r and c.
func NewReadCloser[T any](r Reader[T], c io.Closer) ReadCloser[T] {
	return &readCloser[T]{reader: r, closer: c}
}

type readCloser[T any] struct {
	reader Reader[T]
	closer io.Closer
}

func (r *readCloser[T]) Close() error                 { return r.closer.Close() }
func (r *readCloser[T]) Read(values []T) (int, error) { return r.reader.Read(values) }

// NopCloser constructs a ReadCloser from a Reader.
func NopCloser[T any](r Reader[T]) ReadCloser[T] {
	return &nopCloser[T]{reader: r}
}

type nopCloser[T any] struct{ reader Reader[T] }

func (r *nopCloser[T]) Close() error                 { return nil }
func (r *nopCloser[T]) Read(values []T) (int, error) { return r.reader.Read(values) }

// ReadAll reads all values from r and returns them as a slice, along with any
// error that occurred (other than io.EOF).
func ReadAll[T any](r Reader[T]) ([]T, error) {
	values := make([]T, 0, 1)
	for {
		if len(values) == cap(values) {
			values = append(values, make([]T, 2*len(values))...)[:len(values)]
		}
		n, err := r.Read(values[len(values):cap(values)])
		values = values[:len(values)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return values, err
		}
	}
}

// Writer is an interface implemented by types that write a stream of values of
// type T.
type Writer[T any] interface {
	Write(values []T) (int, error)
}

// WriteCloser represents a closable stream of values of T.
//
// WriteClosers is like io.WriteCloser for values of any type.
type WriteCloser[T any] interface {
	Writer[T]
	io.Closer
}

func NewWriteCloser[T any](w Writer[T], c io.Closer) WriteCloser[T] {
	return &writeCloser[T]{writer: w, closer: c}
}

type writeCloser[T any] struct {
	writer Writer[T]
	closer io.Closer
}

func (w *writeCloser[T]) Write(values []T) (int, error) {
	return w.writer.Write(values)
}

func (w *writeCloser[T]) Close() error {
	return w.closer.Close()
}

// Copy writes values read from r to w, returning the number of values written
// and any error other than io.EOF.
func Copy[T any](w Writer[T], r Reader[T]) (int64, error) {
	b := make([]T, 20)
	n := int64(0)

	for {
		rn, err := r.Read(b)

		if rn > 0 {
			wn, err := w.Write(b[:rn])
			n += int64(wn)
			if err != nil {
				return n, err
			}
			if wn < rn {
				return n, io.ErrNoProgress
			}
		}

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return n, err
		}
	}
}