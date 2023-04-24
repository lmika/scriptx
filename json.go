package scriptx

import (
	"encoding/json"
	"io"
)

// ToJSON is a filter function which converts each line from the source to a JSON structure.
// The result is a line-terminated list of JSON objects.
// The passed in function is to return a Go value that can be marshalled to JSON value.
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
