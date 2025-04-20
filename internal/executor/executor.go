package executor

import "io"

type Executor interface {
	io.Closer
	Do() error
}
