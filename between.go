package scriptx

import (
	"io"
)

// LinesBetween is a filter function that will emit all lines starting from the first that matches
// from (inclusive) up until the first that matches to (exclusive).
func LinesBetween(from LineMatcher, to LineMatcher) func(r io.Reader, w io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		var (
			matcher      = from
			waitingForTo = false
		)

		return eachLine(r, func(n int, line string) error {
			if matcher.Match(n, line) {
				waitingForTo = !waitingForTo

				if waitingForTo {
					matcher = to
				} else {
					matcher = from
				}
			}

			if !waitingForTo {
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
