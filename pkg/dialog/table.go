package dialog

import "io"

type Table interface {
	Headers([]string) Table
	Rows([][]string) Table
	AppendRow([]string) Table
	WriteTo(io.Writer) error
}
