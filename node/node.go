package node

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/ut-astria/spi/index"
	"github.com/ut-astria/spi/misc"
	"github.com/ut-astria/spi/prop"
	"github.com/ut-astria/spi/sgp4"
	"github.com/ut-astria/spi/tle"
)

var DefaultCfg = &Cfg{
	Horizon:             60,
	Resolution:          time.Second,
	ScanDist:            20,
	IndexDist:           50,
	IndexLevel:          5,
	Scan:                true,
	SlowSample:          10,
	SlowSampleThreshold: 0.1,
}

// Cfg is a Node configuration.
type Cfg struct {
	// T0 is the logical starting time for processing.
	T0 time.Time

	// Horizon is the size of the processing horizon in the number
	// of ticks, where a tick has a duration of Resolution.
	Horizon int

	// Resolution is the duration of a tick.
	Resolution time.Duration

	// Scan turns on intra-tick scanning.
	Scan bool

	// ScanDist is scanning's maximum distance for emitting a report.
	ScanDist float32

	// IndexDist is indexing's maximum distance for emitting a report.
	IndexDist float32

	// IndexLevel is the S2 Cell level for all indexes.
	IndexLevel int

	// SlowSample gives a sampling rate when relative speed is
	// less than SlowSampleThreshold.
	SlowSample int

	// SlowSampleThreshold is the maximum relative speed for
	// sampling to occur.
	SlowSampleThreshold float32

	// Logging turns on logging, which uses log.Printf.
	Logging bool
}

type Metrics struct {
	// T is the wall-clock time for this report.
	T time.Time

	// T1 is the virtual time for the Node.
	T1 time.Time

	// In is the total number of TLEs ingested since the Node started.
	In uint64

	// Live is the number of live TLEs.
	Live int

	// Goroutines is the current number of goroutines.
	Goroutines int

	// Lag is difference between when a tick was scheduled and
	// when it was executed.
	//
	// As Lag approaches a the virtual duration of a tick, we see
	// more of a queuing problem.
	Lag time.Duration

	// Strings is the number of interned strings.
	Strings int
}

// PubTLE associates a Publisher with a TLE.
type PubTLE struct {
	// Publisher is an opaque name for the publisher (source) of this TLE.
	Publisher string

	TLE *tle.SGP4TLE
}

func (p *PubTLE) Name() string {
	if p.Publisher == "" {
		return p.TLE.CatNum
	}
	return p.TLE.CatNum + "/" + p.Publisher
}

type Node struct {
	Cfg

	// In receives in-coming TLEs.
	In chan []*PubTLE

	// Out produces the emitted reports.
	Out chan []*Report

	// Errs, if not null, produces (asynchronous) errors.
	Errs chan error

	// Metrics emits Metrics.
	Metrics chan Metrics

	// TimeOffset is the difference between real time and logical time.
	TimeOffset time.Duration

	// Finder finds Cells, and it's used for all Indexes.
	Finder *index.CellFinder

	// All of a node's working state lives in interns and live.

	// interns maps things to ids and back again.
	interns *Interns

	// live the current set of TLEs, which are indexed by their
	// index.Keys.
	live map[index.Key]*IndexInput
}

// NewNode makes a new Node, with cfg defaulting to DefaultCfg.
//
// The Errs channel is nil.
func NewNode(cfg *Cfg) *Node {
	if cfg == nil {
		cfg = DefaultCfg
	}
	return &Node{
		Cfg:     *cfg,
		In:      make(chan []*PubTLE),
		Out:     make(chan []*Report),
		Errs:    nil,
		interns: NewInterns(),
		live:    make(map[index.Key]*IndexInput),
	}
}

// Prepare initializes a few values required before using the Node.
//
// Run calls Prepare if Finder is nil.
func (n *Node) Prepare(ctx context.Context) {
	if n.T0.IsZero() {
		n.T0 = time.Now().UTC()
	}

	n.T0 = RoundTime(n.T0.UTC(), n.Resolution)
	n.TimeOffset = n.T0.Sub(time.Now())
	n.Finder = index.NewCellFinder(n.IndexLevel)
}

// Run executes the main event loop in the current goroutine.
//
// Run calls Prepare if Finder is nil.
//
// This execution should provide backpressure on In.
func (n *Node) Run(ctx context.Context) {

	if n.Finder == nil {
		n.Prepare(ctx)
	}

	var (
		indexes = make(map[time.Time]*Index)
		t0      = RoundTime(time.Now().UTC(), n.Resolution)
		t1      = t0.Add(time.Duration(n.Horizon) * n.Resolution)
		ticker  = time.NewTicker(n.Resolution)

		inCount  = uint64(0)
		inCount0 = inCount
	)

	for t := t0; t.Before(t1); t = t.Add(n.Resolution) {
		i := n.NewIndex(t)
		indexes[t] = i
		go i.Run(ctx)
	}

LOOP:
	for {
		n.logf(ctx, "Node listening (ids:%d, indexes:%d)",
			n.interns.IdCount(),
			len(indexes))
		select {
		case <-ctx.Done():
			break LOOP
		case t := <-ticker.C:
			// Routine ticking of the clock: Give the live
			// TLEs to the next index for this new time.

			// Also updated and emit Metrics.

			inCountDelta := inCount - inCount0
			inCount0 = inCount

			if n.Metrics != nil {
				now := time.Now().UTC()
				m := Metrics{
					T:          now,
					T1:         t1,
					In:         inCountDelta,
					Lag:        now.Sub(t),
					Goroutines: runtime.NumGoroutine(),
					Strings:    n.interns.IdCount(),
					Live:       len(n.live),
				}
				go func() {
					select {
					case <-ctx.Done():
					case n.Metrics <- m:
					}
				}()

			}

			// Remove and terminate earliest index.
			i := indexes[t0]
			delete(indexes, t0)
			i.Stop(ctx)
			t0 = t0.Add(n.Resolution)

			// Make the new index, and give it the live TLEs.
			i = n.NewIndex(t1)
			indexes[t1] = i
			go i.Run(ctx)
			n.logf(ctx, "Processing live sats (%d)", len(n.live))
			n.process(ctx, map[time.Time]*Index{
				t1: i,
			}, n.live)
			n.logf(ctx, "Processed live sats")

			// Increment our virtual clock.
			t1 = t1.Add(n.Resolution)

		case sats := <-n.In:
			// Process in-coming TLEs.
			inCount += uint64(len(sats))
			n.logf(ctx, "Processing %d new TLEs", len(sats))
			n.processNew(ctx, sats, indexes)
			n.logf(ctx, "Processed %d new TLEs", len(sats))
		}
	}

	n.logf(ctx, "Node.Run stopping")
}

// IndexInput represents (interned) data submitted to an Index.
type IndexInput struct {
	Id  index.Id
	Key index.Key
	Sat *PubTLE
}

// NewIndexInput builds an IndexInput.
//
// This method exists for callers to use when not using Run (say when
// doing batch processing).
//
// Also see the function NewIndexInput.
//
// This method indirectly obtains/releases a lock on the interned data.
func (n *Node) NewIndexInput(ctx context.Context, sat *PubTLE) *IndexInput {
	var ii *IndexInput
	f := func(is *Interns) error {
		ii = NewIndexInput(sat, is)
		return nil
	}
	n.interns.Exec(ctx, f)
	return ii
}

// NewIndexInput builds an IndexInput (assuming a lock if required).
func NewIndexInput(sat *PubTLE, is *Interns) *IndexInput {

	id, dup := is.Ids.Intern(sat)

	if dup {
		return nil
	}

	var (
		cat, _ = is.Keys.Intern(sat.TLE.CatNum)
		pub, _ = is.Keys.Intern(sat.Publisher)
	)
	return &IndexInput{
		Id: id,
		Key: index.Key{
			CatalogNum: index.CatalogNum(cat),
			Publisher:  index.Publisher(pub),
		},
		Sat: sat,
	}
}

// IndexOutput represents data produced by an Index.
type IndexOutput struct {
	// Time is the logical time (the time slice) for these events.
	time.Time

	// Novel is the set of novel Conjs.
	Novel []index.Conj

	// Canceled is the set of canceled Conjs.
	Canceled []index.Conj
}

// processNew interns new TLEs, stores them, and submits them to all
// existing indexes.
func (n *Node) processNew(ctx context.Context, sats []*PubTLE, is map[time.Time]*Index) error {
	iis := make(map[index.Key]*IndexInput, len(sats))

	// Make IndexInputs and store the new live TLEs while we're at
	// it.
	internWork := func(is *Interns) error {
		for _, sat := range sats {
			ii := NewIndexInput(sat, n.interns)
			if ii == nil { // Duplicate
				continue
			}
			iis[ii.Key] = ii
			n.live[ii.Key] = ii
		}
		return nil
	}
	if err := n.interns.Exec(ctx, internWork); err != nil {
		return err
	}

	// Submit input to all the tiven indexes, and wait for all of
	// that work to complete.
	if err := n.process(ctx, is, iis); err != nil {
		return err
	}

	// Update each Key's Id.
	internWork = func(is *Interns) error {
		for _, ii := range iis {
			n.interns.Update(ii.Key, ii.Id)
		}
		return nil
	}
	if err := n.interns.Exec(ctx, internWork); err != nil {
		return err
	}

	return nil
}

// IndexWork returns a function that performs the core index operation
// and all subsequent processing.
//
// The given function f is called with the batch of reports produced
// by indexing and subsequent processing.
//
// The returned function should be submitted to an index.
//
// Most of the work after the core index processing runs in a new
// goroutine.
func (n *Node) IndexWork(ctx context.Context, iis map[index.Key]*IndexInput, f func([]*Report)) func(*index.Index, time.Time) {

	return func(i *index.Index, t time.Time) {

		// This work will be done in the index's goroutine.

		ios := make([]*IndexOutput, 0, len(iis))

		for _, ii := range iis {

			// We might want to Prop in a batch in another goroutine.
			e, err := ii.Sat.TLE.Prop(t)
			if err != nil {
				n.logf(ctx, "sat.Prop %s", err)
				continue
			}

			pp := index.ProbPos{
				Pos: index.Pos{
					X: e.ECI.X,
					Y: e.ECI.Y,
					Z: e.ECI.Z,
				},
			}

			pps := []index.ProbPos{pp}

			cans, novs, _, err := i.Update(ii.Id, ii.Key, pps)

			io := &IndexOutput{
				Time:     t,
				Novel:    novs,
				Canceled: cans,
			}

			ios = append(ios, io)
		}

		// We don't need the index anymore.

		// We expect that most output will be empty, so let's
		// perform this initial consolidation without creating
		// a new goroutine.  ToDo: Reconsider.
		ios = n.Consolidate(ios)
		if 0 == len(ios) {
			f([]*Report{})
		} else {
			// We have some work to do, so let's do it in
			// a new goroutine.
			go func() {
				ps := n.GetIndexOutputTLEs(ctx, ios)
				rs := n.generateReports(ctx, ios, ps)
				n.logf(ctx, "Generated %d reports", len(rs))
				f(rs)
			}()
		}
	}
}

// process is constructs an IndexWork and submits it to all given
// indexes.
//
// This method waits for all of that work to complete before
// returning.  That behavior results in backpressure for the Node's
// input channel.
func (n *Node) process(ctx context.Context, indexes map[time.Time]*Index, iis map[index.Key]*IndexInput) error {
	n.logf(ctx, "Index processing: %d (%d indexes)", len(iis), len(indexes))

	var (
		done = ctx.Done()
		wg   = sync.WaitGroup{}

		f = func(rs []*Report) {
			if 0 < len(rs) {
				select {
				case <-done:
				case n.Out <- rs:
					n.logf(ctx, "Emitting %d reports", len(rs))
				}
			}
			wg.Done()
		}

		indexWork = n.IndexWork(ctx, iis, f)
	)

	// Give work to all of these indexes.
LOOP:
	for _, index := range indexes {
		wg.Add(1)
		select {
		case <-done:
			wg.Done() // Undo that Add.
			break LOOP
		case index.in <- indexWork:
		}
	}

	// Backpressure: Wait for all of that work to finish.
	wg.Wait()

	return nil
}

// Index is an index.Index associated with a time slice.
type Index struct {
	// t is the logical time slice for this index.
	t time.Time

	// I is of course the underlying index.Index.
	I *index.Index

	// in is the channel for submitting work to the Index's goroutine.
	in chan func(*index.Index, time.Time)

	// stop will halt the Run loop when stop is closed.
	stop chan bool
}

func (n *Node) NewIndex(t time.Time) *Index {
	return &Index{
		I:    index.NewIndexWithFinder(n.Finder, n.IndexDist),
		t:    t,
		in:   make(chan func(*index.Index, time.Time)),
		stop: make(chan bool),
	}
}

// Stop halts the Index's Run loop (if any).
func (i *Index) Stop(ctx context.Context) {
	close(i.stop)
}

// Run starts the work processing loop in the current goroutine.
//
// Use Stop to stop it.
func (i *Index) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-i.stop:
			return
		case f := <-i.in:
			f(i.I, i.t)
		}
	}
}

func (n *Node) GetIndexOutputTLEs(ctx context.Context, ios []*IndexOutput) map[index.Id]*PubTLE {
	var acc map[index.Id]*PubTLE
	f := func(is *Interns) error {
		acc = getIndexOutputTLEs(is, ios)
		return nil
	}
	n.interns.RExec(ctx, f)
	return acc
}

// getIndexOutputTLEs populates PubTLEs by looking up ids in the intern data.
//
// Assumes a lock.
func getIndexOutputTLEs(is *Interns, ios []*IndexOutput) map[index.Id]*PubTLE {
	acc := make(map[index.Id]*PubTLE)

	for _, io := range ios {

		for _, c := range io.Novel {
			for _, at := range c.Ats {
				if _, have := acc[at.Id]; have {
					continue
				}
				p, have := is.Ids.Find(at.Id)
				if !have {
					panic(fmt.Errorf("id %d missing", at.Id))
				}
				acc[at.Id] = p
			}
		}

		for _, c := range io.Canceled {
			for _, at := range c.Ats {
				if _, have := acc[at.Id]; have {
					continue
				}
				p, have := is.Ids.Find(at.Id)
				if !have {
					panic(fmt.Errorf("id %d missing", at.Id))
				}
				acc[at.Id] = p
			}
		}
	}

	return acc
}

// generateReports constructs Reports by calling ConjToReport.
func (n *Node) generateReports(ctx context.Context, ios []*IndexOutput, ps map[index.Id]*PubTLE) []*Report {
	var (
		rs         = make([]*Report, 0, len(ios))
		novs, cans int
	)

	for _, uo := range ios {
		for _, c := range uo.Novel {
			r, err := n.ConjToReport(uo.Time, &c, n.ScanDist, ps, false)
			if err != nil {
				if sgp4.HasDecayed(err) {
					continue
				}
				n.logf(ctx, "ConjToReport (nov): %s", err)
				continue
			}
			if r == nil {
				continue
			}
			rs = append(rs, r)
			novs++
		}
		for _, c := range uo.Canceled {
			r, err := n.ConjToReport(uo.Time, &c, n.ScanDist, ps, true)
			if err != nil {
				if sgp4.HasDecayed(err) {
					continue
				}
				n.logf(ctx, "ConjToReport (can): %s", err)
			}
			if r == nil {
				continue
			}
			rs = append(rs, r)
			cans++
		}
	}
	n.logf(ctx, "generateReports: novs: %d, cans: %d, total: %d", novs, cans, len(rs))

	return rs
}

// State represents the state of an object at some (implicit) time.
type State struct {
	// Name is the object's Name().
	Name string

	// Obj is the TLE.
	//
	// We say "Obj" to facilitate generalization from TLEs
	// specifically.
	Obj *tle.SGP4TLE

	// Age is the age of the source (TLE) in seconds.
	Age int64 `json:",omitempty"` // Seconds

	// Type is a crude classification of the object.
	//
	// See TLE.GetType().
	Type string

	// ECI is the position vector.
	ECI prop.Vect

	// Vel is relative velocity (m/s).
	Vel prop.Vect

	// LLA is latittude (deg), longitude (deg), and altitude (km)
	LLA LatLonAlt

	// ToDo: Prob (again)
}

// Report is a complete conjunction report: what we are here for.
type Report struct {
	// Id is a logical identifier for this report.
	Id string

	// Sig is the signature for this report without consideration of Canceled and Generated.
	//
	// This data can be used to find a previous report that is canceled by this report.
	Sig string

	// Generated is the real time that this report was Generated.
	Generated time.Time

	// At is the logical time for the event reported here.
	At time.Time

	// Canceled indicates that this report is a cancellation of a
	// previous report (which will have the same Sig).
	Canceled bool `json:",omitempty"`

	// Dist is the estimated Cartesian distance (km) between the two Objs.
	Dist float32

	// Speed is the estimated Relative speed (m/s) between the two Objs.
	Speed float32

	// Objs is an array of the State of the two objects in this event.
	Objs []State
}

// ConjToReport builds the final Reports.
//
// This method doesn't require any external data (other than that
// passed as arguments).
func (n *Node) ConjToReport(t time.Time, c *index.Conj, dist float32, ps map[index.Id]*PubTLE, canceled bool) (*Report, error) {

	// Obtain the (populated) object instances based on their ids.

	o0, have := ps[c.Ats[0].Id]
	if !have {
		return nil, Warningf("Node ids index doesn't have %v", c.Ats[0].Id)
	}

	o1, have := ps[c.Ats[1].Id]
	if !have {
		return nil, Warningf("Node ids index doesn't have %v", c.Ats[1].Id)
	}

	// Possibly scan for a closer approach +/- one tick.  ScanPair
	// won't actually do any real scanning if n.Scan is false.
	d, es, then, err := ScanPair(n.Scan, t, time.Second, 100, o0.TLE, o1.TLE, 0)
	if err != nil {
		return nil, err
	}

	if dist < d {
		return nil, nil
	}

	v := es[0].V.Dist(es[1].V)

	if 0 < n.SlowSample && v < n.SlowSampleThreshold {
		var (
			maxRate = float32(n.SlowSample)
			minRate = float32(1)
			scaled  = 1 - v/n.SlowSampleThreshold
			m       = int(minRate + (maxRate-minRate)*scaled)

			secs = t.Unix()
		)
		if secs%int64(m) != 0 {
			return nil, nil
		}
	}

	// p := c.Ats[0].Prob * c.Ats[1].Prob

	l0, err := ECIToLLA(t, es[0].ECI)
	if err != nil {
		return nil, err
	}
	s0 := State{
		Name: o0.Name(),
		Obj:  o0.TLE,
		Age:  int64(o0.TLE.ApproxAge(t).Seconds()),
		Type: o0.TLE.GetType(),
		// Prob: c.Ats[0].ProbPos.Prob,
		ECI: es[0].ECI,
		Vel: es[0].V,
		LLA: *l0,
	}

	l1, err := ECIToLLA(t, es[1].ECI)
	if err != nil {
		return nil, err
	}
	s1 := State{
		Name: o1.Name(),
		Obj:  o1.TLE,
		Age:  int64(o1.TLE.ApproxAge(t).Seconds()),
		Type: o1.TLE.GetType(),
		// Prob: c.Ats[1].ProbPos.Prob,
		ECI: es[1].ECI,
		Vel: es[1].V,
		LLA: *l1,
	}

	r := &Report{
		At:    then,
		Dist:  d,
		Speed: v,
		Objs:  []State{s0, s1},
	}

	{
		// Sig does not include Canceled or Generated.
		js, err := json.Marshal(r)
		if err != nil {
			return nil, err
		}
		r.Sig = misc.SHA(js)
		r.Canceled = canceled
		r.Generated = time.Now().UTC()

		// Id includes Canceled and Generated.
		js, err = json.Marshal(r)
		if err != nil {
			return nil, err
		}
		r.Id = misc.SHA(js)
	}

	return r, nil
}
