package tle

// Src is the source of a TLE or thing that can be propagated.
//
// ToDo: Move out of this package and into 'prop'?
type Src struct {
	// Publisher is the organization-based origin for this data.
	Publisher string

	// Obj is the name of the object (in an unspecified
	// namespace).
	Obj string
}
