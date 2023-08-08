package linematchers

import (
	"strings"
)

type Matcher struct {
	matchFn func(n int, line string) bool
}

// TrimSpace returns a modified line matcher which will trim the passed in string prior to testing it.
func (lm Matcher) TrimSpace() Matcher {
	return Matcher{
		matchFn: func(n int, line string) bool {
			return lm.matchFn(n, strings.TrimSpace(line))
		},
	}
}

func (lm Matcher) Match(n int, line string) bool {
	return lm.matchFn(n, line)
}

// Equals is a line matcher that will match lines that equal the passed in line.
func Equals(matchLine string) Matcher {
	return Matcher{
		matchFn: func(n int, line string) bool {
			return line == matchLine
		},
	}
}
