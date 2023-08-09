package scriptx_test

import (
	"testing"

	"github.com/bitfield/script"
	"github.com/lmika/scriptx"
	"github.com/lmika/scriptx/linematchers"
)

func TestLinesBetween(t *testing.T) {
	t.Run("should return all lines from 'from' inclusive to 'to' exclusive", func(t *testing.T) {
		verifyPipeLines(t, func(p *script.Pipe) *script.Pipe {
			return p.Filter(scriptx.LinesBetween(
				linematchers.Equals("snow"),
				linematchers.Equals("cloud"),
			))
		}, []string{
			"sun",
			"rain",
			"snow",
			"wind",
			"frost",
			"cloud",
			"fog",
		}, []string{
			"snow",
			"wind",
			"frost",
		})
	})

	t.Run("should return multiple set of lines matching 'from' and to", func(t *testing.T) {
		verifyPipeLines(t, func(p *script.Pipe) *script.Pipe {
			return p.Filter(scriptx.LinesBetween(
				linematchers.Equals("sun"),
				linematchers.Equals("snow"),
			))
		}, []string{
			"sun",
			"rain",
			"snow",
			"sun",
			"frost",
			"snow",
		}, []string{
			"sun",
			"rain",
			"sun",
			"frost",
		})
	})

	t.Run("should return all lines until EOF if no 'to' matches", func(t *testing.T) {
		verifyPipeLines(t, func(p *script.Pipe) *script.Pipe {
			return p.Filter(scriptx.LinesBetween(
				linematchers.Equals("wind"),
				linematchers.Equals("meatballs"),
			))
		}, []string{
			"sun",
			"rain",
			"snow",
			"wind",
			"frost",
			"cloud",
			"fog",
		}, []string{
			"wind",
			"frost",
			"cloud",
			"fog",
		})
	})

	t.Run("should emit nothing if 'from' is not matched", func(t *testing.T) {
		verifyPipeLines(t, func(p *script.Pipe) *script.Pipe {
			return p.Filter(scriptx.LinesBetween(
				linematchers.Equals("meatballs"),
				linematchers.Equals("wind"),
			))
		}, []string{
			"sun",
			"rain",
			"snow",
			"wind",
			"frost",
			"cloud",
			"fog",
		}, []string{})
	})
}

func ExampleLinesBetween() {
	weather := []string{"sun", "rain", "snow", "wind", "cloud"}

	script.Slice(weather).
		Filter(scriptx.LinesBetween(
			linematchers.Equals("rain"),
			linematchers.Equals("cloud"),
		)).
		Stdout()
	// Output:
	// rain
	// snow
	// wind
}
