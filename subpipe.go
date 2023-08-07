package scriptx

import (
	"bufio"
	"io"

	"github.com/bitfield/script"
)

func eachLine(r io.Reader, fn func(n int, line string) error) error {
	scanner := bufio.NewScanner(r)
	n := 1
	for scanner.Scan() {
		if err := fn(n, scanner.Text()); err != nil {
			return err
		}
		n++
	}
	return scanner.Err()
}

func SubPipe(fn func(line string) *script.Pipe) func(r io.Reader, w io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			bts, err := fn(scanner.Text()).Bytes()
			if err != nil {
				return err
			}
			if _, err := w.Write(bts); err != nil {
				return err
			}
		}
		return scanner.Err()
	}
}
