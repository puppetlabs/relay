package config

import "io"

type IO struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}
