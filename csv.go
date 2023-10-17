package scriptx

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"sort"
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

type CSVHeader struct {
	header []string
}

func (ch *CSVHeader) Column(name string) int {
	for c, h := range ch.header {
		if h == name {
			return c
		}
	}
	return -1
}

func (ch *CSVHeader) Value(row []string, name string) string {
	c := ch.Column(name)
	if c == -1 || c >= len(row) {
		return ""
	}
	return row[c]
}

func CSVFilter(fn func(row []string, header *CSVHeader) []string) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		csvr := csv.NewReader(r)
		csvr.FieldsPerRecord = 0

		csvWriter := csv.NewWriter(w)

		header, rErr := csvr.Read()
		if rErr != nil {
			return rErr
		}
		if err := csvWriter.Write(header); err != nil {
			return err
		}
		headerInfo := CSVHeader{header: header}

		row, rErr := csvr.Read()
		for !errors.Is(rErr, io.EOF) {
			newRow := fn(row, &headerInfo)
			if newRow != nil {
				if err := csvWriter.Write(newRow); err != nil {
					return err
				}
			}

			row, rErr = csvr.Read()
		}
		if rErr != nil && !errors.Is(rErr, io.EOF) {
			return rErr
		}

		csvWriter.Flush()
		return nil
	}
}

func CSVSort(lessThan func(row1, row2 []string, header *CSVHeader) bool) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		csvr := csv.NewReader(r)
		csvr.FieldsPerRecord = 0

		header, rErr := csvr.Read()
		if rErr != nil {
			return rErr
		}
		headerInfo := CSVHeader{header: header}

		records, rErr := csvr.ReadAll()
		if rErr != nil {
			return rErr
		}

		sort.Slice(records, func(i, j int) bool {
			return lessThan(records[i], records[j], &headerInfo)
		})

		csvWriter := csv.NewWriter(w)

		if err := csvWriter.Write(header); err != nil {
			return err
		}
		for _, r := range records {
			if err := csvWriter.Write(r); err != nil {
				return err
			}
		}

		csvWriter.Flush()
		return nil
	}
}

func CSVMapToString(fn func(row []string, header *CSVHeader) string) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		csvr := csv.NewReader(r)
		csvr.FieldsPerRecord = 0

		header, rErr := csvr.Read()
		if rErr != nil {
			return rErr
		}
		headerInfo := CSVHeader{header: header}

		row, rErr := csvr.Read()
		for !errors.Is(rErr, io.EOF) {
			line := fn(row, &headerInfo)
			if _, err := io.WriteString(w, line); err != nil {
				return err
			}
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}

			row, rErr = csvr.Read()
		}
		if rErr != nil && !errors.Is(rErr, io.EOF) {
			return rErr
		}

		return nil
	}
}

func CSVMapToJSON(fn func(row []string, header *CSVHeader) any) func(io.Reader, io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		csvr := csv.NewReader(r)
		csvr.FieldsPerRecord = 0

		header, rErr := csvr.Read()
		if rErr != nil {
			return rErr
		}
		headerInfo := CSVHeader{header: header}

		row, rErr := csvr.Read()
		for !errors.Is(rErr, io.EOF) {
			obj := fn(row, &headerInfo)
			jsonBytes, err := json.Marshal(obj)
			if err != nil {
				return err
			}
			if _, err = w.Write(jsonBytes); err != nil {
				return err
			}
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}

			row, rErr = csvr.Read()
		}
		if rErr != nil && !errors.Is(rErr, io.EOF) {
			return rErr
		}

		return nil
	}
}
