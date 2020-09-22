package node

import (
	"time"

	"github.com/ut-astria/spi/prop"

	sat "github.com/jsmorph/go-satellite"
)

func TimeToGST(t time.Time) (float64, float64) {
	var (
		y   = t.Year()
		m   = int(t.Month())
		d   = t.Day()
		h   = t.Hour()
		min = t.Minute()
		sec = t.Second()
		ns  = t.Nanosecond()
	)

	return sat.GSTimeFromDateNano(y, m, d, h, min, sec, ns)
}

type LatLonAlt struct {
	Lat, Lon, Alt float32
}

func ECIToLLA(t time.Time, p prop.Vect) (*LatLonAlt, error) {

	gmst, _ := TimeToGST(t)

	x := sat.Vector3{
		X: float64(p.X),
		Y: float64(p.Y),
		Z: float64(p.Z),
	}

	// sat.ECIToLLA is very slow.
	alt, _, ll := sat.ECIToLLA(x, gmst)

	d, err := sat.LatLongDeg(ll)
	if err != nil {
		return nil, err
	}

	return &LatLonAlt{
		Lat: float32(d.Latitude),
		Lon: float32(d.Longitude),
		Alt: float32(alt),
	}, nil
}
