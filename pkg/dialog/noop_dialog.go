package dialog

import "io"

type noopTable struct{}

func (t noopTable) Headers([]string) Table   { return t }
func (t noopTable) Rows([][]string) Table    { return t }
func (t noopTable) AppendRow([]string) Table { return t }
func (t noopTable) Flush() error             { return nil }

type noopDialog struct{}

func (n noopDialog) WithStdout(io.Writer) Dialog    { return n }
func (n noopDialog) WithStderr(io.Writer) Dialog    { return n }
func (noopDialog) Progress(string)                  {}
func (noopDialog) Progressf(string, ...interface{}) {}
func (noopDialog) Info(string)                      {}
func (noopDialog) Infof(string, ...interface{})     {}
func (noopDialog) Error(string)                     {}
func (noopDialog) Errorf(string, ...interface{})    {}
func (noopDialog) WriteString(string) error         { return nil }
func (noopDialog) Table() Table                     { return noopTable{} }

func newNoopDialog() Dialog {
	return noopDialog{}
}
