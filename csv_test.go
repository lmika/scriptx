package scriptx_test

import (
	"regexp"
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

func ExampleCSVColumn() {
	script.Slice([]string{
		"letter,fruit,word",
		"a,apple,alpha",
		"b,banana,bravo",
		"c,cherry,charlie",
	}).Filter(scriptx.CSVColumn("fruit")).Stdout()
	// Output:
	// apple
	// banana
	// cherry
}

func ExampleToCSV() {
	script.Slice([]string{
		"letter  fruit   word",
		"a       apple   alpha",
		"b       banana  bravo",
		"c       cherry  charlie",
	}).Filter(scriptx.ToCSV(regexp.MustCompile(`\s+`))).Stdout()
	// Output:
	// letter,fruit,word
	// a,apple,alpha
	// b,banana,bravo
	// c,cherry,charlie
}

func ExampleCSVMapToJSON() {
	script.Slice([]string{
		"letter,fruit,word",
		"a,apple,alpha",
		"b,banana,bravo",
		"c,cherry,charlie",
	}).Filter(scriptx.CSVMapToJSON(func(row []string, header *scriptx.CSVHeader) any {
		return struct {
			Letter string `json:"letter"`
			Fruit  string `json:"fruit"`
			Word   string `json:"word"`
		}{
			Letter: header.Value(row, "letter"),
			Fruit:  header.Value(row, "fruit"),
			Word:   header.Value(row, "word"),
		}
	})).Stdout()
	// Output:
	// {"letter":"a","fruit":"apple","word":"alpha"}
	// {"letter":"b","fruit":"banana","word":"bravo"}
	// {"letter":"c","fruit":"cherry","word":"charlie"}
}
