package csv

import "fmt"

const (
	OutOfBoundFormat = "out of bound, aked for %v, but length=%v"
)

type CSVRow struct {
	values []string
}

func NewCSVRow() *CSVRow {
	return NewCSVRowWith(make([]string, 0))
}

func NewCSVRowWith(initialValues []string) *CSVRow {
	return &CSVRow{
		values: initialValues,
	}
}

func (it *CSVRow) Add(item string) *CSVRow {
	it.values = append(it.values, item)
	return it
}

func (it *CSVRow) Length() int {
	return len(it.values)
}

func (it *CSVRow) GetString(position int) (string, error) {
	if position >= it.Length() {
		return "", fmt.Errorf(OutOfBoundFormat, position, it.Length())
	}

	return it.values[position], nil
}

func (it *CSVRow) Values() []string {
	return it.values
}