package tle

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/ut-astria/spi/prop"
	"github.com/ut-astria/spi/sgp4"
)

type SGP4TLE struct {
	CatNum string

	TLE []string

	// ToDo: Deleted flag?

	tle *sgp4.TLE
}

func NewSGP4TLE(line0, line1, line2 string) (prop.Propagator, error) {
	o, err := sgp4.ParseLines(line1, line2)
	if err != nil {
		return nil, err
	}
	cat := "nocat"
	if 8 < len(line1) {
		cat = strings.TrimSpace(line1[2:7])
	}

	return &SGP4TLE{
		CatNum: cat,
		TLE:    []string{line0, line1, line2},
		tle:    o,
	}, nil
}

// Epoch attempts to determine the approximate epoch.  Returns zero on
// failure.
//
// Currently this function probably returns slightly incorrect data.
//
// ToDo: return a proper error.
func (o *SGP4TLE) ApproxEpoch() time.Time {
	var (
		nope = time.Time{}
		line = o.TLE[1]
	)

	if len(line) < 33 {
		return nope
	}

	var (
		y, yerr = strconv.Atoi(line[18:20])
		d, derr = strconv.ParseFloat(line[20:32], 64)
	)

	if yerr != nil || derr != nil {
		return nope
	}

	if 56 < y {
		y += 1900
	} else {
		y += 2000
	}

	t, err := time.Parse("2006", strconv.Itoa(y))
	if err != nil {
		return nope
	}

	t = t.Add(time.Duration((d-1)*24*60*60*1000) * time.Millisecond).UTC()

	return t
}

func (o *SGP4TLE) ApproxAge(t0 time.Time) time.Duration {
	return t0.Sub(o.ApproxEpoch())
}

func (o *SGP4TLE) Type() string {
	return GetType(o.TLE[0])
}

func (o *SGP4TLE) Name() string {
	return o.CatNum
}

func (o *SGP4TLE) Prop(t time.Time) (prop.Ephemeris, error) {
	p, v, err := o.tle.PropUnixMillis(t.UnixNano() / 1000 / 1000)
	var e prop.Ephemeris
	if err == nil {
		e = prop.Ephemeris{
			ECI: prop.Vect{float32(p[0]), float32(p[1]), float32(p[2])},
			V:   prop.Vect{float32(v[0]), float32(v[1]), float32(v[2])},
		}
	}
	return e, err
}

// Legit just calls the function Legit.
func (o *SGP4TLE) Legit() bool {
	return Legit(o)
}

// Legit is a quick check that the propagotor can be propagated as of
// right now.
func Legit(p prop.Propagator) bool {
	_, err := p.Prop(time.Now())
	return err == nil
}

// Check is a quick check that the propagotor can be propagated as of
// right now.
func Check(p prop.Propagator) error {
	_, err := p.Prop(time.Now())
	return err
}

func ParseSGP4TLE(js string) (*SGP4TLE, error) {
	var o SGP4TLE
	if err := json.Unmarshal([]byte(js), &o); err != nil {
		return nil, err
	}
	tle, err := sgp4.ParseLines(o.TLE[1], o.TLE[2])
	if err != nil {
		return nil, err
	}
	o.tle = tle
	return &o, nil
}

func (o *SGP4TLE) GetType() string {
	return GetType(o.TLE[0])
}
