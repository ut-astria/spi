package node

import (
	"github.com/ut-astria/spi/index"
)

// Consolidate eliminates obsolete outputs.
//
// We could have X novel from one ii and X canceled from another ii.
// Say an ii for A results in AB novel and AC canceled, and a
// subsequent ii for C results in AC novel.  We want to remove AC
// novel and AC canceled.
func (n *Node) Consolidate(ios []*IndexOutput) []*IndexOutput {
	acc := make([]*IndexOutput, 0, len(ios))

	for _, io := range ios {
		o := &IndexOutput{
			Time:     io.Time,
			Novel:    make([]index.Conj, 0, len(io.Novel)),
			Canceled: make([]index.Conj, 0, len(io.Canceled)),
		}

		ns := make(map[index.Conj]bool, len(io.Novel))
		for _, n := range io.Novel {
			ns[n] = true
		}
		for _, c := range io.Canceled {
			if _, have := ns[c]; have {
				delete(ns, c)
			} else {
				o.Canceled = append(o.Canceled, c)
			}
		}
		for n, _ := range ns {
			o.Novel = append(o.Novel, n)
		}

		if 0 < len(o.Novel) || 0 < len(o.Canceled) {
			acc = append(acc, o)
		}
	}

	return acc
}
