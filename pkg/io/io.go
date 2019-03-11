package io

import stdio "io"

type IO struct {
	In     stdio.Reader
	Out    stdio.Writer
	ErrOut stdio.Writer
}
