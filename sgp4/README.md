# SGP4 C transpiled to Go

An
[SGP4](http://celestrak.com/publications/AIAA/2006-6753/AIAA-2006-6753-Rev2.pdf)
implementation based on
[`github.com/aholinch/sgp4/tree/master/src/c`](https://github.com/aholinch/sgp4/tree/master/src/c),
which ultimated originated with David Vallado.

```
    This file contains the sgp4 procedures for analytical propagation
    of a satellite. the code was originally released in the 1980 and
    1986 spacetrack papers. a detailed discussion of the theory and
    history may be found in the 2006 aiaa paper by vallado, crawford,
    hujsak, and kelso.

                           companion code for
              fundamentals of astrodynamics and applications
                                   2013
                             by david vallado
     (w) 719-573-2600, email dvallado@agi.com, davallado@gmail.com
```

This implementation is a hand-edited transpilation of [C
sources](https://github.com/aholinch/sgp4/tree/master/src/c) to Go by
[`c2go`](https://github.com/elliotchance/c2go) (version v0.25.9
Dubnium 2018-12-30), and the emitted code was edited by hand.  The
original C implementation test suite was included in this process.

The substantive edits (of Go sources emitted by the transpiler) were
the use of 64-bit integers to address at least one 32-bit overflow and
using floating point constants instead of naked integer constants when
the entire expression consisted of the latter with some division.
Example: "var x2o3 float64 = 2 / 3" was edited to be "var x2o3 float64
= 2.0 / 3.0".

With those changes, the original (transpiled and hand-edited) tests
almost all pass.  The one exception is for objectNum 20413 at
mins=1844335, where 1e-07 < rdist < 1e-06.  The tests have been edited
to tolerate rdist < 1e-07 rather than demand rdist < 1e-06.

