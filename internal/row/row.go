package row

type Row struct {
	rowNumber int
	data      map[string]string
	headers   []string
}

func NewRow(rowNumber int, headers []string, record []string) *Row {
	row := &Row{
		rowNumber: rowNumber,
		data:      make(map[string]string),
		headers:   headers,
	}

	for i, header := range headers {
		row.data[header] = record[i]
	}

	return row
}

func (r *Row) Values() []string {
	values := make([]string, len(r.data))
	i := 0
	for _, key := range r.headers {
		values[i] = r.data[key]
		i++
	}

	return values
}

func (r *Row) LineNumber() int {
	return r.rowNumber
}

func (r *Row) Only(keys []string) *Row {
	if len(keys) == 0 {
		return r
	}

	newFilteredRowData := make([]string, len(keys))

	for idx, key := range keys {
		newFilteredRowData[idx] = r.data[key]
	}

	return NewRow(r.rowNumber, keys, newFilteredRowData)
}

func (r *Row) HasColumn(key string) bool {
	_, ok := r.data[key]
	return ok
}

func (r *Row) GetColumn(key string) string {
	return r.data[key]
}
