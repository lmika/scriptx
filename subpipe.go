package scriptx

import (
	"bufio"
	"io"

	"github.com/bitfield/script"
)

func eachLine(r io.Reader, fn func(line string) error) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := fn(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func SubPipe(fn func(src *script.Pipe) *script.Pipe) func(r io.Reader, w io.Writer) error {
	return func(r io.Reader, w io.Writer) error {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			bts, err := fn(script.Echo(scanner.Text())).Bytes()
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
