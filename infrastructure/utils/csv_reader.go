package utils

import (
	"io"
)

type CsvReader interface {
	ReadCSV(reader io.Reader) ([][]string, error)
}
