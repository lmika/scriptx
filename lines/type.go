package lines

import "github.com/lmika/scriptx"

type lineMatcher struct {
	matchFn func(n int, line string) bool
}

func (lm lineMatcher) Match(n int, line string) bool {
	return lm.matchFn(n, line)
}

func Equals(matchLine string) scriptx.LineMatcher {
	return lineMatcher{
		matchFn: func(n int, line string) bool {
			return line == matchLine
		},
	}
}
