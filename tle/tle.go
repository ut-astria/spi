package tle

import (
	"bufio"
	"io"
	"strings"

	"github.com/ut-astria/spi/prop"
)

// DoTLEs parses Sats from the a reader that offers TLEs.
func DoTLEs(r *bufio.Reader, parser func(line0, line1, line2 string) (prop.Propagator, error), f func(i int, line0 string, p prop.Propagator) error) error {

	if parser == nil {
		parser = NewSGP4TLE
	}

	var (
		tle []string
		i   = 0
	)

	for {
		line, err := r.ReadString('\n')
		if 0 < len(line) {
			if tle == nil {
				tle = make([]string, 3)
			}
			tle[i%3] = line
			if (i+1)%3 == 0 {
				p, err := parser(tle[0], tle[1], tle[2])
				if err != nil {
					return err
				}
				if err = f((i-2)/3, tle[0], p); err != nil {
					return err
				}
				tle = nil
			}
			i++
		}
		if err == io.EOF {
			break
		}
	}

	return nil
}

var ObjTypes = map[string]string{
	"COOLANT":          "debris",
	"DEB":              "debris",
	"DEBRIS":           "debris",
	"SHROUD":           "debris",
	"WESTFORD NEEDLES": "debris",
	"R/B":              "rocket",
	"AKM":              "rocket",
	"PKM":              "rocket",
}

func GetType(line0 string) string {
	name := strings.TrimSpace(line0)
	if len(name) == 0 {
		return "unknown"
	}
	for tag, typ := range ObjTypes {
		if strings.Contains(name, tag) {
			return typ
		}
	}
	return "payload"
}
