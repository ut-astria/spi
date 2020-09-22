// Package prop defines basic types for propagating data to a position.
package prop

import (
	"math"
	"time"
)

// Propagator can compute an Ephemeris for a given time.
type Propagator interface {
	Prop(t time.Time) (Ephemeris, error)
}

// Vect is a 3-vector.
type Vect struct {
	X, Y, Z float32
}

// Ephemeris represents position and velocity.
type Ephemeris struct {
	// V is velocity.
	V Vect

	// C is Cartesian position.
	ECI Vect
}

// Dist is the Cartesian distance metric.
func (v1 Vect) Dist(v2 Vect) float32 {
	var (
		x = float64(v1.X) - float64(v2.X)
		y = float64(v1.Y) - float64(v2.Y)
		z = float64(v1.Z) - float64(v2.Z)
	)
	return float32(math.Sqrt(x*x + y*y + z*z))
}
