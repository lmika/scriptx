package scriptx

import (
	"io"
)

// FirstLinesBetween is a filter function that will emit all lines starting from the first that matches
// from (inclusive) up until the first that matches to (exclusive). Only the first set of lines matching
// the two matchers will be emitted.
func FirstLinesBetween(from LineMatcher, to LineMatcher) func(r io.Reader, w io.Writer) error {
	const (
		stateWaitingForFrom int = iota
		stateWaitingForTo
		stateDone
	)

	return func(r io.Reader, w io.Writer) error {
		var (
			nextMatcher = from
			state       = stateWaitingForFrom
		)

		return eachLine(r, func(n int, line string) error {
			if state >= stateDone {
				return nil
			}

			if nextMatcher.Match(n, line) {
				state++
				nextMatcher = to
			}

			if state != stateWaitingForTo {
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
