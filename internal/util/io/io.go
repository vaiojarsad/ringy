package io

import (
	"io"

	"github.com/vaiojarsad/ringy/internal/appcontext"
)

func CloseWithHandler(c io.Closer, context string, handler func(string, error)) {
	if err := c.Close(); err != nil {
		handler(context, err)
	}
}

func Close(c io.Closer, context string) {
	CloseWithHandler(c, context, func(context string, err error) {
		appcontext.GetInstance().ErrLogger.Println("error closing.", "appcontext", context, "error", err)
	})
}
