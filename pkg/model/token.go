package model

import "fmt"

type Token string

func (t *Token) Bearer() string {
	return fmt.Sprintf("Bearer %s", t)
}

func (t Token) String() string {
	return string(t)
}
