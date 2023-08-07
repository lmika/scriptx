package scriptx_test

import (
	"github.com/bitfield/script"
	"github.com/lmika/scriptx"
	"github.com/lmika/scriptx/lines"
)

func ExampleLines() {
	weather := []string{"sun", "rain", "snow", "wind", "cloud"}

	script.Slice(weather).
		Filter(scriptx.Lines(lines.Equals("rain"), lines.Equals("cloud"))).
		Stdout()
	// Output:
	// rain
	// snow
	// wind
}
