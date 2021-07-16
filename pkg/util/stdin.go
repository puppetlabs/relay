package util

import (
	"bytes"
	"io"
	"os"
)

func PassedStdin() (bool, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}

	if (info.Mode() & os.ModeCharDevice) == 0 {
		return true, nil
	}

	return false, nil
}

func ReadStdin(readLimit int64) ([]byte, error) {
	buf := bytes.Buffer{}
	reader := &io.LimitedReader{R: os.Stdin, N: readLimit}

	n, err := buf.ReadFrom(reader)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	return buf.Bytes(), nil
}
