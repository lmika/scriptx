package scriptx

import (
	"io"
)

// Lines
func Lines(from LineMatcher, to LineMatcher) func(r io.Reader, w io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		var (
			nextMatcher = from
			sendOut     = false
		)

		return eachLine(r, func(n int, line string) error {
			if nextMatcher.Match(n, line) {
				sendOut = !sendOut
				nextMatcher = to
			}

			if !sendOut {
				return nil
			}

			_, err := io.WriteString(w, line)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte{'\n'})
			return err
		})
	}
}

type LineMatcher interface {
	Match(n int, line string) bool
}
