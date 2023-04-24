package scriptx_test

import (
	"testing"

	"github.com/bitfield/script"
	"github.com/lmika/scriptx"
)

func TestCSVColumn(t *testing.T) {
	t.Run("should return named column of CSV input", func(t *testing.T) {
		verifyPipeLines(t, func(p *script.Pipe) *script.Pipe {
			return p.Filter(scriptx.CSVColumn("state"))
		}, []string{
			"city,state,country",
			"Melbourne,Vic,AU",
			"Sydney,NSW,AU",
			"Canberra,\"Territory, Australian Capital\",AU",
			"New York,NY,US,Extra,Comments,Allowed",
			"Missing state is ignored",
			"# Comment is ignored",
			"Seattle,WA",
		}, []string{
			"Vic",
			"NSW",
			"Territory, Australian Capital",
			"NY",
			"WA",
		})
	})
}
