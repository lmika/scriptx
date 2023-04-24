package scriptx_test

import (
	"testing"

	"github.com/bitfield/script"
	"github.com/lmika/scriptx"
)

func TestToJSON(t *testing.T) {
	t.Run("should convert a line to a JSON object", func(t *testing.T) {
		verifyPipeLines(t, func(p *script.Pipe) *script.Pipe {
			return p.Filter(scriptx.ToJSON(func(line string) any {
				return map[string]any{
					"line":        line,
					"line_length": len(line),
				}
			}))
		}, []string{
			"Just three words",
			"She sells sea shells by the \"sea shore\"",
			"Another example of a line",
		}, []string{
			`{"line":"Just three words","line_length":16}`,
			`{"line":"She sells sea shells by the \"sea shore\"","line_length":39}`,
			`{"line":"Another example of a line","line_length":25}`,
		})
	})

	t.Run("should convert a line to a JSON array", func(t *testing.T) {
		verifyPipeLines(t, func(p *script.Pipe) *script.Pipe {
			return p.Filter(scriptx.ToJSON(func(line string) any {
				return []any{line, line}
			}))
		}, []string{
			"line",
			"words words",
			"foobar",
		}, []string{
			`["line","line"]`,
			`["words words","words words"]`,
			`["foobar","foobar"]`,
		})
	})
}
