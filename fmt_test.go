package scriptx_test

import (
	"testing"

	"github.com/bitfield/script"
	"github.com/lmika/scriptx"
	"github.com/stretchr/testify/assert"
)

func TestPrintf(t *testing.T) {
	t.Run("should format each line according to the pattern", func(t *testing.T) {
		verifyPipeLines(t, func(p *script.Pipe) *script.Pipe {
			return p.Filter(scriptx.Printf("Line: [%v]"))
		}, []string{
			"Line 1",
			"Line 2",
			"Line 3",
		}, []string{
			"Line: [Line 1]",
			"Line: [Line 2]",
			"Line: [Line 3]",
		})
	})
}

func verifyPipeLines(t *testing.T, fn func(p *script.Pipe) *script.Pipe, inLines []string, expOut []string) {
	pipe := script.Slice(inLines)
	outLines, err := fn(pipe).Slice()

	assert.NoError(t, err)
	assert.Equal(t, expOut, outLines)
}

func verifyPipeErr(t *testing.T, fn func(p *script.Pipe) *script.Pipe, inLines []string) {
	pipe := script.Slice(inLines)
	_, err := fn(pipe).Slice()

	assert.Error(t, err)
}
