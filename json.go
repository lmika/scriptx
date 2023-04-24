package scriptx

import (
	"encoding/json"
	"io"
)

func ToJSON(fn func(line string) any) func(r io.Reader, w io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		return eachLine(r, func(line string) (err error) {
			jsonBytes, err := json.Marshal(fn(line))
			if err != nil {
				return err
			}
			if _, err = w.Write(jsonBytes); err != nil {
				return err
			}
			_, err = w.Write([]byte{'\n'})
			return err
		})
	}
}
