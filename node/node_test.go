package node

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/ut-astria/spi/prop"
	"github.com/ut-astria/spi/tle"
)

func TestNode(t *testing.T) {

	filenames := []string{
		"planet_mc_20200725.tle",
		"planet_mc_20200727.tle",
		"planet_mc_20200728.tle",
		"planet_mc_20200729.tle",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		// We want to see at least this many emitted reports.
		want = 20

		n = NewNode(nil)
	)

	n.Logging = true

	go n.Run(ctx)
	go n.LogErrors(ctx)

	go func() {
		for _, filename := range filenames {
			log.Printf("filename: %s", filename)

			bs, err := ioutil.ReadFile("../data/planetlabs/" + filename)
			if err != nil {
				log.Printf("error reading %s", filename)
				continue
			}

			r := bytes.NewReader(bs)

			sats := make([]*PubTLE, 0, 1024)

			f := func(i int, line0 string, p prop.Propagator) error {
				if !tle.Legit(p) {
					return nil
				}
				sat := &PubTLE{
					Publisher: "test",
					TLE:       p.(*tle.SGP4TLE),
				}
				sats = append(sats, sat)
				return nil
			}

			if err := tle.DoTLEs(bufio.NewReader(r), tle.NewSGP4TLE, f); err != nil {
				log.Printf("error loading TLEs: %v", err)
				continue
			}

			log.Printf("loading %d", len(sats))

			n.In <- sats
			time.Sleep(time.Second)
		}
	}()

	{
		// Receive reports and count them.
		//
		// This code will run until it gets what it wants.

		count := 1
		for {
			select {
			case <-ctx.Done():
				return
			case rs := <-n.Out:
				for _, r := range rs {
					fmt.Printf("%d %s\n", count, JSON(r))
					count++
					if want <= count {
						return
					}
				}
			}
		}
	}
}
