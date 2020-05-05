package dialog

import (
	"encoding/json"
	"io"
)

type jsonTable struct {
	w       io.Writer
	headers []string
	rows    [][]string
}

func stringsToMap(headers, row []string) map[string]string {
	res := make(map[string]string, len(headers))

	for idx, h := range headers {
		res[h] = row[idx]
	}

	return res
}

func allStringsToMap(headers []string, rows [][]string) []map[string]string {
	res := make([]map[string]string, 0, len(rows))

	for _, row := range rows {
		res = append(res, stringsToMap(headers, row))
	}

	return res
}

func (t *jsonTable) Headers(h []string) Table {
	t.headers = h
	return t
}

func (t *jsonTable) Rows(rows [][]string) Table {
	t.rows = rows
	return t
}

func (t *jsonTable) AppendRow(row []string) Table {
	t.rows = append(t.rows, row)
	return t
}

func (t *jsonTable) Flush() error {
	return json.NewEncoder(t.w).Encode(allStringsToMap(t.headers, t.rows))
}
