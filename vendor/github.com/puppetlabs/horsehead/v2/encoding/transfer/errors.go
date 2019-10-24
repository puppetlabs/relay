package transfer

import "errors"

var (
	ErrUnknownEncodingType = errors.New("transfer: unknown encoding type")
	ErrNotEncodable        = errors.New("transfer: the given value cannot be encoded for this transfer mechanism")
)
