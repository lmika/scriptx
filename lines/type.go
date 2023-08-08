package lines

import (
	"strings"
)

type Matcher struct {
	matchFn func(n int, line string) bool
}

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

func Equals(matchLine string) Matcher {
	return Matcher{
		matchFn: func(n int, line string) bool {
			return line == matchLine
		},
	}
}
