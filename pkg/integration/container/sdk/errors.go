package sdk

import "errors"

var (
	ErrNotSupported = errors.New("sdk: not supported")
	ErrNotFound     = errors.New("sdk: not found")
)
