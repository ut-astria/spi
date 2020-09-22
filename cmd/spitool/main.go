// Package main is a command-line tool for working with TLEs.
//
// Doesn't do any I/O with CSS.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ut-astria/spi/node"
	"github.com/ut-astria/spi/prop"
	"github.com/ut-astria/spi/tle"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {

	usage := func() string {
		return `Usage: csv|vsc|new|old|elements|sample|tag|prop|plot

csv: TLE to CSV
vsc: CSV to TLE
new: emit new TLEs
old: emit old TLEs
elements: some elements as JSON
sample: sample TLEs
tag: add tag to TLE line0
prop: propagate (SGP4)
plot: generate a crude PNG of reports
`
	}

	if len(os.Args) < 2 {
		log.Fatal(usage())
	}

	var (
		defaultFile = "-"
		siz         = 128 * 1024
		cmd         = os.Args[1]
		args        = os.Args[2:]
	)

	switch cmd {
	case "new", "old":
		// Given some old TLEs and some possibly new TLEs,
		// output just the truly new TLEs as CSV.  Doesn't
		// consider TLE times (ToDo).  Just looks for strings
		// that it hasn't seen before.
		//
		// Could have implemented this function with Linux
		// utilities.
		var (
			fs = flag.NewFlagSet("new", flag.PanicOnError)

			prevFile = fs.String("prev", defaultFile, "Previous TLE input filename")
			newFile  = fs.String("new", defaultFile, "New TLE input filename")
			csvOut   = fs.Bool("csv", false, "CSV output")

			older = make([]string, 0, siz)
			newer = make([]string, 0, siz)
		)

		fs.Parse(args)

		err := doTLEs(*prevFile, func(i int, o *tle.SGP4TLE) error {
			csv := strings.Join(o.TLE, ",")
			older = append(older, csv)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Read %d older TLEs in %s.", len(older), *prevFile)

		err = doTLEs(*newFile, func(i int, o *tle.SGP4TLE) error {
			csv := strings.Join(o.TLE, ",")
			newer = append(newer, csv)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Read %d newer TLEs in %s.", len(newer), *newFile)

		sort.Strings(older)

		var (
			novel = make([]string, 0, siz)
			old   = make([]string, 0, siz)
		)

		for _, csv := range newer {
			i := sort.SearchStrings(older, csv)
			if len(older) == i || older[i] != csv {
				novel = append(novel, csv)
			} else {
				old = append(old, csv)
			}
		}

		log.Printf("Found %d novel and %d old TLEs", len(novel), len(old))

		sort.Strings(novel)
		emit := novel
		if cmd == "old" {
			emit = old
		}
		for _, csv := range emit {
			if *csvOut {
				for _, line := range strings.Split(csv, ",") {
					fmt.Printf("%s\n", line)
				}
			} else {
				fmt.Printf("%s\n", csv)
			}
		}

	case "sample":
		var (
			fs     = flag.NewFlagSet("sample", flag.PanicOnError)
			inFile = fs.String("in", defaultFile, "TLE input filename")
			mod    = fs.Int("mod", 10, "Hash modulus")
			rem    = fs.Int("rem", 0, "Hash remainder")
		)

		fs.Parse(args)

		var r io.Reader
		var err error
		if *inFile == "-" {
			r = os.Stdin
		} else {
			r, err = os.Open(*inFile)
		}
		if err != nil {
			log.Fatal(err)
		}

		err = tle.DoTLEs(bufio.NewReader(r), nil, func(i int, line0 string, p prop.Propagator) error {
			var (
				h = Hash(line0)
				r = h % uint64(*mod)
			)
			if r != uint64(*rem) {
				return nil
			}
			o := p.(*tle.SGP4TLE)
			for _, line := range o.TLE {
				fmt.Printf("%s\n", line)
			}

			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

	case "tag":
		var (
			fs     = flag.NewFlagSet("sample", flag.PanicOnError)
			inFile = fs.String("in", defaultFile, "TLE input filename")
			tag    = fs.String("tag", "", "Tag for line0")
		)

		fs.Parse(args)

		if *tag == "" {
			panic("Need a -tag")
		}

		var r io.Reader
		var err error
		if *inFile == "-" {
			r = os.Stdin
		} else {
			r, err = os.Open(*inFile)
		}
		if err != nil {
			log.Fatal(err)
		}

		err = tle.DoTLEs(bufio.NewReader(r), nil, func(i int, line0 string, p prop.Propagator) error {
			var (
				n = len(line0)
			)
			line0 = strings.TrimSpace(line0) + " " + *tag
			for len(line0) < n {
				line0 += " "
			}
			o := p.(*tle.SGP4TLE)
			o.TLE[0] = line0
			for _, line := range o.TLE {
				fmt.Printf("%s\n", line)
			}

			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

	case "elements":
		// Extract some elements as JSON.
		panic("not implemented (for Vallado)")

		var (
			fs     = flag.NewFlagSet("elements", flag.PanicOnError)
			inFile = fs.String("in", defaultFile, "TLE input filename")
		)

		fs.Parse(args)

		var r io.Reader
		var err error
		if *inFile == "-" {
			r = os.Stdin
		} else {
			r, err = os.Open(*inFile)
		}
		if err != nil {
			log.Fatal(err)
		}

		err = tle.DoTLEs(bufio.NewReader(r), nil, func(i int, line0 string, p prop.Propagator) error {
			// m := make(map[string]interface{})
			// o := p.(*tle.SGP4TLE)
			// es := o.Sat.Elements()
			// m["elements"] = es
			// m["tle"] = o.TLE
			// js, err := json.Marshal(&m)
			// if err != nil {
			// 	log.Fatalf("json.Marshal %s on %#v", err, m)
			// }
			// fmt.Printf("%s\n", js)
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

	case "csv":
		// TLEs in and CSV representation of TLEs out.
		var (
			fs     = flag.NewFlagSet("csv", flag.PanicOnError)
			inFile = fs.String("in", defaultFile, "TLE input filename")
		)

		fs.Parse(args)

		var r io.Reader
		var err error
		if *inFile == "-" {
			r = os.Stdin
		} else {
			r, err = os.Open(*inFile)
		}
		if err != nil {
			log.Fatal(err)
		}

		err = tle.DoTLEs(bufio.NewReader(r), nil, func(i int, line0 string, p prop.Propagator) error {
			o := p.(*tle.SGP4TLE)
			csv := strings.Join(o.TLE, `","`)
			fmt.Printf("\"%s\"\n", csv)
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

	case "prop":
		// Propagate
		var (
			fs       = flag.NewFlagSet("prop", flag.PanicOnError)
			inFile   = fs.String("in", defaultFile, "TLE input filename")
			from     = fs.String("from", "", "Start time")
			duration = fs.Duration("horizon", 600*time.Second, "Duration")
			interval = fs.Duration("interval", 20*time.Second, "Interval")
			// vallado  = flag.Bool("vallado", true, "Use Vallado SGP4 implementation")
		)

		fs.Parse(args)

		var r io.Reader
		var err error
		if *inFile == "-" {
			r = os.Stdin
		} else {
			r, err = os.Open(*inFile)
		}
		if err != nil {
			log.Fatal(err)
		}

		if *from == "" {
			*from = time.Now().UTC().Format(time.RFC3339)
		}
		now, err := time.Parse(time.RFC3339, *from)
		if err != nil {
			log.Fatalf("Bad 'from': %s %s", *from, err)
		}
		then := now.Add(*duration)

		err = tle.DoTLEs(bufio.NewReader(r), nil, func(i int, line0 string, p prop.Propagator) error {
			o := p.(*tle.SGP4TLE)
			for t := now; t.Before(then); t = t.Add(*interval) {
				e, err := o.Prop(t)
				if err != nil {
					return err
				}

				lla, err := node.ECIToLLA(t, e.ECI)
				if err != nil {
					return err
				}

				m := map[string]interface{}{
					"At":    t,
					"State": e,
					"TLE":   o.TLE,
					"LLA":   lla,
					"Age":   o.ApproxAge(t).Seconds(),
				}
				js, err := json.Marshal(&m)
				if err != nil {
					log.Fatalf("prop json.Marshal error %s on %#v", err, m)
				}
				fmt.Printf("%s\n", js)
			}
			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

	case "vsc":
		// CSV representation of TLEs in and TLEs out.
		var (
			fs     = flag.NewFlagSet("vcs", flag.PanicOnError)
			inFile = fs.String("in", defaultFile, "TLE CSV input filename")
		)

		fs.Parse(args)

		var bs []byte
		var err error

		if *inFile == "-" {
			bs, err = ioutil.ReadAll(os.Stdin)
		} else {
			bs, err = ioutil.ReadFile(*inFile)
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, line := range bytes.Split(bs, []byte("\n")) {
			for _, s := range strings.Split(string(line), ",") {
				s = strings.Trim(s, `"`)
				fmt.Printf("%s\r\n", s)
			}
		}

	case "plot":
		// CSV representation of TLEs in and TLEs out.
		var (
			fs      = flag.NewFlagSet("plot", flag.PanicOnError)
			inFile  = fs.String("in", defaultFile, "reports (JSON) input filename")
			outFile = fs.String("out", "reports.png", "filename for PNG")
		)

		fs.Parse(args)

		var bs []byte
		var err error

		if *inFile == "-" {
			bs, err = ioutil.ReadAll(os.Stdin)
		} else {
			bs, err = ioutil.ReadFile(*inFile)
		}
		if err != nil {
			log.Fatal(err)
		}

		type Point struct {
			T time.Time
			D float32
		}

		var (
			t0    time.Time
			t1    time.Time
			rs    = make([]node.Report, 0, 1024*16)
			pairs = make(map[string][]Point)
		)

		for i, bs := range bytes.Split(bs, []byte("\n")) {
			if len(bs) == 0 {
				continue
			}
			var r node.Report
			if err = json.Unmarshal(bs, &r); err != nil {
				log.Fatalf("error parsing line %d: %s: '%s'", i, err, bs)
			}
			if t0.IsZero() || r.At.Before(t0) {
				t0 = r.At
			}
			if r.At.After(t1) {
				t1 = r.At
			}
			pair := r.Objs[0].Name + "/" + r.Objs[1].Name
			ps, have := pairs[pair]
			if !have {
				ps = make([]Point, 0, 32)
			}
			ps = append(ps, Point{
				T: r.At,
				D: r.Dist,
			})
			pairs[pair] = ps

			rs = append(rs, r)
		}

		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}

		p.Title.Text = "Reports"
		p.X.Label.Text = fmt.Sprintf("Seconds since %s", t0.Format(time.RFC3339Nano))
		p.Y.Label.Text = "Distance (km)"
		p.Y.Min = 0

		x := func(t time.Time) float64 {
			return t.Sub(t0).Seconds()

		}

		xyss := make([]interface{}, 0, len(pairs))

		for pair, ps := range pairs {
			xys := make(plotter.XYs, len(ps))
			for i, p := range ps {
				xys[i].X = x(p.T)
				xys[i].Y = float64(p.D)
			}
			xyss = append(xyss, pair)
			xyss = append(xyss, xys)
		}

		if err = plotutil.AddLinePoints(p, xyss...); err != nil {
			panic(err)
		}

		if 10 < len(pairs) {
			// Don't attempt a legend.
			nope, _ := plot.NewLegend()
			p.Legend = nope
		}

		// Save the plot to a PNG file.
		if err := p.Save(8*vg.Inch, 4*vg.Inch, *outFile); err != nil {
			panic(err)
		}

	default:
		fmt.Printf("%s\n", usage())
		os.Exit(1)
	}
}

func doTLEs(filename string, f func(i int, o *tle.SGP4TLE) error) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}

	err = tle.DoTLEs(bufio.NewReader(r), nil, func(i int, line1 string, p prop.Propagator) error {
		return f(i, p.(*tle.SGP4TLE))
	})

	return err
}

func Hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}
