package row

import "slices"

type Row struct {
	rowNumber int
	data      []string
	headers   []string
}

func NewRow(rowNumber int, headers []string, record []string) *Row {
	row := &Row{
		rowNumber: rowNumber,
		data:      record,
		headers:   headers,
	}

	return row
}

func (r *Row) Values() []string {
	return r.data
}

func (r *Row) LineNumber() int {
	return r.rowNumber
}

func (r *Row) Only(keys []string) *Row {
	if len(keys) == 0 {
		return r
	}

	newFilteredRowData := make([]string, len(keys))

	for _, key := range keys {
		idx := slices.Index(r.headers, key)
		newFilteredRowData[idx] = r.data[idx]
	}

	return NewRow(r.rowNumber, keys, newFilteredRowData)
}

func (r *Row) HasColumn(key string) bool {
	idx := slices.Index(r.headers, key)
	return idx != -1
}

func (r *Row) GetColumn(key string) string {
	idx := slices.Index(r.headers, key)
	return r.data[idx]
}
