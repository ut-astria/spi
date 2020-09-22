package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/pkg/profile"

	"github.com/ut-astria/spi/node"
	"github.com/ut-astria/spi/prop"
	"github.com/ut-astria/spi/tle"
)

func main() {

	log.SetFlags(log.Lmicroseconds | log.LUTC)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		n              = node.NewNode(nil)
		inputBatchSize = 128

		duration = flag.Duration("duration", 0, "Duration")
		horizon  = flag.Int("horizon", n.Horizon, "Horizon in ticks")

		indexDist       = flag.Float64("index-dist", float64(n.IndexDist), "Distance limit")
		scanDist        = flag.Float64("scan-dist", float64(n.ScanDist), "Distance limit")
		level           = flag.Int("level", n.IndexLevel, "Cells level")
		scan            = flag.Bool("scan", n.Scan, "Scan")
		sample          = flag.Int("slow-sample", n.SlowSample, "Slow sampling rate")
		sampleThreshold = flag.Float64("slow-sample-threshold", float64(n.SlowSampleThreshold), "Slow sample threshold")

		sampleMod   = flag.Int("sample-mod", 0, "Sample modulus")
		sampleRem   = flag.Int("sample-rem", 0, "Sample remainder")
		batchOutput = flag.Bool("batch-output", false, "output batches")

		logging          = flag.Bool("v", false, "Logging")
		memProf          = flag.Bool("prof-mem", false, "Enable memory profiling")
		blockProfileRate = flag.Int("block-profile-rate", 0, "Blocking profile rate")
		wg               sync.WaitGroup
	)

	flag.Parse()

	{ // Configure runtime
		if *memProf {
			defer profile.Start(profile.MemProfile).Stop()
		}

		if 0 < *blockProfileRate {
			runtime.SetBlockProfileRate(*blockProfileRate)
		}
	}

	pub := func(topic, msg string) {
		js := fmt.Sprintf(`{"%s":%s"}`, topic, msg)
		if _, err := fmt.Printf("%s\n", js); err != nil {
			panic(err)
		}
	}

	read := func(publisher string, r *bufio.Reader) []*node.PubTLE {
		var batch []*node.PubTLE
		t := time.NewTimer(time.Second)

		emit := func(force bool) {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if 0 < len(batch) {
					log.Printf("incoming batch: %d", len(batch))
					select {
					case <-ctx.Done():
					case n.In <- batch:
					}
				}
				t = time.NewTimer(time.Second)
			default:
				if 0 < len(batch) {
					if force || len(batch) == inputBatchSize {
						log.Printf("incoming batch: %d", len(batch))
						select {
						case <-ctx.Done():
						case n.In <- batch:
						}
					}
				}
			}

			batch = nil
		}

		f := func(i int, line0 string, p prop.Propagator) error {
			_, err := p.Prop(time.Now())
			if err != nil {
				log.Printf("sat error at %d %s: %s", i, strings.TrimSpace(line0), err)
				return nil
			}
			t := p.(*tle.SGP4TLE)

			if *sampleMod != 0 {
				k := []byte(t.CatNum)
				hash := fnv.New32()
				hash.Write(k)
				h := int(hash.Sum32())
				if *sampleRem != h%*sampleMod {
					return nil
				}
			}

			sat := &node.PubTLE{
				Publisher: publisher,
				TLE:       t,
			}

			if batch == nil {
				batch = make([]*node.PubTLE, 0, inputBatchSize)
			}

			batch = append(batch, sat)

			if len(batch) == inputBatchSize {
				emit(false)
			}

			return nil
		}

		if err := tle.DoTLEs(r, tle.NewSGP4TLE, f); err != nil {
			panic(err)
		}

		emit(true)

		return nil
	}

	n.Horizon = *horizon
	n.IndexLevel = *level
	n.IndexDist = float32(*indexDist)
	n.ScanDist = float32(*scanDist)
	n.Logging = *logging
	n.SlowSample = *sample
	n.SlowSampleThreshold = float32(*sampleThreshold)
	n.Scan = *scan

	n.Metrics = make(chan node.Metrics)
	n.Errs = make(chan error)
	n.Out = make(chan []*node.Report, 32)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-n.Errs:
				js := JSON(map[string]interface{}{
					"error": err.Error(),
				}, false)
				pub("errs", js)
			case m := <-n.Metrics:
				pub("metrics", JSON(m, false))
			case rs := <-n.Out:
				if *batchOutput {
					pub("reports", JSON(rs, false))
				} else {
					for _, r := range rs {
						pub("report", JSON(r, false))
					}
				}
			}
		}
	}()

	go n.Run(ctx)

	go func() {
		r := bufio.NewReader(os.Stdin)
		defer os.Stdin.Close()
		if err := read("stdin", r); err != nil {
			panic(err)
		}
		log.Printf("Reader done")
	}()

	if 0 < *duration {
		time.Sleep(*duration)
	} else {
		wg.Add(1)
		wg.Wait()
	}

	log.Printf("main done")
}

func SatBatches(sats []*node.PubTLE, n int) [][]*node.PubTLE {
	acc := make([][]*node.PubTLE, 0, 1+len(sats)/n)
	var batch []*node.PubTLE
	for _, sat := range sats {
		if batch == nil {
			batch = make([]*node.PubTLE, 0, n)
		}
		batch = append(batch, sat)
		if len(batch) == n {
			acc = append(acc, batch)
			batch = nil
		}
	}
	if 0 < len(batch) {
		acc = append(acc, batch)
	}
	return acc
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
