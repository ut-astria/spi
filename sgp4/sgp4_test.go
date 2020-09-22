package sgp4

import (
	"log"
	"testing"
	"time"
)

func BenchmarkSGP4Vallado(b *testing.B) {
	var (
		mins     = float64(60 * 24)
		line1    = "1 39132U PLANET   20016.08334491  .00000000  00000+0 -47542-3 0    07"
		line2    = "2 39132 064.8760 163.6520 0036285 284.0373 175.5769 15.07452065    00"
		tle, err = ParseLines(line1, line2)
	)

	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r, v, err := tle.PropForMins(mins)
		if err != nil {
			b.Fatal(err)
		}
		if r[0] == 0 {
			b.Fatal(r)
		}
		if v[0] == 0 {
			b.Fatal(v)
		}
	}
}

func TestError1(t *testing.T) {
	var (
		line1    = "1 44246U 19029M   19348.91667824  .00396868  00000-0  14121-2 0  9990"
		line2    = "2 44246  52.9936 253.6898 0006659 333.1525 350.4837 15.87356006 31853"
		tle, err = ParseLines(line1, line2)

		now     = time.Now().UTC()
		horizon = 60 * time.Minute
		inc     = time.Second
	)

	if err != nil {
		t.Fatal(err)
	}

	limit := 10
	for mins := float64(0); mins < 100000; mins += 10 {
		if _, _, err = tle.PropForMins(mins); err != nil {
			log.Println(err)
			limit--
			if limit < 0 {
				break
			}
		}
	}

	limit = 10
	for t := now; t.Before(now.Add(horizon)); t = t.Add(inc) {
		ms := t.UnixNano() / 1000 / 1000
		if _, _, err = tle.PropUnixMillis(ms); err != nil {
			log.Println(err)
			limit--
			if limit < 0 {
				break
			}
		}
	}
}
