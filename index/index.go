package index

import (
	"fmt"
	"math"

	"github.com/golang/geo/s2"
)

// Id represents the distinct report in its totality (as far as this
// package goes).  Two reports that are "equal" should have the same
// id.
//
// We are using only four bytes to conserve RAM.
//
// Re "Id" vs "ID": I'm in the "Id" is an abbreviation, not an acronym
// camp, and I've been there (rightly or wrongly) for a long time.
// https://github.com/golang/lint/issues/89 says I'm wrong -- at least
// as far as Go sources go.  Had I known at the time, I probably would
// have complied.
type Id uint32

// CatalogNum is shared identity of the focal object.
//
// This whole concept is perilous, but we'll attempt to use it as best
// we can.
//
// Currently these identifiers are assigned by the United States Space
// Command; however, we plan to be able to extend this namespace at
// least a little (within the current confines of four bytes, which
// could of course change in the future).
//
// A publisher of an unassigned object needs to take care to use a
// number that won't conflict with use of that number by a different
// publisher!  Yes, this approach will become problematic, but we need
// this notion here in order to avoid reporting the conjuction of two
// views of the same thing.
//
// The number zero represents an unknown value.  In this case, any
// other position report is a candidate for an emitted conjunction.
// With that exception, no conjunction report is emitted based on two
// reports that have the same CatalogNum.
type CatalogNum uint32

// Publisher is the id for the publisher of a position.
type Publisher uint32

// Key is the key for storing positions in the index.
//
// In an Index, a CatalogNum+Publisher has can have a single opinion
// of that object's position(s).  Note that that opinion can include
// multiple possible positions.
type Key struct {
	CatalogNum
	Publisher
}

func (k Key) String() string {
	return fmt.Sprintf("K(cat=%d,pub=%d)", k.CatalogNum, k.Publisher)
}

// Pos is the object's state.
//
// Currently this state only includes position (ECI).
//
// If we add velocity, then the index could sample Conjs with low
// relative speeds.  For objects parked together, every index would
// emit an Conj.  The resulting downstream relative computational load
// is unknown but possibly significant.  However, most of that load is
// (I think) typically spread across all cores.  Adding velocity
// results in a memory increase of maybe 40%.  The guiding principle
// was to optimize for memory, so we'll stick to that principle for
// now.  ToDo: Reconsider.
type Pos struct {
	X float32
	Y float32
	Z float32

	// Experiment: include velocity
	// VX, VY, VZ float32
}

func (p Pos) String() string {
	return fmt.Sprintf("Pos(%0.4f,%0.4f,%0.4f)", p.X, p.Y, p.Z)
}

// Dist computes the Cartesian distance.
func (p Pos) Dist(q Pos) float32 {
	var (
		x = p.X - q.X
		y = p.Y - q.Y
		z = p.Z - q.Z
	)
	return float32(math.Sqrt(float64(x*x + y*y + z*z)))
}

// Prob represents uncertainty (without any index-imposed meaning).
//
// Can update this definition to support more dimensions, but of
// course we're mindful of memory consumption.
type Prob float32

// ProbPos represents one candidate Pos.
//
// Currently we do not include Prob here, but this type exists to make
// it easy to add that and other data.
type ProbPos struct {
	Pos
	// Prob
}

// At represents an opinion about an object's position.
type At struct {
	Id
	ProbPos
}

// Conj represents an opinion about the proximity of two objects (at
// some unspecified time).
type Conj struct {
	Ats  [2]At
	Dist float32
}

// AtLess specifies the canonical order on At instances.
//
// We make the computation completely explicit rather than rely on a
// more indirect algorithm.
func AtLess(a, b At) bool {
	if a.Id < b.Id {
		return true
	}
	if a.Id > b.Id {
		return false
	}

	if a.X < b.X {
		return true
	}
	if a.X > b.X {
		return false
	}
	if a.Y < b.Y {
		return true
	}
	if a.Y > b.Y {
		return false
	}
	if a.Z < b.Z {
		return true
	}
	if a.Z > b.Z {
		return false
	}

	if a.Id < b.Id {
		return true
	}
	if a.Id > b.Id {
		return false
	}

	if a.X < b.X {
		return true
	}
	if a.X > b.X {
		return false
	}
	if a.Y < b.Y {
		return true
	}
	if a.Y > b.Y {
		return false
	}
	if a.Z < b.Z {
		return true
	}
	if a.Z > b.Z {
		return false
	}

	return false
}

func NewConj(a, b At, d float32) Conj {
	if AtLess(b, a) {
		a, b = b, a
	}
	c := Conj{
		Ats:  [2]At{a, b},
		Dist: d,
	}

	return c
}

type IdProbPos struct {
	Id
	CatalogNum
	ProbPos
}

func (i IdProbPos) String() string {
	return fmt.Sprintf("IdProbPos(%d,%d,%v)", i.Id, i.CatalogNum, i.ProbPos)
}

// Cell represents a set of IdProbPos in an S2 cell.
type Cell struct {
	items []IdProbPos
}

func (c Cell) String() string {
	return fmt.Sprintf("Cell(%d:%v)", len(c.items), c.items)
}

func NewCell() *Cell {
	return &Cell{
		items: make([]IdProbPos, 0, 2),
	}
}

func (c *Cell) Count() int {
	return len(c.items)
}

// IdProbPoss is a set of IdProbPos for a given object.
type IdProbPoss struct {
	Id
	CatalogNum
	PPS []ProbPos
}

// Index receives object positions and returns "conjunctions".
//
// This implementation is not safe for concurrent use.
type Index struct {
	Dist       float32
	CellFinder *CellFinder

	IPPS  map[Key]*IdProbPoss
	Cells map[CellId]*Cell
}

// NewIndex creates a new index based on cells with the given level
// and uses the given distance as the maximum distance for
// "conjunction" determination.
func NewIndex(level int, dist float32) *Index {
	return &Index{
		Dist:       dist,
		IPPS:       make(map[Key]*IdProbPoss),
		CellFinder: NewCellFinder(level),
		Cells:      make(map[CellId]*Cell),
	}
}

// NewIndexWithFinder allows the caller to provide a CellFinder that
// can be used for multiple indexes.
func NewIndexWithFinder(finder *CellFinder, dist float32) *Index {
	return &Index{
		Dist:       dist,
		IPPS:       make(map[Key]*IdProbPoss),
		CellFinder: finder,
		Cells:      make(map[CellId]*Cell),
	}
}

// Search is the core method for finding Conjs.
//
// This implementation is not safe for concurrent use.
func (i *Index) Search(cid CellId, spp IdProbPos, d float32) []Conj {

	var (
		neighbors = append(i.CellFinder.Neighbors(cid), s2.CellID(cid))
		cells     = make([]*Cell, 0, len(neighbors)+1)
	)

	for _, n := range neighbors {
		id := CellId(n)
		cell, have := i.Cells[id]
		if !have {
			continue
		}
		cells = append(cells, cell)
	}

	cs := make([]Conj, 0, 2)
	for _, c := range cells {
		for _, spp0 := range c.items {
			if spp0 == spp {
				continue
			}
			if spp0.CatalogNum == spp.CatalogNum {
				if spp0.CatalogNum != 0 && spp.CatalogNum != 0 {
					continue
				}
			}
			d0 := spp0.Dist(spp.ProbPos.Pos)
			if d0 <= d {
				c := NewConj(
					At{
						Id:      spp0.Id,
						ProbPos: spp0.ProbPos,
					},
					At{
						Id:      spp.Id,
						ProbPos: spp.ProbPos,
					},
					d0)
				cs = append(cs, c)
			}
		}
	}
	return cs
}

func (c *Cell) Add(ipp IdProbPos) {
	c.items = append(c.items, ipp)
}

func (c *Cell) Rem(ipp IdProbPos) {
	n := len(c.items)
	for i, ipp0 := range c.items {
		if ipp == ipp0 {
			c.items[i] = c.items[n-1]
			c.items = c.items[0 : n-1]
			break
		}
	}
}

// IndexData represents some statistics about an Index.
type IndexData struct {
	CellCount   int
	KeyCount    int
	KeysPerCell float64
}

func (i *Index) Data() *IndexData {
	d := &IndexData{
		CellCount: len(i.Cells),
		KeyCount:  len(i.IPPS),
	}
	d.KeysPerCell = float64(d.KeyCount) / float64(d.CellCount)

	return d
}

func (i *Index) Print() {
	fmt.Printf("Cells: %d\n", len(i.Cells))
	for id, c := range i.Cells {
		fmt.Printf("  %v\n", id)
		for _, spp := range c.items {
			fmt.Printf("    %v\n", spp)
		}
	}
	fmt.Printf("PPS: %d\n", len(i.IPPS))
	for key, pps := range i.IPPS {
		fmt.Printf("  %v id=%v\n", key, pps.Id)
		for _, pp := range pps.PPS {
			fmt.Printf("    %v\n", pp)
		}
	}
}

func (i *Index) GetCell(p Pos) (*Cell, CellId, error) {
	id, err := i.CellFinder.Find(p)
	if err != nil {
		return nil, 0, err
	}

	cell, _ := i.Cells[id]

	return cell, id, nil
}

// Update updates the index and returns canceled Conj(s) and new
// Conj(s).
//
// This implementation is not safe for concurrent use.
//
// The key is for the given set of ProbPos, which should be complete
// for the key.  For example, if the key represents the combination of
// a satellite and a publisher of information about that satellite,
// then an Update call with that key should provide the publisher's
// complete knowledge for that key.
//
// The id should represent data that was used to generate this
// knowledge.  These ids are included in the returned Conj, so a
// caller can do additional processing (like scanning for closer
// approaches) after retrieving the source data that an id identifies.
//
// The returned Id (if any) is the id of the previous Update call for
// the given key.
//
// Some of the constants in this method body relate to tuning (such as
// Level).  The coarser the cells, the larger the initial allocations.
// ToDo: Expose these values.
func (i *Index) Update(id Id, key Key, pps []ProbPos) ([]Conj, []Conj, Id, error) {

	// Remove all ps previously associated with key.
	var (
		ipps  = i.IPPS[key]
		oldCs = make([]Conj, 0, 2)
	)

	if ipps == nil {
		ipps = &IdProbPoss{}
	}
	oldId := ipps.Id

	for _, pp := range ipps.PPS {
		cell, cid, err := i.GetCell(pp.Pos)
		if err != nil {
			return nil, nil, oldId, err
		}
		if cell == nil {
			continue
		}
		ipp := IdProbPos{
			Id:         ipps.Id,
			CatalogNum: key.CatalogNum,
			ProbPos:    pp,
		}
		cell.Rem(ipp)
		if 0 == len(cell.items) {
			delete(i.Cells, cid)
		}
		cs := i.Search(cid, ipp, i.Dist)
		oldCs = append(oldCs, cs...)
	}

	// Write pps and gather new Conjs.
	newCs := make([]Conj, 0, 8)
	for _, pp := range pps {
		cell, cid, err := i.GetCell(pp.Pos)
		if err != nil {
			return nil, nil, oldId, err
		}
		if cell == nil {
			cell = NewCell()
			i.Cells[cid] = cell
		}
		ipp := IdProbPos{
			Id:         id,
			CatalogNum: key.CatalogNum,
			ProbPos:    pp,
		}
		cell.Add(ipp)

		cs := i.Search(cid, ipp, i.Dist)
		newCs = append(newCs, cs...)
	}

	// Remember these PPS and the id.
	i.IPPS[key] = &IdProbPoss{
		Id:         id,
		CatalogNum: key.CatalogNum,
		PPS:        pps,
	}

	canceledCs, novelCs := Diff(oldCs, newCs)

	return canceledCs, novelCs, oldId, nil
}

// Diff returns (1) the set of Conjs in oldCs that are not in newCs
// and (2) the set of newCs that are not in oldCs.
func Diff(oldCs, newCs []Conj) ([]Conj, []Conj) {
	var (
		canceledCs = make([]Conj, 0, len(oldCs))
		novelCs    = make([]Conj, 0, len(newCs))
	)

CANCELED:
	for _, c := range oldCs {
		for _, c0 := range newCs {
			if c == c0 {
				continue CANCELED
			}
		}
		canceledCs = append(canceledCs, c)
	}

NOVEL:
	for _, c := range newCs {
		for _, c0 := range oldCs {
			if c == c0 {
				continue NOVEL
			}
		}
		novelCs = append(novelCs, c)
	}

	return canceledCs, novelCs
}

func PrintCs(cs []Conj, prefix0, prefix string) {
	if len(cs) == 0 {
		fmt.Printf("%s[]\n", prefix0)
		return
	}
	for i, c := range cs {
		s := prefix
		if i == 0 {
			s = prefix0
		}
		fmt.Printf("%s%v\n", s, c)
	}
}
