package scriptx

import (
	"fmt"
	"os"
)

type Runnable interface {
	Stdout() (int, error)
}

func Run(pipes ...Runnable) {
	for _, pipe := range pipes {
		if _, err := pipe.Stdout(); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
}
