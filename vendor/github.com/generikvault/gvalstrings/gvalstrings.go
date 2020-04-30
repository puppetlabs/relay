package gvalstrings

// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"context"
	"fmt"
	"text/scanner"

	"github.com/PaesslerAG/gval"
)

// SingleQuoted returns a gval.Language for parsing single quoted strings.
func SingleQuoted() gval.Language {
	return gval.PrefixExtension(scanner.Char, parseSingleQuotedString)
}

func parseSingleQuotedString(c context.Context, p *gval.Parser) (gval.Evaluable, error) {
	s, err := UnquoteSingleQuoted(p.TokenText())
	if err != nil {
		return nil, fmt.Errorf("could not parse string: %v", err)
	}
	return p.Const(s), nil
}
