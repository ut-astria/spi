package index

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/ut-astria/spi/prop"
	"github.com/ut-astria/spi/tle"
)

var (
	// Filenames are used in performance tests.
	//
	// Tests look for thes files in FilenameDir.
	Filenames = []string{
		"2020-04-10.tle",
		"2020-04-11.tle",
		"2020-04-12.tle",
		"2020-04-13.tle",
		"2020-04-14.tle",
		"2020-04-15.tle",
		"2020-04-16.tle",
		"2020-04-17.tle",
		"2020-04-18.tle",
		"2020-04-19.tle",
		"2020-04-20.tle",
	}

	// FilenameDir should name the directory that contains
	// Filenames.
	FilenameDir = "../data/astria"
)

// asKey makes a Key with the given catalog number (and zero for the
// publisher).
func asKey(n int) Key {
	return Key{
		CatalogNum: CatalogNum(n),
	}
}

func TestCell(t *testing.T) {
	var (
		id Id = 1

		key Key

		pos = Pos{
			X: 10,
			Y: 11,
			Z: 12,
		}

		pps = func() []ProbPos {
			pp := ProbPos{
				pos,
			}
			fmt.Printf("update %#v\n", pp)
			return []ProbPos{pp}
		}

		i = NewIndex(7, 10)

		label string
	)

	check := func(nc, nn int, canceledCs, novelCs []Conj, err error) {
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("\n===> %s\n", label)
		PrintCs(canceledCs, "  can ", "      ")
		PrintCs(novelCs, "  nov ", "      ")
		fmt.Printf("\n")
		i.Print()
		if n := len(canceledCs); n != nc {
			t.Fatalf("%s; can: %d", label, n)
		}
		if n := len(novelCs); n != nn {
			t.Fatalf("%s; nov: %d", label, n)
		}
		fmt.Printf("<===\n\n")
	}

	label = "first"
	fmt.Printf("*** %s\n\n", label)
	id++
	key = asKey(1)
	c, n, _, err := i.Update(id, key, pps())
	check(0, 0, c, n, err)

	label = "second: get a report"
	fmt.Printf("*** %s\n\n", label)
	id++
	key = asKey(2)
	pos.X = 10.0001
	c, n, _, err = i.Update(id, key, pps())
	check(0, 1, c, n, err)

	label = "move first: get a can and a nov"
	fmt.Printf("*** %s\n\n", label)
	id++
	key = asKey(1)
	pos.X = 10.0002
	c, n, _, err = i.Update(id, key, pps())
	check(1, 1, c, n, err)

	label = "update first: get a can and nov"
	fmt.Printf("*** %s\n\n", label)
	id++
	key = asKey(1)
	c, n, _, err = i.Update(id, key, pps())
	check(1, 1, c, n, err)

	label = "add third: get nothing"
	fmt.Printf("*** %s\n\n", label)
	id++
	key = asKey(3)
	pos.X = 20.0002
	c, n, _, err = i.Update(id, key, pps())
	check(0, 0, c, n, err)

	label = "move third: get two new"
	fmt.Printf("*** %s\n\n", label)
	id++
	key = asKey(3)
	pos.X = 10.0003
	c, n, _, err = i.Update(id, key, pps())
	check(0, 2, c, n, err)

	label = "move first again: get two cans and novs"
	fmt.Printf("*** %s\n\n", label)
	id++
	key = asKey(1)
	pos.X = 10.0003
	c, n, _, err = i.Update(id, key, pps())
	check(2, 2, c, n, err)
}

func BenchmarkIndex(b *testing.B) {

	log.Printf("BenchmarkIndex %d", b.N)

	filenames := []string{
		"2020-04-10.tle",
		"2020-04-11.tle",
		"2020-04-12.tle",
		"2020-04-13.tle",
		"2020-04-14.tle",
		"2020-04-15.tle",
		"2020-04-16.tle",
		"2020-04-17.tle",
		"2020-04-18.tle",
		"2020-04-19.tle",
		"2020-04-20.tle",
	}

	now := time.Now().UTC()

	type input struct {
		Id  Id
		Key Key
		PPS []ProbPos
	}

	inputss := make([][]input, 0, len(filenames))

	for _, filename := range filenames {
		var (
			is = make([]input, 0, 32*1024)
			id = 0
		)

		bs, err := ioutil.ReadFile(FilenameDir + "/" + filename)
		if err != nil {
			b.Skip(err)
		}

		r := bytes.NewReader(bs)
		f := func(i int, line0 string, p prop.Propagator) error {
			e, err := p.Prop(now)
			if err != nil {
				return nil
			}
			var (
				key = Key{
					CatalogNum: CatalogNum(i),
				}

				pos = Pos{
					X: e.ECI.X,
					Y: e.ECI.Y,
					Z: e.ECI.Z,
				}
				pps = []ProbPos{
					ProbPos{
						pos,
					},
				}
			)
			is = append(is, input{
				Id:  Id(id),
				Key: key,
				PPS: pps,
			})

			id++

			return nil
		}

		if err := tle.DoTLEs(bufio.NewReader(r), tle.NewSGP4TLE, f); err != nil {
			b.Fatal(err)
		}

		inputss = append(inputss, is)
	}

	var (
		index      = NewIndex(5, 50)
		i          = 0
		cans, novs int
		limit      = b.N * 1000
	)

	b.ResetTimer()

LOOP:
	for {
		for _, is := range inputss {
			for _, input := range is {
				if i == limit {
					break LOOP
				}
				can, nov, _, err := index.Update(input.Id, input.Key, input.PPS)
				if err != nil {
					b.Fatal(err)
				}
				cans += len(can)
				novs += len(nov)
				i++
			}
		}
	}

	log.Printf("b.N=%d %d %d %d", b.N, i, cans, novs)
}

func TestSpeed(t *testing.T) {

	now := time.Now().UTC()

	type input struct {
		Id  Id
		Key Key
		PPS []ProbPos
	}

	inputss := make([][]input, 0, len(Filenames))

	for _, filename := range Filenames {
		now = now.Add(60 * time.Second)

		var (
			is = make([]input, 0, 32*1024)
			id = 0
		)

		filename := FilenameDir + "/" + filename
		bs, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("TestSpeed: skipping (%s: %s)", filename, err)
			t.Skipf("TestSpeed skipping %s: %s", filename, err)
			continue
		}

		r := bytes.NewReader(bs)
		f := func(i int, line0 string, p prop.Propagator) error {
			e, err := p.Prop(now)
			if err != nil {
				return nil
			}
			var (
				key = Key{
					CatalogNum: CatalogNum(i),
				}

				pos = Pos{
					X: e.ECI.X,
					Y: e.ECI.Y,
					Z: e.ECI.Z,
				}
				pps = []ProbPos{
					ProbPos{
						pos,
					},
				}
			)
			is = append(is, input{
				Id:  Id(id),
				Key: key,
				PPS: pps,
			})

			id++

			return nil
		}

		if err := tle.DoTLEs(bufio.NewReader(r), tle.NewSGP4TLE, f); err != nil {
			t.Fatal(err)
		}

		inputss = append(inputss, is)
	}

	var (
		index      = NewIndex(5, 50)
		i          = 0
		cans, novs int
		limit      = 1000 * 100
		then       = time.Now()
	)

LOOP:
	for {
		for _, is := range inputss {
			for _, input := range is {
				if i == limit {
					break LOOP
				}
				can, nov, _, err := index.Update(input.Id, input.Key, input.PPS)
				if err != nil {
					t.Fatal(err)
				}
				cans += len(can)
				novs += len(nov)
				i++
			}
		}
	}

	elapsed := time.Now().Sub(then).Seconds()
	secsPer := elapsed / float64(limit)
	micsPer := secsPer * 1000 * 1000
	log.Printf("b.N=%d %d %d %d %v mics/op", limit, i, cans, novs, micsPer)
}
