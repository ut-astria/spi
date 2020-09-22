package index

import (
	"fmt"
	"log"

	"github.com/golang/geo/r3"
	"github.com/golang/geo/s2"
)

type CellId s2.CellID

// CellFinder maps a position to a cell.
type CellFinder struct {
	Level *Level
}

func NewCellFinder(level int) *CellFinder {
	return &CellFinder{
		Level: NewLevel(level),
	}
}

// Cell returns the cell that contains the given position.
func (i *CellFinder) Find(p Pos) (CellId, error) {

	ll := s2.LatLngFromPoint(s2.Point{
		Vector: r3.Vector{
			X: float64(p.X),
			Y: float64(p.Y),
			Z: float64(p.Z),
		},
	})

	c, err := i.Level.Find(ll)
	if err != nil {
		return 0, err
	}

	return CellId(*c), nil
}

// Neighbors returns the cells that neighbor the cell that contains
// the given position.
func (i *CellFinder) Neighbors(cid CellId) []s2.CellID {
	// Don't want to allocate a new array to shift types.
	return i.Level.Neighbors(s2.CellID(cid))
}

type Level struct {
	Level   int
	Coverer s2.RegionCoverer
}

// NewLevel generates a Level with a Coverer with min and max levels
// both at given level.
func NewLevel(level int) *Level {
	l := &Level{
		Level: level,
		Coverer: s2.RegionCoverer{
			MinLevel: level,
			MaxLevel: level,
		},
	}

	return l
}

// Find hopefully returns the cell that contains the given point.
func (l *Level) Find(ll s2.LatLng) (*s2.CellID, error) {
	p := s2.PointFromLatLng(ll)
	union := l.Coverer.CellUnion(p)
	switch len(union) {
	case 0:
		return nil, fmt.Errorf("bad union for point %s", ll)
	case 1:
		return &union[0], nil
	default:
		// https://github.com/golang/geo/issues/58
		for _, c := range union {
			cell := s2.CellFromCellID(c)
			loop := s2.LoopFromCell(cell)
			if loop.ContainsPoint(p) {
				return &c, nil
			}
		}
		log.Printf("warning %#v -> covering had no covering cell %#v", ll, union)
		return nil, fmt.Errorf("bad union for point %s", ll)
	}
}

// Neighbors returns the given cell's neighors (via
// Cell.AllNeighbors().
func (l *Level) Neighbors(c s2.CellID) []s2.CellID {
	return c.AllNeighbors(l.Level)
}
