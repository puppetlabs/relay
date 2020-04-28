package gvalstrings

// Some function are under Copyright 2009 The Go Authors.

import (
	"strconv"
	"unicode/utf8"
)

func unhex(b byte) (v rune, ok bool) {

	c := rune(b)

	switch {

	case '0' <= c && c <= '9':
		return c - '0', true

	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true

	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true

	}

	return
}

// UnquoteChar decodes the first character or byte in the escaped string
// or character literal represented by the string s.
// It returns four values:
//
//	1) value, the decoded Unicode code point or byte value;
//	2) multibyte, a boolean indicating whether the decoded character requires a multibyte UTF-8 representation;
//	3) tail, the remainder of the string after the character; and
//	4) an error that will be nil if the character is syntactically valid.
//  It permits the sequence \' and disallows unescaped '.
func unquoteChar(s string) (value rune, multibyte bool, tail string, err error) {

	// easy cases
	if len(s) == 0 {
		err = strconv.ErrSyntax
		return
	}

	switch c := s[0]; {

	case c == '\'':
		err = strconv.ErrSyntax
		return

	case c >= utf8.RuneSelf:
		r, size := utf8.DecodeRuneInString(s)
		return r, true, s[size:], nil

	case c != '\\':
		return rune(s[0]), false, s[1:], nil
	}

	// hard case: c is backslash

	if len(s) <= 1 {
		err = strconv.ErrSyntax
		return
	}

	c := s[1]
	s = s[2:]

	switch c {

	case 'a':
		value = '\a'

	case 'b':
		value = '\b'

	case 'f':
		value = '\f'

	case 'n':
		value = '\n'

	case 'r':
		value = '\r'

	case 't':
		value = '\t'

	case 'v':
		value = '\v'

	case 'x', 'u', 'U':
		n := 0

		switch c {

		case 'x':
			n = 2

		case 'u':
			n = 4

		case 'U':
			n = 8

		}

		var v rune

		if len(s) < n {
			err = strconv.ErrSyntax
			return

		}

		for j := 0; j < n; j++ {

			x, ok := unhex(s[j])

			if !ok {
				err = strconv.ErrSyntax
				return
			}

			v = v<<4 | x
		}

		s = s[n:]

		if c == 'x' {
			// single-byte string, possibly not UTF-8
			value = v
			break
		}

		if v > utf8.MaxRune {
			err = strconv.ErrSyntax
			return
		}

		value = v
		multibyte = true

	case '0', '1', '2', '3', '4', '5', '6', '7':

		v := rune(c) - '0'

		if len(s) < 2 {
			err = strconv.ErrSyntax
			return
		}

		for j := 0; j < 2; j++ { // one digit already; two more
			x := rune(s[j]) - '0'
			if x < 0 || x > 7 {
				err = strconv.ErrSyntax
				return
			}
			v = (v << 3) | x
		}

		s = s[2:]

		if v > 255 {
			err = strconv.ErrSyntax
			return
		}

		value = v

	case '\\':
		value = '\\'

	case '\'':
		value = '\''

	default:
		err = strconv.ErrSyntax
		return

	}

	tail = s

	return

}

// UnquoteSingleQuoted interprets s as a single-quoted
// Go like string literal, returning the string value
// that s quotes.
func UnquoteSingleQuoted(s string) (string, error) {

	n := len(s)

	if n < 2 {
		return "", strconv.ErrSyntax
	}
	quote := s[0]

	if quote != '\'' {
		return strconv.Unquote(s)
	}

	if quote != s[n-1] {
		return "", strconv.ErrSyntax
	}

	s = s[1 : n-1]

	if contains(s, '\n') {
		return "", strconv.ErrSyntax
	}

	// Is it trivial? Avoid allocation.

	if !contains(s, '\\') && !contains(s, quote) {
		if utf8.ValidString(s) {
			return s, nil
		}
	}

	var runeTmp [utf8.UTFMax]byte

	buf := make([]byte, 0, 3*len(s)/2) // Try to avoid more allocations.

	for len(s) > 0 {

		c, multibyte, ss, err := unquoteChar(s)

		if err != nil {
			return "", err
		}

		s = ss

		if c < utf8.RuneSelf || !multibyte {
			buf = append(buf, byte(c))

		} else {
			n := utf8.EncodeRune(runeTmp[:], c)
			buf = append(buf, runeTmp[:n]...)

		}

	}
	return string(buf), nil

}

// contains reports whether the string contains the byte c.

func contains(s string, c byte) bool {

	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}

	return false

}
