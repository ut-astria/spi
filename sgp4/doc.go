// Package sgp4 provides a method for analytical propagation of a
// satellite.
//
// The code here was hand-edited from a transpilation of C sources.
//
// The transpiler was c2go from https://github.com/elliotchance/c2go
// (version v0.25.9 Dubnium 2018-12-30).
//
// The original C code was from
// https://github.com/aholinch/sgp4/tree/master/src/c, which itself
// originated from code authored by David Vallado:
//
//    This file contains the sgp4 procedures for analytical propagation
//    of a satellite. the code was originally released in the 1980 and
//    1986 spacetrack papers. a detailed discussion of the theory and
//    history may be found in the 2006 aiaa paper by vallado, crawford,
//    hujsak, and kelso.
//
//                           companion code for
//              fundamentals of astrodynamics and applications
//                                   2013
//                             by david vallado
//     (w) 719-573-2600, email dvallado@agi.com, davallado@gmail.com
//
// The substantive edits (of Go sources emitted by the transpiler)
// were the use of 64-bit integers to address at least one 32-bit
// overflow and using floating point constants instead of naked
// integer constants when the entire expression consisted of the
// latter with some division.  Example: "var x2o3 float64 = 2 / 3" was
// edited to be "var x2o3 float64 = 2.0 / 3.0".
//
// With those changes, the original included tests almost all pass.
// The one exception is for objectNum 20413 at mins=1844335, where
// 1e-07 < rdist < 1e-06.  The tests have been edited to tolerate
// rdist <= 1e-06.
//
//
// Error codes:
//
//   1 - mean elements, ecc >= 1.0 or ecc < -0.001 or a < 0.95 er
//   2 - mean motion less than 0.0
//   3 - pert elements, ecc < 0.0  or  ecc > 1.0
//   4 - semi-latus rectum < 0.0
//   5 - epoch elements are sub-orbital
//   6 - satellite has decayed
//
package sgp4
