// Package index implements a stream-able geospatial index on
// positions in Cartesian 3-space.
//
// Object-positions go in, and proximity (or "conjunction") reports
// come out.
//
// Indexing is based on cells that cover S2.
//
// The primary function is NewIndex, and the primary method is Update.
package index
