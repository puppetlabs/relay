package defwalker

import (
	"fmt"
	"strings"
)

type WalkError struct {
	Path  string
	Cause error
}

func (e *WalkError) Error() string {
	return fmt.Sprintf("defwalker: failed to load %q: %+v", e.Path, e.Cause)
}

type WalkErrors []*WalkError

func (e WalkErrors) Error() string {
	wes := make([]string, len(e))
	for i, we := range e {
		wes[i] = strings.ReplaceAll(we.Error(), "\n", "\n  ")
	}

	return fmt.Sprintf("defwalker: failed to walk:\n* %s", strings.Join(wes, "\n * "))
}
