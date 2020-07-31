package dev

import (
	"bufio"
	"bytes"
	"io"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/streaming"
)

// multiDocDecoder basically does the same thing as yaml.NewDocumentDecoder,
// but without a super small buffer that blows up all the time. I know they are
// looking to fix this in some future version, but for now, we are on 1.17.2, which is
// completely broken for _most_ manifests I tried. I don't even know how kubernetes deals
// with this...
//
// TODO: remove when k8s apis are updated.
type multiDocDecoder struct {
	r         io.ReadCloser
	scanner   *bufio.Scanner
	remaining []byte
}

func (d *multiDocDecoder) Read(data []byte) (int, error) {
	left := len(d.remaining)
	if left == 0 {
		if !d.scanner.Scan() {
			err := d.scanner.Err()
			if err == nil {
				err = io.EOF
			}

			return 0, err
		}
		out := d.scanner.Bytes()
		d.remaining = out
		left = len(out)
	}

	if left <= len(data) {
		copy(data, d.remaining)
		d.remaining = nil

		return left, nil
	}

	copy(data, d.remaining[:len(data)])
	d.remaining = d.remaining[len(data):]

	return len(data), io.ErrShortBuffer
}

func (d *multiDocDecoder) Close() error {
	return d.r.Close()
}

func newMultiDocDecoder(r io.ReadCloser) io.ReadCloser {
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 4*1024)
	scanner.Buffer(buf, 5*1024*1024)
	scanner.Split(splitYAMLDocument)

	return &multiDocDecoder{
		r:       r,
		scanner: scanner,
	}
}

// TODO used to support multiDocDecoder: remove when k8s apis are updated.
func splitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var yamlSeparator = "---"

	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	sep := len([]byte(yamlSeparator))
	if i := bytes.Index(data, []byte(yamlSeparator)); i >= 0 {
		// We have a potential document terminator
		i += sep
		after := data[i:]
		if len(after) == 0 {
			// we can't read any more characters
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}

			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + j + 1, data[0 : i-sep], nil
		}

		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}

	// Request more data.
	return 0, nil, nil
}

// relay-core provides a testing utility that does this already. I would like
// to eventually pull that out into another package and use it here so we aren't
// duplicating as much functionality.
func parseManifest(r io.ReadCloser, patchers ...objectPatcherFunc) ([]runtime.Object, error) {
	defer r.Close()

	yamldecoder := newMultiDocDecoder(r)
	ud := serializer.NewCodecFactory(scheme).UniversalDeserializer()
	sd := streaming.NewDecoder(yamldecoder, ud)

	var objs []runtime.Object

	for {
		obj, _, err := sd.Decode(nil, nil)
		if err != nil {
			if err == io.EOF {
				break
			}

			// only return an error if it's not a missing Kind error.
			// These seem to only stem from comments being parsed as documents.
			if !runtime.IsMissingKind(err) {
				return nil, err
			}

			continue
		}

		objs = append(objs, obj)
	}

	return objs, nil
}
