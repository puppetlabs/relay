package dialog

import (
	"io"

	"github.com/jedib0t/go-pretty/table"
)

type textTable struct {
	w       io.Writer
	headers []string
	rows    [][]string
}

func stringsToRow(arr []string) table.Row {
	// TODO: Is there a more idiomatic way of converting from []string to []interface?
	ifaces := make([]interface{}, 0, len(arr))

	for _, str := range arr {
		ifaces = append(ifaces, str)
	}

	return table.Row(ifaces)
}

func (t *textTable) Headers(h []string) Table {
	t.headers = h
	return t
}

func (t *textTable) Rows(rows [][]string) Table {
	t.rows = rows
	return t
}

func (t *textTable) AppendRow(row []string) Table {
	t.rows = append(t.rows, row)
	return t
}

func (t *textTable) Flush() error {
	ta := table.NewWriter()
	ta.SetOutputMirror(t.w)

	if t.headers != nil {
		ta.AppendHeader(stringsToRow(t.headers))
	}

	if t.rows != nil {
		for _, row := range t.rows {
			ta.AppendRow(stringsToRow(row))
		}
	}

	ta.Render()
	return nil
}
