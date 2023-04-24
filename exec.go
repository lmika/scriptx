package scriptx

import (
	"fmt"
	"os"
)

// Runnable is a pipeline that can be executed by Run. The *script.Pipe type implements this interface
type Runnable interface {
	// Stdout invokes the runnable task with the result written to stdout, if successful.  It returns the number
	// of bytes written out, or an error if the task was unsuccessful.
	Stdout() (int, error)
}

// Run takes a number of pipes and executes them in order, until they all succeed or the first error is encountered.
// If a runnable task fails with an error, it will write the error to stderr and terminate the program with exit code 1.
// It's designed to be used within the "main()" function.
func Run(pipes ...Runnable) {
	for _, pipe := range pipes {
		if _, err := pipe.Stdout(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
}
