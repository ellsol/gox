package gox

import (
	"os"
	"encoding/csv"
)

type CSVFile struct {
	filename string
	handle *os.File
	writer *csv.Writer
}

func (it *CSVFile) Open(filename string) (*CSVFile, error) {

	handle, err := OpenFile(filename)

	if err != nil {
		return nil, err
	}

	it.filename = filename
	it.handle = handle
	it.writer = csv.NewWriter(handle)

	return it, nil
}

func (it *CSVFile) Append(row []string) error {
	return it.writer.Write(row)
}


func (it *CSVFile) AppendRow(row *CSVRow) error {
	return it.writer.Write(row.Values())
}

func(it *CSVFile) Flush()  {
	it.writer.Flush()
}

func (it *CSVFile) Close() error {
	it.writer.Flush()
	return it.handle.Close()
}

