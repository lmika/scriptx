package scriptx

import (
	"io"

	"github.com/bitfield/script"
)

// EachLinesSubpipe runs the subpipe for the collection of lines matching the line matcher
func LinesBetweenSubpipe(from LineMatcher, to LineMatcher, subpipe func(lines []string) *script.Pipe) func(r io.Reader, w io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		var (
			matcher      = from
			waitingForTo = false
			lines        = []string{}
		)

		if err := eachLine(r, func(n int, line string) error {
			if matcher.Match(n, line) {
				if waitingForTo && len(lines) > 0 {
					bts, err := subpipe(lines).Bytes()
					if err != nil {
						return err
					}

					_, err = w.Write(bts)
					if err != nil {
						return err
					}
					// _, err = w.Write([]byte{'\n'})
					// if err != nil {
					// return err
					// }
				}

				lines = []string{}
				waitingForTo = !waitingForTo

				if waitingForTo {
					matcher = to
				} else {
					matcher = from
				}
			}

			if !waitingForTo {
				return nil
			} else {
				lines = append(lines, line)
			}
			return nil

			// _, err := io.WriteString(w, line)
			// if err != nil {
			// 	return err
			// }
			// _, err = w.Write([]byte{'\n'})
			// return err
		}); err != nil {
			return err
		}

		if len(lines) > 0 {
			bts, err := subpipe(lines).Bytes()
			if err != nil {
				return err
			}

			_, err = w.Write(bts)
			if err != nil {
				return err
			}
			// if _, err = w.Write([]byte{'\n'}); err != nil {
				// return err
			// }
		}
		return nil
	}
}

// LinesBetween is a filter function that will emit all lines starting from the first that matches
// from (inclusive) up until the first that matches to (exclusive).
func LinesBetween(from LineMatcher, to LineMatcher) func(r io.Reader, w io.Writer) error {
	return LinesBetweenSubpipe(from, to, func(lines []string) *script.Pipe {
		return script.Slice(lines)
	})
	/*
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
	*/
}

type LineMatcher interface {
	Match(n int, line string) bool
}
