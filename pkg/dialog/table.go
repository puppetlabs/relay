package dialog

type Table interface {
	Headers([]string) Table
	Rows([][]string) Table
	AppendRow([]string) Table
	Flush() error
}
