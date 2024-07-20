package csv_parser

import "strings"

type Row struct {
	data map[string]string
	headers []string
}

func NewRow(headers []string, record []string) *Row {
	row := &Row{
		data: make(map[string]string),
		headers: headers,
	}

	for i, header := range headers {
		row.data[header] = record[i]
	}

	return row
}

func (r *Row) Str() string {
	return strings.Join(r.Values(), ",")
}

func (r *Row) Values() []string{
	values := make([]string, 0, len(r.data))

	for _, key := range r.headers {
		values = append(values, r.data[key])
	}

	return values
}

func (r *Row) Only(keys ...string) *Row {
	if len(keys) == 0 {
		return r
	}

	newFilteredRowData:= make([]string, len(keys))

	for idx, key:= range keys {
		newFilteredRowData[idx] = r.data[key]
	}

	return NewRow(keys, newFilteredRowData) 
}