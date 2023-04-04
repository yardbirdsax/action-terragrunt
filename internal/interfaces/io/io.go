//go:generate mockgen -source=io.go -destination=../../mock/io/io.go -package=io
package interfaces

import (
	"io"
)

type ReadCloser interface {
	io.ReadCloser
}

type Writer interface {
	io.Writer
}
