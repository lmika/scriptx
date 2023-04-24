package scriptx

import (
	"fmt"
	"io"
)

func Printf(ptrn string) func(r io.Reader, w io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		return eachLine(r, func(line string) error {
			_, err := io.WriteString(w, fmt.Sprintf(ptrn, line))
			if err != nil {
				return err
			}
			_, err = w.Write([]byte{'\n'})
			return err
		})
	}
}