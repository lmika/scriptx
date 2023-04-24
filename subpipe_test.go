package scriptx_test

import (
	"testing"

	"github.com/bitfield/script"
	"github.com/lmika/scriptx"
)

func TestSubPipe(t *testing.T) {
	t.Run("should run the subpipe for each line in the parent pipe", func(t *testing.T) {
		verifyPipeLines(t, func(p *script.Pipe) *script.Pipe {
			return p.Filter(scriptx.SubPipe(func (src *script.Pipe) *script.Pipe {
				return src.Replace("Line", "Subpipe")
			}))
		}, []string{
			"Line 1",
			"Line 2",
			"Line 3",
		}, []string{
			"Subpipe 1",
			"Subpipe 2",
			"Subpipe 3",
		})
	})
}