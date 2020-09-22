package node

import (
	"math/rand"

	"github.com/ut-astria/spi/index"
)

type InternTLE struct {
	m   map[string]index.Id
	inv map[index.Id]*PubTLE
}

func NewInternTLE() *InternTLE {
	siz := 64 * 1024
	return &InternTLE{
		m:   make(map[string]index.Id, siz),
		inv: make(map[index.Id]*PubTLE, siz),
	}
}

func (p *PubTLE) Key() string {
	k := p.Publisher + "/"
	lines := p.TLE.TLE
	k += lines[0] + "/"
	k += lines[1] + "/"
	k += lines[2]
	return k
}

func (is *InternTLE) Count() int {
	return len(is.m)
}

func (is *InternTLE) Rem(id index.Id) bool {
	p, have := is.inv[id]
	if have {
		delete(is.inv, id)
		delete(is.m, p.Key())
	}
	return have
}

func (is *InternTLE) Find(id index.Id) (*PubTLE, bool) {
	p, have := is.inv[id]
	return p, have
}

func (is *InternTLE) genId() index.Id {
	for {
		// ToDo: Limit!

		n := index.Id(rand.Uint32())
		// We hope that rand.Uint32() has sufficient entropy.
		// Since we don't know, let's scan a little, too.
		for i := index.Id(0); i < 32; i++ {
			id := n + i
			// Overflow okay.
			if _, have := is.inv[id]; !have {
				return id
			}
		}
	}

	return 0
}

func (is *InternTLE) Intern(p *PubTLE) (index.Id, bool) {
	k := p.Key()
	id, have := is.m[k]
	if !have {
		id = is.genId()
		is.m[k] = id
		is.inv[id] = p
	}
	return id, have
}
