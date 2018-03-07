// +build go1.10

package rokka

import "math"

// round calls math.Round for go 1.10.
// TODO: when removing support for go 1.8/1.9, remove this function and instead use `math.Round` directly.
func round(x float64) float64 {
	return math.Round(x)
}
