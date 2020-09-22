package tle

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ut-astria/spi/prop"
)

func TestDoTLEsBasic(t *testing.T) {
	txt := `0 DOVE 2 0505
1 39132U PLANET   20053.08335648  .00000000  00000+0  26885-3 0    09
2 39132 064.8781 046.1432 0037470 277.8993 081.9081 15.07480396    00
0 DOVE 3 0711
1 39429U PLANET   20053.08335648  .00000000  00000+0 -71889-3 0    03
2 39429 097.7682 002.8290 0145653 174.4230 093.9191 14.61124738    03
0 FLOCK 1C 1 0903
1 40027U PLANET   20053.08335648  .00000000  00000+0  34079-4 0    09
2 40027 097.9376 341.9424 0011991 325.7224 261.4362 14.90036668    06
0 FLOCK 1C 2 0904
1 40029U PLANET   20053.08335648  .00000000  00000+0 -69446-4 0    08
2 40029 097.9353 341.7057 0012795 321.8632 209.5180 14.90002937    09
0 FLOCK 1C 3 0905
1 40041U PLANET   20053.08335648  .00000000  00000+0  13294-3 0    00
2 40041 097.9374 341.8638 0012061 313.9023 333.3021 14.89832895    06
0 FLOCK 1C 4 0906
1 40031U PLANET   20053.08335648  .00000000  00000+0 -88381-4 0    00
2 40031 097.9366 341.3057 0013505 318.8100 219.1319 14.89604241    03
`

	r := strings.NewReader(txt)
	t0 := time.Now().UTC()

	f := func(i int, line0 string, p prop.Propagator) error {
		e, err := p.Prop(time.Now())
		if err != nil {
			return err
		}
		s := p.(*SGP4TLE)
		fmt.Printf("%d %s %v epoch=%v age=%v\n",
			i, s.Type(), e, s.ApproxEpoch(), s.ApproxAge(t0))
		return nil
	}

	if err := DoTLEs(bufio.NewReader(r), NewSGP4TLE, f); err != nil {
		t.Fatal(err)
	}
}

func TestDoTLEsConcurrent(t *testing.T) {
	bs, err := ioutil.ReadFile("../data/test.tle")
	if err != nil {
		t.Skip(err)
	}
	txt := string(bs)

	var (
		r    = strings.NewReader(txt)
		t0   = time.Now().UTC()
		es   = make([]prop.Ephemeris, 0, 32)
		ps   = make([]prop.Propagator, 0, 32)
		back = make([]prop.Propagator, 0, 32)
		ss   = make([]string, 0, 32)

		str = func(p prop.Propagator) string {
			return fmt.Sprintf("%#v", p.(*SGP4TLE).tle)
		}
	)

	f := func(i int, line0 string, p prop.Propagator) error {
		ps = append(ps, p)
		p0 := p
		back = append(ps, p0)
		ss = append(ss, str(p))
		e, err := p.Prop(t0)
		if err != nil {
			return err
		}
		e1, err := p.Prop(t0)
		if err != nil {
			return err
		}
		if e != e1 {
			log.Fatalf("%v %v", e, e1)
		}
		es = append(es, e)

		return nil
	}

	if err := DoTLEs(bufio.NewReader(r), NewSGP4TLE, f); err != nil {
		t.Fatal(err)
	}
	var (
		wg    = sync.WaitGroup{}
		count = int64(0)
		diff  = int64(0)
		same  = int64(0)

		delta = func(s1, s2 string) string {
			var acc string
			ss1 := strings.Split(s1, ",")
			ss2 := strings.Split(s2, ",")
			for i, line := range ss1 {
				if line != ss2[i] {
					acc += fmt.Sprintf("%d '%s' '%s'\n", i, line, ss2[i])
				}
			}
			return acc
		}
	)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			for k := 0; k < 100; k++ {
				for j, p := range ps {
					atomic.AddInt64(&count, 1)
					e, err := p.Prop(t0)
					if err != nil {
						log.Fatal(err)
					}
					if e != es[j] {
						e2, _ := p.Prop(t0)
						e3, _ := back[j].Prop(t0)
						fmt.Printf("%v\n%v\n%v\n%v\n%v\n%#v\n%v\n%v\n",
							e, es[j], e2, e3,
							p == back[j],
							str(p) == str(back[j]),
							true,
							delta(str(p), ss[j]))
						atomic.AddInt64(&diff, 1)
					} else {
						atomic.AddInt64(&same, 1)
					}
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	log.Printf("count: %d, same: %d, diff: %d", count, same, diff)

	for j, p := range ps {
		if p != back[j] {
			log.Fatalf("%v %v", p, back[j])
		}
	}

}
