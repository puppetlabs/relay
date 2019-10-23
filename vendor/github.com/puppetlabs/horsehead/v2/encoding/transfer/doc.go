/*
Package transfer provides an interface for encoding and decoding values
for storage. The utility in this package is transparent to the user
and it is used to maintain byte integrity on values used in workflows.

Basic use when encoding a value:
	encoder := transfer.Encoders[transfer.DefaultEncodingType]()

	result, err := encoder.EncodeForTransfer([]byte("super secret token"))
	if err != nil {
		// handle error
	}

Basic use when decoding a value:
	encodingType, value := transfer.ParseEncodedValue("base64:c3VwZXIgc2VjcmV0IHRva2Vu")
	encoder := transfer.Encoders[encoderType]()

	result, err := encoder.DecodeFromTransfer(value)
	if err != nil {
		// handle error
	}

In addition, it provides types for automatically encoding and decoding non-UTF-8 strings in an expanded JSON
format.
*/

package transfer
