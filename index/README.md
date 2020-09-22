# A geospatial index

This code provides an index for object proximity.  The index structure
is based on [S2 Geometry](https://s2geometry.io/), and this code uses
the [`github.com/golang/geo`](https://github.com/golang/geo)
implementation.

The primary method is `Update`, which takes a publisher's opinion
about the set of possible positions, each with an associated
probability, for an object key in the publisher's namespace.  Input
also includes the identifier for data the publisher used to generate
this input.  `Update` returns the set of obsolete conjunction reports
(if any) and a set of new conjunction reports (if any).  This method
also returns the id (if any) the publisher previously used for the
object key.

This implementation is not safe for concurrent use.
