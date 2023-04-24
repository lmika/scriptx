package scriptx

import (
	"encoding/csv"
	"errors"
	"io"
)

func CSVColumn(name string) func(r io.Reader, w io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		cr := csv.NewReader(r)

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
			if len(rec) < colIdx {
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