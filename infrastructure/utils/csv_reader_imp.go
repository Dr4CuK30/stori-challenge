package utils

import (
	"encoding/csv"
	"io"
)

type CsvReaderImp struct{}

func (r *CsvReaderImp) ReadCSV(reader io.Reader) ([][]string, error) {
	var records [][]string
	csvReader := csv.NewReader(reader)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records[1:len(records)], nil
}
