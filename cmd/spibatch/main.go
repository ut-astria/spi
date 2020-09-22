package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ut-astria/spi/index"
	"github.com/ut-astria/spi/node"
	"github.com/ut-astria/spi/prop"
	"github.com/ut-astria/spi/tle"
)

func main() {

	log.SetFlags(log.Lmicroseconds | log.LUTC)

	n := node.NewNode(nil)

	flag.BoolVar(&n.Scan, "scan", n.Scan, "Scan for sub-tick proximity")
	flag.IntVar(&n.SlowSample, "slow-sample", n.SlowSample, "Sample during slow approaches")
	flag.IntVar(&n.Horizon, "horizon", n.Horizon, "Horizon in number of ticks")
	flag.DurationVar(&n.Resolution, "resolution", n.Resolution, "Tick duration")
	flag.IntVar(&n.IndexLevel, "index-level", n.IndexLevel, "Index's cells level")

	var (
		// Vars that aren't direct Node.Cfg fields.

		filename = flag.String("filename", "", "TLE filename (or stdin if empty)")

		// Argh. No flag.Float32Var!

		scanDist = flag.Float64("scan-dist", float64(n.ScanDist),
			"Threshold for sub-tick scan distance")

		indexDist = flag.Float64("index-dist", float64(n.IndexDist),
			"Threshold for index search distance")

		slowSampleThreshold = flag.Float64("slow-sample-threshold", float64(n.SlowSampleThreshold),
			"Maximum relative speed (m/s) to trigger slow approach sampling")

		numWorkers = flag.Int("workers", runtime.NumCPU(), "Nummber of workers")

		cfg = flag.String("cfg", "", "Filename for JSON configuration; overrides any other args")
		// ToDo: Command-line args override cfg file.

		ts = flag.String("t0", "now", "Logical starting time (example: \"2020-09-18T17:31:16Z\")")
	)

	flag.Parse()

	if *ts == "now" {
		*ts = time.Now().UTC().Format(time.RFC3339)
	}

	t0, err := time.Parse(time.RFC3339Nano, *ts)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("t0: %s", *ts)

	n.ScanDist = float32(*scanDist)
	n.IndexDist = float32(*indexDist)
	n.SlowSampleThreshold = float32(*slowSampleThreshold)

	if *cfg != "" {
		js, err := ioutil.ReadFile(*cfg)
		if err != nil {
			log.Fatal(err)
		}
		if err = json.Unmarshal(js, &n.Cfg); err != nil {
			log.Fatal(err)
		}
	}

	var bs []byte
	if *filename == "" {
		bs, err = ioutil.ReadAll(os.Stdin)
	} else {
		bs, err = ioutil.ReadFile(*filename)
	}
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(bytes.NewReader(bs))

	sats := make([]*node.PubTLE, 0, 1024)

	f := func(i int, line0 string, p prop.Propagator) error {
		if err = tle.Check(p); err != nil {
			log.Printf("sat propagation error at TLE %d (%s): %s", i, strings.TrimSpace(line0), err)
			return nil
		}
		sat := &node.PubTLE{
			TLE: p.(*tle.SGP4TLE),
		}
		sats = append(sats, sat)
		return nil
	}

	if err := tle.DoTLEs(r, tle.NewSGP4TLE, f); err != nil {
		log.Fatal(err)
	}

	log.Printf("loaded %d TLEs", len(sats))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n.Prepare(ctx)

	iis := make([]*node.IndexInput, 0, len(sats))
	for _, sat := range sats {
		ii := n.NewIndexInput(ctx, sat)
		if ii == nil {
			continue
		}
		iis = append(iis, ii)
	}

	var (
		tick  = n.Resolution
		ticks = int(float64(n.Horizon) / float64(n.Resolution/time.Second))
		t     = t0
		then  = t.Add(tick * time.Duration(ticks))

		workers = make(chan bool, *numWorkers)
	)

	log.Printf("%d ticks over %d secs", ticks, n.Horizon)

	log.Printf("using %d workers", *numWorkers)
	for i := 0; i < *numWorkers; i++ {
		workers <- true
	}

	for t.Before(then) {
		work := func(t time.Time) {
			// log.Printf("starting %v", t)

			i := n.NewIndex(t)
			ios := make([]*node.IndexOutput, 0, len(iis))
			for _, ii := range iis {

				e, err := ii.Sat.TLE.Prop(t)
				if err != nil {
					log.Printf("sat.Prop %s", err)
					continue
				}

				pp := index.ProbPos{
					Pos: index.Pos{
						X: e.ECI.X,
						Y: e.ECI.Y,
						Z: e.ECI.Z,
					},
				}

				pps := []index.ProbPos{pp}

				cans, novs, _, err := i.I.Update(ii.Id, ii.Key, pps)

				io := &node.IndexOutput{
					Time:     t,
					Novel:    novs,
					Canceled: cans,
				}
				ios = append(ios, io)
			}

			ps := n.GetIndexOutputTLEs(ctx, ios)
			rs := make([]*node.Report, 0, 4)

			for _, io := range ios {

				for _, c := range io.Novel {
					r, err := n.ConjToReport(io.Time, &c, n.ScanDist, ps, false)
					if err != nil {
						panic(err)
					}
					if r == nil {
						continue
					}
					rs = append(rs, r)
				}

				for _, c := range io.Canceled {
					r, err := n.ConjToReport(io.Time, &c, n.ScanDist, ps, true)
					if err != nil {
						panic(err)
					}
					if r == nil {
						continue
					}
					r.Canceled = true
					rs = append(rs, r)
				}

			}
			for _, r := range rs {
				fmt.Printf("%s\n", JSON(r, false))
			}

			select {
			case <-ctx.Done():
			case workers <- true:
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-workers:
			go work(t)
		}

		t = t.Add(tick)
	}
}

func JSON(x interface{}, pretty bool) string {
	var js []byte
	var err error
	if pretty {
		js, err = json.MarshalIndent(&x, "", "  ")
	} else {
		js, err = json.Marshal(&x)
	}
	if err != nil {
		return fmt.Sprintf("%#v", x)
	}
	return string(js)
}
