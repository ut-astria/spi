package node

import (
	"fmt"
	"time"

	"github.com/ut-astria/spi/prop"
)

func ScanPair(doit bool, t time.Time, tick time.Duration, samples uint64, pa, pb prop.Propagator, cutDist float32) (float32, []prop.Ephemeris, time.Time, error) {

	var (
		half      = tick / 2
		t0        = t.Add(-half)
		t1        = t.Add(half)
		increment = tick / time.Duration(samples)

		closest = float32(-1)

		a, b prop.Ephemeris
		then time.Time
	)

	// ToDo: Use speed to plan the search.
	//
	// ToDo: Even something fancier, but we still need to be fast.
	// Leave high-precision analysis to downstream processing.
	for at := t0; at.Before(t1); at = at.Add(increment) {
		var (
			a0, erra = pa.Prop(at)
			b0, errb = pb.Prop(at)
		)
		if erra != nil {
			return 0, nil, at, erra
		}

		if errb != nil {
			return 0, nil, at, errb
		}

		d := a0.ECI.Dist(b0.ECI)

		if closest < 0 || d < closest {
			closest = d
			then = at
			a, b = a0, b0
			if d < cutDist {
				break
			}
		}

		if !doit && at == t {
			// Should avoid preceeding loop iterations.
			return d, []prop.Ephemeris{a0, b0}, at, nil
		}
	}

	if !doit {
		return 0, nil, t, fmt.Errorf("shouldn't have scanned")
	}

	return closest, []prop.Ephemeris{a, b}, then, nil
}
