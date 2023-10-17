package scriptx

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
)

type Splitter interface {
	Split(line string, n int) []string
}

// ToCSV is a filter function that reads the source as a series of lines, splits them
// into tokens using the passed in Splitter, and writes them to the output as a CSV.
func ToCSV(splitter Splitter) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		csvWriter := csv.NewWriter(w)

		scnr := bufio.NewScanner(r)
		for scnr.Scan() {
			line := scnr.Text()
			if err := csvWriter.Write(splitter.Split(line, -1)); err != nil {
				return err
			}
		}
		if err := scnr.Err(); err != nil {
			return err
		}

		csvWriter.Flush()
		return nil
	}
}

// CSVColumn is a filter function that reads the source as a CSV file and extracts the cell
// values of the named column, excluding the header itself. If the column cannot be found,
// the filter will produce nothing. If the column index is beyond the number of columns
// for a particular row, it will be skipped.
func CSVColumn(name string) func(r io.Reader, w io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		cr := csv.NewReader(r)
		cr.FieldsPerRecord = -1

		header, err := cr.Read()
		if err != nil {
			return err
		}

		var colIdx = -1
		for i, ch := range header {
			if ch == name {
				colIdx = i
			}
		}
		if colIdx == -1 {
			return nil
		}

		for {
			rec, err := cr.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				} else {
					return err
				}
			}
			if len(rec) <= colIdx {
				continue
			}

			if _, err := io.WriteString(w, rec[colIdx]); err != nil {
				return err
			}
			if _, err := w.Write([]byte{'\n'}); err != nil {
				return err
			}
		}
	}
}
