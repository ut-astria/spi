package node

import (
	"context"
	"math/rand"
	"sync"

	"github.com/ut-astria/spi/index"
)

type Intern struct {
	m   map[string]index.Id
	inv map[index.Id]string
}

func NewIntern() *Intern {
	return &Intern{
		m:   make(map[string]index.Id),
		inv: make(map[index.Id]string),
	}
}

// Count returns the number of interned strings.
//
// This method is not safe for concurrent use.
func (is *Intern) Count() int {
	return len(is.m)
}

// Rem removes the string (if any) with the given id.
//
// This method is not safe for concurrent use.
func (is *Intern) Rem(id index.Id) bool {
	s, have := is.inv[id]
	if have {
		delete(is.inv, id)
		delete(is.m, s)
	}
	return have
}

// Get finds the id for the give string or returns false if the string
// isn't interned.
//
// This method is not safe for concurrent use.
func (is *Intern) Get(s string) (index.Id, bool) {
	id, have := is.m[s]
	return id, have
}

// Find returns the string for the given id or returns false if the id
// wasn't found.
//
// This method is not safe for concurrent use.
func (is *Intern) Find(id index.Id) (string, bool) {
	s, have := is.inv[id]
	return s, have
}

func (is *Intern) genId() index.Id {
	for {
		// ToDo: Limit!

		n := index.Id(rand.Uint32())
		// We hope that rand.Uint32() has sufficient entropy.
		// Since we don't know, let's scan a little, too.
		for i := index.Id(0); i < 128; i++ {
			id := n + i
			// Overflow okay.
			if _, have := is.inv[id]; !have {
				return id
			}
		}
	}

	return 0
}

// Intern interns the given string and returns its id and whether the
// string was previously known.
//
// This method is not safe for concurrent use.
func (is *Intern) Intern(k string) (index.Id, bool) {
	id, have := is.m[k]
	if !have {
		id = is.genId()
		is.m[k] = id
		is.inv[id] = k
	}
	return id, have
}

type Interns struct {
	Ids   *InternTLE
	Keys  *Intern
	KeyId map[index.Key]index.Id

	sync.RWMutex
}

func NewInterns() *Interns {
	return &Interns{
		Ids:   NewInternTLE(),
		Keys:  NewIntern(),
		KeyId: make(map[index.Key]index.Id),
	}
}

func (is *Interns) RExec(ctx context.Context, f func(*Interns) error) error {
	is.RLock()
	err := f(is)
	is.RUnlock()
	return err
}

func (is *Interns) Exec(ctx context.Context, f func(*Interns) error) error {
	is.Lock()
	err := f(is)
	is.Unlock()
	return err
}

func (is *Interns) IdCount() int {
	is.RLock()
	n := len(is.Ids.m)
	is.RUnlock()
	return n
}

func (is *Interns) Update(k index.Key, id index.Id) {
	if old, have := is.KeyId[k]; have {
		if old != id {
			is.Ids.Rem(old)
		}
	}
	is.KeyId[k] = id
}
