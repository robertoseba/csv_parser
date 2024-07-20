package reader

import "strings"

type Row struct {
	data map[string]string
}

func NewRow(headers []string, record []string) *Row {
	row := &Row{
		data: make(map[string]string),
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

	for _, value := range r.data {
		values = append(values, value)
	}

	return values
}