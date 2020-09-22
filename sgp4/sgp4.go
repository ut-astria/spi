/*
  Originally transpiled from

     https://github.com/aholinch/sgp4/tree/master/src/c

  by c2go

     https://github.com/elliotchance/c2go

  (version v0.25.9 Dubnium 2018-12-30).

  Then edited by hand.

  Original C code credit:

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

*/

/*     ----------------------------------------------------------------
*
*                               sgp4unit.cpp
*
*    this file contains the sgp4 procedures for analytical propagation
*    of a satellite. the code was originally released in the 1980 and 1986
*    spacetrack papers. a detailed discussion of the theory and history
*    may be found in the 2006 aiaa paper by vallado, crawford, hujsak,
*    and kelso.
*
*                            companion code for
*               fundamentals of astrodynamics and applications
*                                    2013
*                              by david vallado
*
*     (w) 719-573-2600, email dvallado@agi.com, davallado@gmail.com
*
*    current :
*               7 dec 15  david vallado
*                           fix jd, jdfrac
*    changes :
*               3 nov 14  david vallado
*                           update to msvs2013 c++
*              30 aug 10  david vallado
*                           delete unused variables in initl
*                           replace pow integer 2, 3 with multiplies for speed
*               3 nov 08  david vallado
*                           put returns in for error codes
*              29 sep 08  david vallado
*                           fix atime for faster operation in dspace
*                           add operationmode for afspc (a) or improved (i)
*                           performance mode
*              16 jun 08  david vallado
*                           update small eccentricity check
*              16 nov 07  david vallado
*                           misc fixes for better compliance
*              20 apr 07  david vallado
*                           misc fixes for constants
*              11 aug 06  david vallado
*                           chg lyddane choice back to strn3, constants, misc doc
*              15 dec 05  david vallado
*                           misc fixes
*              26 jul 05  david vallado
*                           fixes for paper
*                           note that each fix is preceded by a
*                           comment with "sgp4fix" and an explanation of
*                           what was changed
*              10 aug 04  david vallado
*                           2nd printing baseline working
*              14 may 01  david vallado
*                           2nd edition baseline
*                     80  norad
*                           original baseline
*       ----------------------------------------------------------------      */ //

package sgp4

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"unsafe"
	// "github.com/elliotchance/c2go/noarch"
)

// cs2s takes a C string and returns a Go string.
//
// Almost certainly something better exists.
func cs2s(bs []byte) string {
	for i, b := range bs {
		if b == 0 {
			return string(bs[0:i])
		}
	}
	return string(bs)
}

func sscanf(s string, format string, args ...interface{}) int {
	n, _ := fmt.Sscanf(s, format, args)
	return n
}

type bool_ int64

type ElsetRec struct {
	whichconst     int64
	satnum         int64
	epochyr        int64
	epochtynumrev  int64
	error          int64
	operationmode  byte
	init_          byte
	method         byte
	a              float64
	altp           float64
	alta           float64
	epochdays      float64
	jdsatepoch     float64
	jdsatepochF    float64
	nddot          float64
	ndot           float64
	bstar          float64
	rcse           float64
	inclo          float64
	nodeo          float64
	ecco           float64
	argpo          float64
	mo             float64
	no_kozai       float64
	classification byte
	intldesg       [12]byte
	ephtype        int64
	elnum          int64
	revnum         int64
	no_unkozai     float64
	am             float64
	em             float64
	im             float64
	Om             float64
	om             float64
	mm             float64
	nm             float64
	t              float64
	tumin          float64
	mu             float64
	radiusearthkm  float64
	xke            float64
	j2             float64
	j3             float64
	j4             float64
	j3oj2          float64
	dia_mm         int64
	period_sec     float64
	active         byte
	not_orbital    byte
	rcs_m2         float64
	ep             float64
	inclp          float64
	nodep          float64
	argpp          float64
	mp             float64
	isimp          int64
	aycof          float64
	con41          float64
	cc1            float64
	cc4            float64
	cc5            float64
	d2             float64
	d3             float64
	d4             float64
	delmo          float64
	eta            float64
	argpdot        float64
	omgcof         float64
	sinmao         float64
	t2cof          float64
	t3cof          float64
	t4cof          float64
	t5cof          float64
	x1mth2         float64
	x7thm1         float64
	mdot           float64
	nodedot        float64
	xlcof          float64
	xmcof          float64
	nodecf         float64
	irez           int64
	d2201          float64
	d2211          float64
	d3210          float64
	d3222          float64
	d4410          float64
	d4422          float64
	d5220          float64
	d5232          float64
	d5421          float64
	d5433          float64
	dedt           float64
	del1           float64
	del2           float64
	del3           float64
	didt           float64
	dmdt           float64
	dnodt          float64
	domdt          float64
	e3             float64
	ee2            float64
	peo            float64
	pgho           float64
	pho            float64
	pinco          float64
	plo            float64
	se2            float64
	se3            float64
	sgh2           float64
	sgh3           float64
	sgh4           float64
	sh2            float64
	sh3            float64
	si2            float64
	si3            float64
	sl2            float64
	sl3            float64
	sl4            float64
	gsto           float64
	xfact          float64
	xgh2           float64
	xgh3           float64
	xgh4           float64
	xh2            float64
	xh3            float64
	xi2            float64
	xi3            float64
	xl2            float64
	xl3            float64
	xl4            float64
	xlamo          float64
	zmol           float64
	zmos           float64
	atime          float64
	xli            float64
	xni            float64
	snodm          float64
	cnodm          float64
	sinim          float64
	cosim          float64
	sinomm         float64
	cosomm         float64
	day            float64
	emsq           float64
	gam            float64
	rtemsq         float64
	s1             float64
	s2             float64
	s3             float64
	s4             float64
	s5             float64
	s6             float64
	s7             float64
	ss1            float64
	ss2            float64
	ss3            float64
	ss4            float64
	ss5            float64
	ss6            float64
	ss7            float64
	sz1            float64
	sz2            float64
	sz3            float64
	sz11           float64
	sz12           float64
	sz13           float64
	sz21           float64
	sz22           float64
	sz23           float64
	sz31           float64
	sz32           float64
	sz33           float64
	z1             float64
	z2             float64
	z3             float64
	z11            float64
	z12            float64
	z13            float64
	z21            float64
	z22            float64
	z23            float64
	z31            float64
	z32            float64
	z33            float64
	argpm          float64
	inclm          float64
	nodem          float64
	dndt           float64
	eccsq          float64
	ainv           float64
	ao             float64
	con42          float64
	cosio          float64
	cosio2         float64
	omeosq         float64
	posq           float64
	rp             float64
	rteosq         float64
	sinio          float64
}
type TLE struct {
	sync.Mutex

	Rec       ElsetRec
	line1     [70]byte
	line2     [70]byte
	intlid    [12]byte
	objectNum int64
	epoch     int64
	ndot      float64
	nddot     float64
	bstar     float64
	elnum     int64
	incDeg    float64
	raanDeg   float64
	ecc       float64
	argpDeg   float64
	maDeg     float64
	n         float64
	revnum    int64
	sgp4Error int64
}

// parseLines - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:16
// parse the double
// parse the double with implied decimal
// copy the lines
//          1         2         3         4         5         6
//0123456789012345678901234567890123456789012345678901234567890123456789
//line1="1 00005U 58002B   00179.78495062  .00000023  00000-0  28098-4 0  4753";
//line2="2 00005  34.2682 348.7242 1859667 331.7664  19.3264 10.82419157413667";
// intlid
//
func parseLines(tle *TLE, line1 *byte, line2 *byte) {
	(*tle).Rec.whichconst = int64(2)
	Strncpy64(&(*tle).line1[0], line1, int64(uint64(int64(69))))
	Strncpy64(&(*tle).line2[0], line2, int64(uint64(int64(69))))
	*((*byte)(func() unsafe.Pointer {
		tempVar := &(*tle).line1[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(69))*unsafe.Sizeof(*tempVar))
	}())) = byte(int64(0))
	*((*byte)(func() unsafe.Pointer {
		tempVar := &(*tle).line2[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(69))*unsafe.Sizeof(*tempVar))
	}())) = byte(int64(0))
	Strncpy64(&(*tle).intlid[0], &*((*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(line1)) + (uintptr)(int64(9))*unsafe.Sizeof(*line1)))), int64(uint64(int64(8))))
	(*tle).Rec.classification = *((*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(line1)) + (uintptr)(int64(7))*unsafe.Sizeof(*line1))))
	(*tle).objectNum = int64(gd(line1, int64(2), int64(7)))
	(*tle).ndot = gdi(line1, int64(35), int64(44))
	if int64(*((*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(line1)) + (uintptr)(int64(33))*unsafe.Sizeof(*line1))))) == int64('-') {
		(*tle).ndot *= -1
	}
	(*tle).nddot = gdi(line1, int64(45), int64(50))
	if int64(*((*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(line1)) + (uintptr)(int64(44))*unsafe.Sizeof(*line1))))) == int64('-') {
		(*tle).nddot *= -1
	}
	(*tle).nddot *= math.Pow(10, gd(line1, int64(50), int64(52)))
	(*tle).bstar = gdi(line1, int64(54), int64(59))
	if int64(*((*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(line1)) + (uintptr)(int64(53))*unsafe.Sizeof(*line1))))) == int64('-') {
		(*tle).bstar *= -1
	}
	(*tle).bstar *= math.Pow(10, gd(line1, int64(59), int64(61)))
	(*tle).elnum = int64(gd(line1, int64(64), int64(68)))
	(*tle).incDeg = gd(line2, int64(8), int64(16))
	(*tle).raanDeg = gd(line2, int64(17), int64(25))
	(*tle).ecc = gdi(line2, int64(26), int64(33))
	(*tle).argpDeg = gd(line2, int64(34), int64(42))
	(*tle).maDeg = gd(line2, int64(43), int64(51))
	(*tle).n = gd(line2, int64(52), int64(63))
	(*tle).revnum = int64(gd(line2, int64(63), int64(68)))
	(*tle).sgp4Error = int64(0)
	(*tle).epoch = parseEpoch(&(*tle).Rec, &*((*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(line1)) + (uintptr)(int64(18))*unsafe.Sizeof(*line1)))))
	setValsToRec(tle, &(*tle).Rec)
}

// isLeap - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:66
func isLeap(year int64) bool_ {
	if year%int64(4) != int64(0) {
		return bool_((int64(0)))
	}
	if year%int64(100) == int64(0) {
		if year%int64(400) == int64(0) {
			return bool_((int64(1)))
		} else {
			return bool_((int64(0)))
		}
	}
	return bool_((int64(1)))
}

// parseEpoch - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:88
// convert doy to mon, day
//
func parseEpoch(rec *ElsetRec, str *byte) int64 {
	var tmp []byte = make([]byte, 16, 16)
	Strncpy64(&tmp[0], str, int64(uint64(int64(14))))
	*((*byte)(func() unsafe.Pointer {
		tempVar := &tmp[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(15))*unsafe.Sizeof(*tempVar))
	}())) = byte(int64(0))
	var tmp2 []byte = make([]byte, 16, 16)
	Strncpy64(&tmp2[0], &tmp[0], int64(uint64(int64(2))))
	*((*byte)(func() unsafe.Pointer {
		tempVar := &tmp2[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(2))*unsafe.Sizeof(*tempVar))
	}())) = byte(int64(0))

	// var year int64 = noarch.Atoi64(&tmp2[0])
	// noarch
	year, _ := strconv.ParseInt(string(tmp2[0:2]), 10, 64)

	(*rec).epochyr = year
	if year > int64(56) {
		year += int64(1900)
	} else {
		year += int64(2000)
	}
	Strncpy64(&tmp2[0], &*((*byte)(func() unsafe.Pointer {
		tempVar := &tmp[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(2))*unsafe.Sizeof(*tempVar))
	}())), int64(uint64(int64(3))))
	*((*byte)(func() unsafe.Pointer {
		tempVar := &tmp2[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(3))*unsafe.Sizeof(*tempVar))
	}())) = byte(int64(0))

	// var doy int64 = noarch.Atoi64(&tmp2[0])
	// noarch
	doy, _ := strconv.ParseInt(string(tmp2[0:3]), 10, 64)

	*&tmp2[0] = '0'
	Strncpy64(&*((*byte)(func() unsafe.Pointer {
		tempVar := &tmp2[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(1))*unsafe.Sizeof(*tempVar))
	}())), &*((*byte)(func() unsafe.Pointer {
		tempVar := &tmp[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(5))*unsafe.Sizeof(*tempVar))
	}())), int64(uint64(int64(9))))
	*((*byte)(func() unsafe.Pointer {
		tempVar := &tmp2[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(11))*unsafe.Sizeof(*tempVar))
	}())) = byte(int64(0))

	// var dfrac float64 = noarch.Strtod(&tmp2[0], nil)
	// noarch
	dfrac, _ := strconv.ParseFloat(cs2s(tmp2), 64)

	// var odfrac float64 = dfrac
	(*rec).epochdays = float64(doy)
	(*rec).epochdays += dfrac
	dfrac *= 24
	var hr int64 = int64(dfrac)
	dfrac = 60 * (dfrac - float64(hr))
	var mn int64 = int64(dfrac)
	dfrac = 60 * (dfrac - float64(mn))
	var sc int64 = int64(dfrac)
	dfrac = 1000 * (dfrac - float64(sc))
	// var milli int64 = int64(dfrac)
	var sec float64 = float64(sc) + dfrac/1000

	var mon int64 = int64(0)
	var day int64 = int64(0)
	var days []int64 = []int64{int64(31), int64(28), int64(31), int64(30), int64(31), int64(30), int64(31), int64(31), int64(30), int64(31), int64(30), int64(31)}
	if isLeap(year) == 1 {
		*((*int64)(func() unsafe.Pointer {
			tempVar := &days[0]
			return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(1))*unsafe.Sizeof(*tempVar))
		}())) = int64(29)
	}
	var ind int64 = int64(0)
	for ind < int64(12) && doy > *((*int64)(func() unsafe.Pointer {
		tempVar := &days[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(ind)*unsafe.Sizeof(*tempVar))
	}())) {
		doy -= *((*int64)(func() unsafe.Pointer {
			tempVar := &days[0]
			return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(ind)*unsafe.Sizeof(*tempVar))
		}()))
		ind += 1
	}
	mon = ind + int64(1)
	day = doy
	jday(year, mon, day, hr, mn, sec, &(*rec).jdsatepoch, &(*rec).jdsatepochF)

	var diff float64 = (*rec).jdsatepoch - 2.4405875e+06
	var diff2 float64 = 8.64e+07 * (*rec).jdsatepochF
	diff *= 8.64e+07
	var epoch = int64(diff2)
	epoch += int64(diff)

	return epoch
}

// getRVForDate - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:164
func getRVForDate(tle *TLE, millisSince1970 int64, r *float64, v *float64) {
	var diff float64 = float64(millisSince1970) - float64((*tle).epoch)
	diff /= 60000
	getRV(tle, diff, r, v)
}

// getRV - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:171
func getRV(tle *TLE, minutesAfterEpoch float64, r *float64, v *float64) {
	(*tle).Rec.error = int64(0)
	sgp4(&(*tle).Rec, minutesAfterEpoch, r, v)
	(*tle).sgp4Error = (*tle).Rec.error
}

// gd - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:178
func gd(str *byte, ind1 int64, ind2 int64) float64 {
	var num float64 = float64(int64(0))
	var tmp []byte = make([]byte, 50, 50)
	var cnt int64 = ind2 - ind1
	Strncpy64(&tmp[0], &*((*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(str)) + (uintptr)(ind1)*unsafe.Sizeof(*str)))), int64(uint64(cnt)))
	*((*byte)(func() unsafe.Pointer {
		tempVar := &tmp[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(cnt)*unsafe.Sizeof(*tempVar))
	}())) = byte(int64(0))

	// num = noarch.Strtod(&tmp[0], nil)
	// noarch
	num, _ = strconv.ParseFloat(strings.TrimSpace(cs2s(tmp)), 64)

	return num
}

// gdi - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:190
// parse with an implied decimal place
//
func gdi(str *byte, ind1 int64, ind2 int64) float64 {
	var num float64 = float64(int64(0))
	var tmp []byte = make([]byte, 52, 52)
	*&tmp[0] = '0'
	*((*byte)(func() unsafe.Pointer {
		tempVar := &tmp[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(1))*unsafe.Sizeof(*tempVar))
	}())) = '.'
	var cnt int64 = ind2 - ind1
	Strncpy64(&*((*byte)(func() unsafe.Pointer {
		tempVar := &tmp[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(2))*unsafe.Sizeof(*tempVar))
	}())), &*((*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(str)) + (uintptr)(ind1)*unsafe.Sizeof(*str)))), int64(uint64(cnt)))
	*((*byte)(func() unsafe.Pointer {
		tempVar := &tmp[0]
		return unsafe.Pointer(uintptr(unsafe.Pointer(tempVar)) + (uintptr)(int64(2)+cnt)*unsafe.Sizeof(*tempVar))
	}())) = byte(int64(0))

	// num = noarch.Strtod(&tmp[0], nil)
	// noarch
	num, _ = strconv.ParseFloat(cs2s(tmp), 64)

	return num
}

// setValsToRec - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:203
// 229.1831180523293
//
func setValsToRec(tle *TLE, rec *ElsetRec) {
	var xpdotp float64 = 1440.0 / (2.0 * 3.141592653589793)
	(*rec).elnum = (*tle).elnum
	(*rec).revnum = (*tle).revnum
	(*rec).satnum = (*tle).objectNum
	(*rec).bstar = (*tle).bstar
	(*rec).inclo = (*tle).incDeg * (3.141592653589793 / 180.0)
	(*rec).nodeo = (*tle).raanDeg * (3.141592653589793 / 180.0)
	(*rec).argpo = (*tle).argpDeg * (3.141592653589793 / 180.0)
	(*rec).mo = (*tle).maDeg * (3.141592653589793 / 180.0)
	(*rec).ecco = (*tle).ecc
	(*rec).no_kozai = (*tle).n / xpdotp
	(*rec).ndot = (*tle).ndot / (xpdotp * 1440.0)
	(*rec).nddot = (*tle).nddot / (xpdotp * 1440.0 * 1440.0)
	sgp4init('a', rec)
}

// dpper - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:348
/*     ----------------------------------------------------------------
*
*                               sgp4unit.cpp
*
*    this file contains the sgp4 procedures for analytical propagation
*    of a satellite. the code was originally released in the 1980 and 1986
*    spacetrack papers. a detailed discussion of the theory and history
*    may be found in the 2006 aiaa paper by vallado, crawford, hujsak,
*    and kelso.
*
*                            companion code for
*               fundamentals of astrodynamics and applications
*                                    2013
*                              by david vallado
*
*     (w) 719-573-2600, email dvallado@agi.com, davallado@gmail.com
*
*    current :
*               7 dec 15  david vallado
*                           fix jd, jdfrac
*    changes :
*               3 nov 14  david vallado
*                           update to msvs2013 c++
*              30 aug 10  david vallado
*                           delete unused variables in initl
*                           replace pow integer 2, 3 with multiplies for speed
*               3 nov 08  david vallado
*                           put returns in for error codes
*              29 sep 08  david vallado
*                           fix atime for faster operation in dspace
*                           add operationmode for afspc (a) or improved (i)
*                           performance mode
*              16 jun 08  david vallado
*                           update small eccentricity check
*              16 nov 07  david vallado
*                           misc fixes for better compliance
*              20 apr 07  david vallado
*                           misc fixes for constants
*              11 aug 06  david vallado
*                           chg lyddane choice back to strn3, constants, misc doc
*              15 dec 05  david vallado
*                           misc fixes
*              26 jul 05  david vallado
*                           fixes for paper
*                           note that each fix is preceded by a
*                           comment with "sgp4fix" and an explanation of
*                           what was changed
*              10 aug 04  david vallado
*                           2nd printing baseline working
*              14 may 01  david vallado
*                           2nd edition baseline
*                     80  norad
*                           original baseline
*       ----------------------------------------------------------------      */ //
/* -----------------------------------------------------------------------------
 *
 *                           procedure dpper
 *
 *  this procedure provides deep space long period periodic contributions
 *    to the mean elements.  by design, these periodics are zero at epoch.
 *    this used to be dscom which included initialization, but it's really a
 *    recurring function.
 *
 *  author        : david vallado                  719-573-2600   28 jun 2005
 *
 *  inputs        :
 *    e3          -
 *    ee2         -
 *    peo         -
 *    pgho        -
 *    pho         -
 *    pinco       -
 *    plo         -
 *    se2 , se3 , sgh2, sgh3, sgh4, sh2, sh3, si2, si3, sl2, sl3, sl4 -
 *    t           -
 *    xh2, xh3, xi2, xi3, xl2, xl3, xl4 -
 *    zmol        -
 *    zmos        -
 *    ep          - eccentricity                           0.0 - 1.0
 *    inclo       - inclination - needed for lyddane modification
 *    nodep       - right ascension of ascending node
 *    argpp       - argument of perigee
 *    mp          - mean anomaly
 *
 *  outputs       :
 *    ep          - eccentricity                           0.0 - 1.0
 *    inclp       - inclination
 *    nodep        - right ascension of ascending node
 *    argpp       - argument of perigee
 *    mp          - mean anomaly
 *
 *  locals        :
 *    alfdp       -
 *    betdp       -
 *    cosip  , sinip  , cosop  , sinop  ,
 *    dalf        -
 *    dbet        -
 *    dls         -
 *    f2, f3      -
 *    pe          -
 *    pgh         -
 *    ph          -
 *    pinc        -
 *    pl          -
 *    sel   , ses   , sghl  , sghs  , shl   , shs   , sil   , sinzf , sis   ,
 *    sll   , sls
 *    xls         -
 *    xnoh        -
 *    zf          -
 *    zm          -
 *
 *  coupling      :
 *    none.
 *
 *  references    :
 *    hoots, roehrich, norad spacetrack report #3 1980
 *    hoots, norad spacetrack report #6 1986
 *    hoots, schumacher and glover 2004
 *    vallado, crawford, hujsak, kelso  2006
 ----------------------------------------------------------------------------*/ //
/* --------------------- local variables ------------------------ */ //
/* ---------------------- constants ----------------------------- */ //
/* --------------- calculate time varying periodics ----------- */ //
// be sure that the initial call has time set to zero
/* ----------------- apply periodics directly ------------ */ //
//  sgp4fix for lyddane choice
//  strn3 used original inclination - this is technically feasible
//  gsfc used perturbed inclination - also technically feasible
//  probably best to readjust the 0.2 limit value and limit discontinuity
//  0.2 rad = 11.45916 deg
//  use next line for original strn3 approach and original inclination
//  if (inclo >= 0.2)
//  use next line for gsfc version and perturbed inclination
/* ---- apply periodics with lyddane modification ---- */ //
//  sgp4fix for afspc written intrinsic functions
// nodep used without a trigonometric function ahead
//  sgp4fix for afspc written intrinsic functions
// nodep used without a trigonometric function ahead
// if init == 'n'
// dpper
//
func dpper(e3 float64, ee2 float64, peo float64, pgho float64, pho float64, pinco float64, plo float64, se2 float64, se3 float64, sgh2 float64, sgh3 float64, sgh4 float64, sh2 float64, sh3 float64, si2 float64, si3 float64, sl2 float64, sl3 float64, sl4 float64, t float64, xgh2 float64, xgh3 float64, xgh4 float64, xh2 float64, xh3 float64, xi2 float64, xi3 float64, xl2 float64, xl3 float64, xl4 float64, zmol float64, zmos float64, init_ byte, rec *ElsetRec, opsmode byte) {
	var alfdp float64
	var betdp float64
	var cosip float64
	var cosop float64
	var dalf float64
	var dbet float64
	var dls float64
	var f2 float64
	var f3 float64
	var pe float64
	var pgh float64
	var ph float64
	var pinc float64
	var pl float64
	var sel float64
	var ses float64
	var sghl float64
	var sghs float64
	var shll float64
	var shs float64
	var sil float64
	var sinip float64
	var sinop float64
	var sinzf float64
	var sis float64
	var sll float64
	var sls float64
	var xls float64
	var xnoh float64
	var zf float64
	var zm float64
	var zel float64
	var zes float64
	var znl float64
	var zns float64
	zns = 1.19459e-05
	zes = 0.01675
	znl = 0.00015835218
	zel = 0.0549
	zm = zmos + zns*t
	if int64(init_) == int64('y') {
		zm = zmos
	}
	zf = zm + 2*zes*math.Sin(zm)
	sinzf = math.Sin(zf)
	f2 = 0.5*sinzf*sinzf - 0.25
	f3 = -0.5 * sinzf * math.Cos(zf)
	ses = se2*f2 + se3*f3
	sis = si2*f2 + si3*f3
	sls = sl2*f2 + sl3*f3 + sl4*sinzf
	sghs = sgh2*f2 + sgh3*f3 + sgh4*sinzf
	shs = sh2*f2 + sh3*f3
	zm = zmol + znl*t
	if int64(init_) == int64('y') {
		zm = zmol
	}
	zf = zm + 2*zel*math.Sin(zm)
	sinzf = math.Sin(zf)
	f2 = 0.5*sinzf*sinzf - 0.25
	f3 = -0.5 * sinzf * math.Cos(zf)
	sel = ee2*f2 + e3*f3
	sil = xi2*f2 + xi3*f3
	sll = xl2*f2 + xl3*f3 + xl4*sinzf
	sghl = xgh2*f2 + xgh3*f3 + xgh4*sinzf
	shll = xh2*f2 + xh3*f3
	pe = ses + sel
	pinc = sis + sil
	pl = sls + sll
	pgh = sghs + sghl
	ph = shs + shll
	if int64(init_) == int64('n') {
		pe = pe - peo
		pinc = pinc - pinco
		pl = pl - plo
		pgh = pgh - pgho
		ph = ph - pho
		(*rec).inclp = (*rec).inclp + pinc
		(*rec).ep = (*rec).ep + pe
		sinip = math.Sin((*rec).inclp)
		cosip = math.Cos((*rec).inclp)
		if (*rec).inclp >= 0.2 {
			ph = ph / sinip
			pgh = pgh - cosip*ph
			(*rec).argpp = (*rec).argpp + pgh
			(*rec).nodep = (*rec).nodep + ph
			(*rec).mp = (*rec).mp + pl
		} else {
			sinop = math.Sin((*rec).nodep)
			cosop = math.Cos((*rec).nodep)
			alfdp = sinip * sinop
			betdp = sinip * cosop
			dalf = ph*cosop + pinc*cosip*sinop
			dbet = -ph*sinop + pinc*cosip*cosop
			alfdp = alfdp + dalf
			betdp = betdp + dbet
			(*rec).nodep = math.Mod((*rec).nodep, (2 * 3.141592653589793))
			if (*rec).nodep < 0 && int64(opsmode) == int64('a') {
				(*rec).nodep = (*rec).nodep + 2*3.141592653589793
			}
			xls = (*rec).mp + (*rec).argpp + cosip*(*rec).nodep
			dls = pl + pgh - pinc*(*rec).nodep*sinip
			xls = xls + dls
			xls = math.Mod(xls, (2 * 3.141592653589793))
			xnoh = (*rec).nodep
			(*rec).nodep = math.Atan2(alfdp, betdp)
			if (*rec).nodep < 0 && int64(opsmode) == int64('a') {
				(*rec).nodep = (*rec).nodep + 2*3.141592653589793
			}
			if math.Abs(xnoh-(*rec).nodep) > 3.141592653589793 {
				if (*rec).nodep < xnoh {
					(*rec).nodep = (*rec).nodep + 2*3.141592653589793
				} else {
					(*rec).nodep = (*rec).nodep - 2*3.141592653589793
				}
			}
			(*rec).mp = (*rec).mp + pl
			(*rec).argpp = xls - (*rec).mp - cosip*(*rec).nodep
		}
	}
}

// dscom - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:543
/*-----------------------------------------------------------------------------
 *
 *                           procedure dscom
 *
 *  this procedure provides deep space common items used by both the secular
 *    and periodics subroutines.  input is provided as shown. this routine
 *    used to be called dpper, but the functions inside weren't well organized.
 *
 *  author        : david vallado                  719-573-2600   28 jun 2005
 *
 *  inputs        :
 *    epoch       -
 *    ep          - eccentricity
 *    argpp       - argument of perigee
 *    tc          -
 *    inclp       - inclination
 *    nodep       - right ascension of ascending node
 *    np          - mean motion
 *
 *  outputs       :
 *    sinim  , cosim  , sinomm , cosomm , snodm  , cnodm
 *    day         -
 *    e3          -
 *    ee2         -
 *    em          - eccentricity
 *    emsq        - eccentricity squared
 *    gam         -
 *    peo         -
 *    pgho        -
 *    pho         -
 *    pinco       -
 *    plo         -
 *    rtemsq      -
 *    se2, se3         -
 *    sgh2, sgh3, sgh4        -
 *    sh2, sh3, si2, si3, sl2, sl3, sl4         -
 *    s1, s2, s3, s4, s5, s6, s7          -
 *    ss1, ss2, ss3, ss4, ss5, ss6, ss7, sz1, sz2, sz3         -
 *    sz11, sz12, sz13, sz21, sz22, sz23, sz31, sz32, sz33        -
 *    xgh2, xgh3, xgh4, xh2, xh3, xi2, xi3, xl2, xl3, xl4         -
 *    nm          - mean motion
 *    z1, z2, z3, z11, z12, z13, z21, z22, z23, z31, z32, z33         -
 *    zmol        -
 *    zmos        -
 *
 *  locals        :
 *    a1, a2, a3, a4, a5, a6, a7, a8, a9, a10         -
 *    betasq      -
 *    cc          -
 *    ctem, stem        -
 *    x1, x2, x3, x4, x5, x6, x7, x8          -
 *    xnodce      -
 *    xnoi        -
 *    zcosg  , zsing  , zcosgl , zsingl , zcosh  , zsinh  , zcoshl , zsinhl ,
 *    zcosi  , zsini  , zcosil , zsinil ,
 *    zx          -
 *    zy          -
 *
 *  coupling      :
 *    none.
 *
 *  references    :
 *    hoots, roehrich, norad spacetrack report #3 1980
 *    hoots, norad spacetrack report #6 1986
 *    hoots, schumacher and glover 2004
 *    vallado, crawford, hujsak, kelso  2006
 ----------------------------------------------------------------------------*/ //
/* -------------------------- constants ------------------------- */ //
/* --------------------- local variables ------------------------ */ //
/* ----------------- initialize lunar solar terms --------------- */ //
/* ------------------------- do solar terms --------------------- */ //
/* ----------------------- do lunar terms ------------------- */ //
/* ------------------------ do solar terms ---------------------- */ //
/* ------------------------ do lunar terms ---------------------- */ //
// dscom
//
func dscom(epoch float64, ep float64, argpp float64, tc float64, inclp float64, nodep float64, np float64, rec *ElsetRec) {
	var zes float64 = 0.01675
	var zel float64 = 0.0549
	var c1ss float64 = 2.9864797e-06
	var c1l float64 = 4.7968065e-07
	var zsinis float64 = 0.39785416
	var zcosis float64 = 0.91744867
	var zcosgs float64 = 0.1945905
	var zsings float64 = -0.98088458
	var lsflg int64
	var a1 float64
	var a2 float64
	var a3 float64
	var a4 float64
	var a5 float64
	var a6 float64
	var a7 float64
	var a8 float64
	var a9 float64
	var a10 float64
	var betasq float64
	var cc float64
	var ctem float64
	var stem float64
	var x1 float64
	var x2 float64
	var x3 float64
	var x4 float64
	var x5 float64
	var x6 float64
	var x7 float64
	var x8 float64
	var xnodce float64
	var xnoi float64
	var zcosg float64
	var zcosgl float64
	var zcosh float64
	var zcoshl float64
	var zcosi float64
	var zcosil float64
	var zsing float64
	var zsingl float64
	var zsinh float64
	var zsinhl float64
	var zsini float64
	var zsinil float64
	var zx float64
	var zy float64
	(*rec).nm = np
	(*rec).em = ep
	(*rec).snodm = math.Sin(nodep)
	(*rec).cnodm = math.Cos(nodep)
	(*rec).sinomm = math.Sin(argpp)
	(*rec).cosomm = math.Cos(argpp)
	(*rec).sinim = math.Sin(inclp)
	(*rec).cosim = math.Cos(inclp)
	(*rec).emsq = (*rec).em * (*rec).em
	betasq = 1 - (*rec).emsq
	(*rec).rtemsq = math.Sqrt(betasq)
	(*rec).peo = 0
	(*rec).pinco = 0
	(*rec).plo = 0
	(*rec).pgho = 0
	(*rec).pho = 0
	(*rec).day = epoch + 18261.5 + tc/1440
	xnodce = math.Mod(4.523602-0.00092422029*(*rec).day, (2 * 3.141592653589793))
	stem = math.Sin(xnodce)
	ctem = math.Cos(xnodce)
	zcosil = 0.91375164 - 0.03568096*ctem
	zsinil = math.Sqrt(1 - zcosil*zcosil)
	zsinhl = 0.089683511 * stem / zsinil
	zcoshl = math.Sqrt(1 - zsinhl*zsinhl)
	(*rec).gam = 5.8351514 + 0.001944368*(*rec).day
	zx = 0.39785416 * stem / zsinil
	zy = zcoshl*ctem + 0.91744867*zsinhl*stem
	zx = math.Atan2(zx, zy)
	zx = (*rec).gam + zx - xnodce
	zcosgl = math.Cos(zx)
	zsingl = math.Sin(zx)
	zcosg = zcosgs
	zsing = zsings
	zcosi = zcosis
	zsini = zsinis
	zcosh = (*rec).cnodm
	zsinh = (*rec).snodm
	cc = c1ss
	xnoi = 1.0 / (*rec).nm
	for lsflg = int64(1); lsflg <= int64(2); lsflg++ {
		a1 = zcosg*zcosh + zsing*zcosi*zsinh
		a3 = -zsing*zcosh + zcosg*zcosi*zsinh
		a7 = -zcosg*zsinh + zsing*zcosi*zcosh
		a8 = zsing * zsini
		a9 = zsing*zsinh + zcosg*zcosi*zcosh
		a10 = zcosg * zsini
		a2 = (*rec).cosim*a7 + (*rec).sinim*a8
		a4 = (*rec).cosim*a9 + (*rec).sinim*a10
		a5 = -(*rec).sinim*a7 + (*rec).cosim*a8
		a6 = -(*rec).sinim*a9 + (*rec).cosim*a10
		x1 = a1*(*rec).cosomm + a2*(*rec).sinomm
		x2 = a3*(*rec).cosomm + a4*(*rec).sinomm
		x3 = -a1*(*rec).sinomm + a2*(*rec).cosomm
		x4 = -a3*(*rec).sinomm + a4*(*rec).cosomm
		x5 = a5 * (*rec).sinomm
		x6 = a6 * (*rec).sinomm
		x7 = a5 * (*rec).cosomm
		x8 = a6 * (*rec).cosomm
		(*rec).z31 = 12*x1*x1 - 3*x3*x3
		(*rec).z32 = 24*x1*x2 - 6*x3*x4
		(*rec).z33 = 12*x2*x2 - 3*x4*x4
		(*rec).z1 = 3*(a1*a1+a2*a2) + (*rec).z31*(*rec).emsq
		(*rec).z2 = 6*(a1*a3+a2*a4) + (*rec).z32*(*rec).emsq
		(*rec).z3 = 3*(a3*a3+a4*a4) + (*rec).z33*(*rec).emsq
		(*rec).z11 = -6*a1*a5 + (*rec).emsq*(-24*x1*x7-6*x3*x5)
		(*rec).z12 = -6*(a1*a6+a3*a5) + (*rec).emsq*(-24*(x2*x7+x1*x8)-6*(x3*x6+x4*x5))
		(*rec).z13 = -6*a3*a6 + (*rec).emsq*(-24*x2*x8-6*x4*x6)
		(*rec).z21 = 6*a2*a5 + (*rec).emsq*(24*x1*x5-6*x3*x7)
		(*rec).z22 = 6*(a4*a5+a2*a6) + (*rec).emsq*(24*(x2*x5+x1*x6)-6*(x4*x7+x3*x8))
		(*rec).z23 = 6*a4*a6 + (*rec).emsq*(24*x2*x6-6*x4*x8)
		(*rec).z1 = (*rec).z1 + (*rec).z1 + betasq*(*rec).z31
		(*rec).z2 = (*rec).z2 + (*rec).z2 + betasq*(*rec).z32
		(*rec).z3 = (*rec).z3 + (*rec).z3 + betasq*(*rec).z33
		(*rec).s3 = cc * xnoi
		(*rec).s2 = -0.5 * (*rec).s3 / (*rec).rtemsq
		(*rec).s4 = (*rec).s3 * (*rec).rtemsq
		(*rec).s1 = -15 * (*rec).em * (*rec).s4
		(*rec).s5 = x1*x3 + x2*x4
		(*rec).s6 = x2*x3 + x1*x4
		(*rec).s7 = x2*x4 - x1*x3
		if lsflg == int64(1) {
			(*rec).ss1 = (*rec).s1
			(*rec).ss2 = (*rec).s2
			(*rec).ss3 = (*rec).s3
			(*rec).ss4 = (*rec).s4
			(*rec).ss5 = (*rec).s5
			(*rec).ss6 = (*rec).s6
			(*rec).ss7 = (*rec).s7
			(*rec).sz1 = (*rec).z1
			(*rec).sz2 = (*rec).z2
			(*rec).sz3 = (*rec).z3
			(*rec).sz11 = (*rec).z11
			(*rec).sz12 = (*rec).z12
			(*rec).sz13 = (*rec).z13
			(*rec).sz21 = (*rec).z21
			(*rec).sz22 = (*rec).z22
			(*rec).sz23 = (*rec).z23
			(*rec).sz31 = (*rec).z31
			(*rec).sz32 = (*rec).z32
			(*rec).sz33 = (*rec).z33
			zcosg = zcosgl
			zsing = zsingl
			zcosi = zcosil
			zsini = zsinil
			zcosh = zcoshl*(*rec).cnodm + zsinhl*(*rec).snodm
			zsinh = (*rec).snodm*zcoshl - (*rec).cnodm*zsinhl
			cc = c1l
		}
	}
	(*rec).zmol = math.Mod(4.7199672+0.2299715*(*rec).day-(*rec).gam, (2 * 3.141592653589793))
	(*rec).zmos = math.Mod(6.2565837+0.017201977*(*rec).day, (2 * 3.141592653589793))
	(*rec).se2 = 2 * (*rec).ss1 * (*rec).ss6
	(*rec).se3 = 2 * (*rec).ss1 * (*rec).ss7
	(*rec).si2 = 2 * (*rec).ss2 * (*rec).sz12
	(*rec).si3 = 2 * (*rec).ss2 * ((*rec).sz13 - (*rec).sz11)
	(*rec).sl2 = -2 * (*rec).ss3 * (*rec).sz2
	(*rec).sl3 = -2 * (*rec).ss3 * ((*rec).sz3 - (*rec).sz1)
	(*rec).sl4 = -2 * (*rec).ss3 * (-21 - 9*(*rec).emsq) * zes
	(*rec).sgh2 = 2 * (*rec).ss4 * (*rec).sz32
	(*rec).sgh3 = 2 * (*rec).ss4 * ((*rec).sz33 - (*rec).sz31)
	(*rec).sgh4 = -18 * (*rec).ss4 * zes
	(*rec).sh2 = -2 * (*rec).ss2 * (*rec).sz22
	(*rec).sh3 = -2 * (*rec).ss2 * ((*rec).sz23 - (*rec).sz21)
	(*rec).ee2 = 2 * (*rec).s1 * (*rec).s6
	(*rec).e3 = 2 * (*rec).s1 * (*rec).s7
	(*rec).xi2 = 2 * (*rec).s2 * (*rec).z12
	(*rec).xi3 = 2 * (*rec).s2 * ((*rec).z13 - (*rec).z11)
	(*rec).xl2 = -2 * (*rec).s3 * (*rec).z2
	(*rec).xl3 = -2 * (*rec).s3 * ((*rec).z3 - (*rec).z1)
	(*rec).xl4 = -2 * (*rec).s3 * (-21 - 9*(*rec).emsq) * zel
	(*rec).xgh2 = 2 * (*rec).s4 * (*rec).z32
	(*rec).xgh3 = 2 * (*rec).s4 * ((*rec).z33 - (*rec).z31)
	(*rec).xgh4 = -18 * (*rec).s4 * zel
	(*rec).xh2 = -2 * (*rec).s2 * (*rec).z22
	(*rec).xh3 = -2 * (*rec).s2 * ((*rec).z23 - (*rec).z21)
}

// dsinit - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:806
/*-----------------------------------------------------------------------------
 *
 *                           procedure dsinit
 *
 *  this procedure provides deep space contributions to mean motion dot due
 *    to geopotential resonance with half day and one day orbits.
 *
 *  author        : david vallado                  719-573-2600   28 jun 2005
 *
 *  inputs        :
 *    xke         - reciprocal of tumin
 *    cosim, sinim-
 *    emsq        - eccentricity squared
 *    argpo       - argument of perigee
 *    s1, s2, s3, s4, s5      -
 *    ss1, ss2, ss3, ss4, ss5 -
 *    sz1, sz3, sz11, sz13, sz21, sz23, sz31, sz33 -
 *    t           - time
 *    tc          -
 *    gsto        - greenwich sidereal time                   rad
 *    mo          - mean anomaly
 *    mdot        - mean anomaly dot (rate)
 *    no          - mean motion
 *    nodeo       - right ascension of ascending node
 *    nodedot     - right ascension of ascending node dot (rate)
 *    xpidot      -
 *    z1, z3, z11, z13, z21, z23, z31, z33 -
 *    eccm        - eccentricity
 *    argpm       - argument of perigee
 *    inclm       - inclination
 *    mm          - mean anomaly
 *    xn          - mean motion
 *    nodem       - right ascension of ascending node
 *
 *  outputs       :
 *    em          - eccentricity
 *    argpm       - argument of perigee
 *    inclm       - inclination
 *    mm          - mean anomaly
 *    nm          - mean motion
 *    nodem       - right ascension of ascending node
 *    irez        - flag for resonance           0-none, 1-one day, 2-half day
 *    atime       -
 *    d2201, d2211, d3210, d3222, d4410, d4422, d5220, d5232, d5421, d5433    -
 *    dedt        -
 *    didt        -
 *    dmdt        -
 *    dndt        -
 *    dnodt       -
 *    domdt       -
 *    del1, del2, del3        -
 *    ses  , sghl , sghs , sgs  , shl  , shs  , sis  , sls
 *    theta       -
 *    xfact       -
 *    xlamo       -
 *    xli         -
 *    xni
 *
 *  locals        :
 *    ainv2       -
 *    aonv        -
 *    cosisq      -
 *    eoc         -
 *    f220, f221, f311, f321, f322, f330, f441, f442, f522, f523, f542, f543  -
 *    g200, g201, g211, g300, g310, g322, g410, g422, g520, g521, g532, g533  -
 *    sini2       -
 *    temp        -
 *    temp1       -
 *    theta       -
 *    xno2        -
 *
 *  coupling      :
 *    getgravconst- no longer used
 *
 *  references    :
 *    hoots, roehrich, norad spacetrack report #3 1980
 *    hoots, norad spacetrack report #6 1986
 *    hoots, schumacher and glover 2004
 *    vallado, crawford, hujsak, kelso  2006
 ----------------------------------------------------------------------------*/ //
/* --------------------- local variables ------------------------ */ //
// this equates to 7.29211514668855e-5 rad/sec
// sgp4fix identify constants and allow alternate values
// just xke is used here so pass it in rather than have multiple calls
// getgravconst( whichconst, tumin, mu, radiusearthkm, xke, j2, j3, j4, j3oj2 );
/* -------------------- deep space initialization ------------ */ //
/* ------------------------ do solar terms ------------------- */ //
// sgp4fix for 180 deg incl
/* ------------------------- do lunar terms ------------------ */ //
// sgp4fix for 180 deg incl
/* ----------- calculate deep space resonance effects -------- */ //
//   sgp4fix for negative inclinations
//   the following if statement should be commented out
//if (inclm < 0.0)
//  {
//    inclm  = -inclm;
//    argpm  = argpm - pi;
//    nodem = nodem + pi;
//  }
/* -------------- initialize the resonance terms ------------- */ //
/* ---------- geopotential resonance for 12 hour orbits ------ */ //
/* ---------------- synchronous resonance terms -------------- */ //
/* ------------ for sgp4, initialize the integrator ---------- */ //
// dsinit
//
func dsinit(tc float64, xpidot float64, rec *ElsetRec) {
	var ainv2 float64
	var aonv float64 = 0
	var cosisq float64
	var eoc float64
	var f220 float64
	var f221 float64
	var f311 float64
	var f321 float64
	var f322 float64
	var f330 float64
	var f441 float64
	var f442 float64
	var f522 float64
	var f523 float64
	var f542 float64
	var f543 float64
	var g200 float64
	var g201 float64
	var g211 float64
	var g300 float64
	var g310 float64
	var g322 float64
	var g410 float64
	var g422 float64
	var g520 float64
	var g521 float64
	var g532 float64
	var g533 float64
	var ses float64
	var sgs float64
	var sghl float64
	var sghs float64
	var shs float64
	var shll float64
	var sis float64
	var sini2 float64
	var sls float64
	var temp float64
	var temp1 float64
	var theta float64
	var xno2 float64
	var q22 float64
	var q31 float64
	var q33 float64
	var root22 float64
	var root44 float64
	var root54 float64
	var rptim float64
	var root32 float64
	var root52 float64
	var x2o3 float64
	var znl float64
	var emo float64
	var zns float64
	var emsqo float64
	q22 = 1.7891679e-06
	q31 = 2.1460748e-06
	q33 = 2.2123015e-07
	root22 = 1.7891679e-6
	root44 = 7.3636953e-9
	root54 = 2.1765803e-9
	rptim = 0.0043752690880113
	root32 = 3.7393792e-7
	root52 = 1.1428639e-7
	x2o3 = 2.0 / 3.0
	znl = 0.00015835218
	zns = 1.19459e-05
	(*rec).irez = int64(0)
	if (*rec).nm < 0.0052359877 && (*rec).nm > 0.0034906585 {
		(*rec).irez = int64(1)
	}
	if (*rec).nm >= 0.00826 && (*rec).nm <= 0.00924 && (*rec).em >= 0.5 {
		(*rec).irez = int64(2)
	}
	ses = (*rec).ss1 * zns * (*rec).ss5
	sis = (*rec).ss2 * zns * ((*rec).sz11 + (*rec).sz13)
	sls = -zns * (*rec).ss3 * ((*rec).sz1 + (*rec).sz3 - 14 - 6*(*rec).emsq)
	sghs = (*rec).ss4 * zns * ((*rec).sz31 + (*rec).sz33 - 6)
	shs = -zns * (*rec).ss2 * ((*rec).sz21 + (*rec).sz23)
	if (*rec).inclm < 0.052359877 || (*rec).inclm > 3.141592653589793-0.052359877 {
		shs = 0
	}
	if (*rec).sinim != 0 {
		shs = shs / (*rec).sinim
	}
	sgs = sghs - (*rec).cosim*shs
	(*rec).dedt = ses + (*rec).s1*znl*(*rec).s5
	(*rec).didt = sis + (*rec).s2*znl*((*rec).z11+(*rec).z13)
	(*rec).dmdt = sls - znl*(*rec).s3*((*rec).z1+(*rec).z3-14-6*(*rec).emsq)
	sghl = (*rec).s4 * znl * ((*rec).z31 + (*rec).z33 - 6)
	shll = -znl * (*rec).s2 * ((*rec).z21 + (*rec).z23)
	if (*rec).inclm < 0.052359877 || (*rec).inclm > 3.141592653589793-0.052359877 {
		shll = 0
	}
	(*rec).domdt = sgs + sghl
	(*rec).dnodt = shs
	if (*rec).sinim != 0 {
		(*rec).domdt = (*rec).domdt - (*rec).cosim/(*rec).sinim*shll
		(*rec).dnodt = (*rec).dnodt + shll/(*rec).sinim
	}
	(*rec).dndt = 0
	theta = math.Mod((*rec).gsto+tc*rptim, (2 * 3.141592653589793))
	(*rec).em = (*rec).em + (*rec).dedt*(*rec).t
	(*rec).inclm = (*rec).inclm + (*rec).didt*(*rec).t
	(*rec).argpm = (*rec).argpm + (*rec).domdt*(*rec).t
	(*rec).nodem = (*rec).nodem + (*rec).dnodt*(*rec).t
	(*rec).mm = (*rec).mm + (*rec).dmdt*(*rec).t
	if (*rec).irez != int64(0) {
		aonv = math.Pow((*rec).nm/(*rec).xke, x2o3)
		if (*rec).irez == int64(2) {
			cosisq = (*rec).cosim * (*rec).cosim
			emo = (*rec).em
			(*rec).em = (*rec).ecco
			emsqo = (*rec).emsq
			(*rec).emsq = (*rec).eccsq
			eoc = (*rec).em * (*rec).emsq
			g201 = -0.306 - ((*rec).em-0.64)*0.44
			if (*rec).em <= 0.65 {
				g211 = 3.616 - 13.247*(*rec).em + 16.29*(*rec).emsq
				g310 = -19.302 + 117.39*(*rec).em - 228.419*(*rec).emsq + 156.591*eoc
				g322 = -18.9068 + 109.7927*(*rec).em - 214.6334*(*rec).emsq + 146.5816*eoc
				g410 = -41.122 + 242.694*(*rec).em - 471.094*(*rec).emsq + 313.953*eoc
				g422 = -146.407 + 841.88*(*rec).em - 1629.014*(*rec).emsq + 1083.435*eoc
				g520 = -532.114 + 3017.977*(*rec).em - 5740.032*(*rec).emsq + 3708.276*eoc
			} else {
				g211 = -72.099 + 331.819*(*rec).em - 508.738*(*rec).emsq + 266.724*eoc
				g310 = -346.844 + 1582.851*(*rec).em - 2415.925*(*rec).emsq + 1246.113*eoc
				g322 = -342.585 + 1554.908*(*rec).em - 2366.899*(*rec).emsq + 1215.972*eoc
				g410 = -1052.797 + 4758.686*(*rec).em - 7193.992*(*rec).emsq + 3651.957*eoc
				g422 = -3581.69 + 16178.11*(*rec).em - 24462.77*(*rec).emsq + 12422.52*eoc
				if (*rec).em > 0.715 {
					g520 = -5149.66 + 29936.92*(*rec).em - 54087.36*(*rec).emsq + 31324.56*eoc
				} else {
					g520 = 1464.74 - 4664.75*(*rec).em + 3763.64*(*rec).emsq
				}
			}
			if (*rec).em < 0.7 {
				g533 = -919.2277 + 4988.61*(*rec).em - 9064.77*(*rec).emsq + 5542.21*eoc
				g521 = -822.71072 + 4568.6173*(*rec).em - 8491.4146*(*rec).emsq + 5337.524*eoc
				g532 = -853.666 + 4690.25*(*rec).em - 8624.77*(*rec).emsq + 5341.4*eoc
			} else {
				g533 = -37995.78 + 161616.52*(*rec).em - 229838.2*(*rec).emsq + 109377.94*eoc
				g521 = -51752.104 + 218913.95*(*rec).em - 309468.16*(*rec).emsq + 146349.42*eoc
				g532 = -40023.88 + 170470.89*(*rec).em - 242699.48*(*rec).emsq + 115605.82*eoc
			}
			sini2 = (*rec).sinim * (*rec).sinim
			f220 = 0.75 * (1 + 2*(*rec).cosim + cosisq)
			f221 = 1.5 * sini2
			f321 = 1.875 * (*rec).sinim * (1 - 2*(*rec).cosim - 3*cosisq)
			f322 = -1.875 * (*rec).sinim * (1 + 2*(*rec).cosim - 3*cosisq)
			f441 = 35 * sini2 * f220
			f442 = 39.375 * sini2 * sini2
			f522 = 9.84375 * (*rec).sinim * (sini2*(1-2*(*rec).cosim-5*cosisq) + 0.33333333*(-2+4*(*rec).cosim+6*cosisq))
			f523 = (*rec).sinim * (4.92187512*sini2*(-2-4*(*rec).cosim+10*cosisq) + 6.56250012*(1+2*(*rec).cosim-3*cosisq))
			f542 = 29.53125 * (*rec).sinim * (2 - 8*(*rec).cosim + cosisq*(-12+8*(*rec).cosim+10*cosisq))
			f543 = 29.53125 * (*rec).sinim * (-2 - 8*(*rec).cosim + cosisq*(12+8*(*rec).cosim-10*cosisq))
			xno2 = (*rec).nm * (*rec).nm
			ainv2 = aonv * aonv
			temp1 = 3 * xno2 * ainv2
			temp = temp1 * root22
			(*rec).d2201 = temp * f220 * g201
			(*rec).d2211 = temp * f221 * g211
			temp1 = temp1 * aonv
			temp = temp1 * root32
			(*rec).d3210 = temp * f321 * g310
			(*rec).d3222 = temp * f322 * g322
			temp1 = temp1 * aonv
			temp = 2 * temp1 * root44
			(*rec).d4410 = temp * f441 * g410
			(*rec).d4422 = temp * f442 * g422
			temp1 = temp1 * aonv
			temp = temp1 * root52
			(*rec).d5220 = temp * f522 * g520
			(*rec).d5232 = temp * f523 * g532
			temp = 2 * temp1 * root54
			(*rec).d5421 = temp * f542 * g521
			(*rec).d5433 = temp * f543 * g533
			(*rec).xlamo = math.Mod((*rec).mo+(*rec).nodeo+(*rec).nodeo-theta-theta, (2 * 3.141592653589793))
			(*rec).xfact = (*rec).mdot + (*rec).dmdt + 2*((*rec).nodedot+(*rec).dnodt-rptim) - (*rec).no_unkozai
			(*rec).em = emo
			(*rec).emsq = emsqo
		}
		if (*rec).irez == int64(1) {
			g200 = 1 + (*rec).emsq*(-2.5+0.8125*(*rec).emsq)
			g310 = 1 + 2*(*rec).emsq
			g300 = 1 + (*rec).emsq*(-6+6.60937*(*rec).emsq)
			f220 = 0.75 * (1 + (*rec).cosim) * (1 + (*rec).cosim)
			f311 = 0.9375*(*rec).sinim*(*rec).sinim*(1+3*(*rec).cosim) - 0.75*(1+(*rec).cosim)
			f330 = 1 + (*rec).cosim
			f330 = 1.875 * f330 * f330 * f330
			(*rec).del1 = 3 * (*rec).nm * (*rec).nm * aonv * aonv
			(*rec).del2 = 2 * (*rec).del1 * f220 * g200 * q22
			(*rec).del3 = 3 * (*rec).del1 * f330 * g300 * q33 * aonv
			(*rec).del1 = (*rec).del1 * f311 * g310 * q31 * aonv
			(*rec).xlamo = math.Mod((*rec).mo+(*rec).nodeo+(*rec).argpo-theta, (2 * 3.141592653589793))
			(*rec).xfact = (*rec).mdot + xpidot - rptim + (*rec).dmdt + (*rec).domdt + (*rec).dnodt - (*rec).no_unkozai
		}
		(*rec).xli = (*rec).xlamo
		(*rec).xni = (*rec).no_unkozai
		(*rec).atime = 0
		(*rec).nm = (*rec).no_unkozai + (*rec).dndt
	}
}

// dspace - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:1086
/*-----------------------------------------------------------------------------
 *
 *                           procedure dspace
 *
 *  this procedure provides deep space contributions to mean elements for
 *    perturbing third body.  these effects have been averaged over one
 *    revolution of the sun and moon.  for earth resonance effects, the
 *    effects have been averaged over no revolutions of the satellite.
 *    (mean motion)
 *
 *  author        : david vallado                  719-573-2600   28 jun 2005
 *
 *  inputs        :
 *    d2201, d2211, d3210, d3222, d4410, d4422, d5220, d5232, d5421, d5433 -
 *    dedt        -
 *    del1, del2, del3  -
 *    didt        -
 *    dmdt        -
 *    dnodt       -
 *    domdt       -
 *    irez        - flag for resonance           0-none, 1-one day, 2-half day
 *    argpo       - argument of perigee
 *    argpdot     - argument of perigee dot (rate)
 *    t           - time
 *    tc          -
 *    gsto        - gst
 *    xfact       -
 *    xlamo       -
 *    no          - mean motion
 *    atime       -
 *    em          - eccentricity
 *    ft          -
 *    argpm       - argument of perigee
 *    inclm       - inclination
 *    xli         -
 *    mm          - mean anomaly
 *    xni         - mean motion
 *    nodem       - right ascension of ascending node
 *
 *  outputs       :
 *    atime       -
 *    em          - eccentricity
 *    argpm       - argument of perigee
 *    inclm       - inclination
 *    xli         -
 *    mm          - mean anomaly
 *    xni         -
 *    nodem       - right ascension of ascending node
 *    dndt        -
 *    nm          - mean motion
 *
 *  locals        :
 *    delt        -
 *    ft          -
 *    theta       -
 *    x2li        -
 *    x2omi       -
 *    xl          -
 *    xldot       -
 *    xnddt       -
 *    xndt        -
 *    xomi        -
 *
 *  coupling      :
 *    none        -
 *
 *  references    :
 *    hoots, roehrich, norad spacetrack report #3 1980
 *    hoots, norad spacetrack report #6 1986
 *    hoots, schumacher and glover 2004
 *    vallado, crawford, hujsak, kelso  2006
 ----------------------------------------------------------------------------*/ //
// this equates to 7.29211514668855e-5 rad/sec
/* ----------- calculate deep space resonance effects ----------- */ //
//   sgp4fix for negative inclinations
//   the following if statement should be commented out
//  if (inclm < 0.0)
// {
//    inclm = -inclm;
//    argpm = argpm - pi;
//    nodem = nodem + pi;
//  }
/* - update resonances : numerical (euler-maclaurin) integration - */ //
/* ------------------------- epoch restart ----------------------  */ //
//   sgp4fix for propagator problems
//   the following integration works for negative time steps and periods
//   the specific changes are unknown because the original code was so convoluted
// sgp4fix take out atime = 0.0 and fix for faster operation
// sgp4fix streamline check
// sgp4fix move check outside loop
// added for do loop
/* ------------------- dot terms calculated ------------- */ //
/* ----------- near - synchronous resonance terms ------- */ //
/* --------- near - half-day resonance terms -------- */ //
/* ----------------------- integrator ------------------- */ //
// sgp4fix move end checks to end of routine
// exit here
// while iretn = 381
// dsspace
//
func dspace(tc float64, rec *ElsetRec) {
	var iretn int64
	var delt float64
	var ft float64
	var theta float64
	var x2li float64
	var x2omi float64
	var xl float64
	var xldot float64
	var xnddt float64
	var xndt float64
	var xomi float64
	var g22 float64
	var g32 float64
	var g44 float64
	var g52 float64
	var g54 float64
	var fasx2 float64
	var fasx4 float64
	var fasx6 float64
	var rptim float64
	var step2 float64
	var stepn float64
	var stepp float64
	xndt = float64(int64(0))
	xnddt = float64(int64(0))
	xldot = float64(int64(0))
	fasx2 = 0.13130908
	fasx4 = 2.8843198
	fasx6 = 0.37448087
	g22 = 5.7686396
	g32 = 0.95240898
	g44 = 1.8014998
	g52 = 1.050833
	g54 = 4.4108898
	rptim = 0.0043752690880113
	stepp = 720
	stepn = -720
	step2 = 259200
	(*rec).dndt = 0
	theta = math.Mod((*rec).gsto+tc*rptim, (2 * 3.141592653589793))
	(*rec).em = (*rec).em + (*rec).dedt*(*rec).t
	(*rec).inclm = (*rec).inclm + (*rec).didt*(*rec).t
	(*rec).argpm = (*rec).argpm + (*rec).domdt*(*rec).t
	(*rec).nodem = (*rec).nodem + (*rec).dnodt*(*rec).t
	(*rec).mm = (*rec).mm + (*rec).dmdt*(*rec).t
	ft = 0
	if (*rec).irez != int64(0) {
		if (*rec).atime == 0 || (*rec).t*(*rec).atime <= 0 || math.Abs((*rec).t) < math.Abs((*rec).atime) {
			(*rec).atime = 0
			(*rec).xni = (*rec).no_unkozai
			(*rec).xli = (*rec).xlamo
		}
		if (*rec).t > 0 {
			delt = stepp
		} else {
			delt = stepn
		}
		iretn = int64(381)
		for iretn == int64(381) {
			if (*rec).irez != int64(2) {
				xndt = (*rec).del1*math.Sin((*rec).xli-fasx2) + (*rec).del2*math.Sin(2*((*rec).xli-fasx4)) + (*rec).del3*math.Sin(3*((*rec).xli-fasx6))
				xldot = (*rec).xni + (*rec).xfact
				xnddt = (*rec).del1*math.Cos((*rec).xli-fasx2) + 2*(*rec).del2*math.Cos(2*((*rec).xli-fasx4)) + 3*(*rec).del3*math.Cos(3*((*rec).xli-fasx6))
				xnddt = xnddt * xldot
			} else {
				xomi = (*rec).argpo + (*rec).argpdot*(*rec).atime
				x2omi = xomi + xomi
				x2li = (*rec).xli + (*rec).xli
				xndt = (*rec).d2201*math.Sin(x2omi+(*rec).xli-g22) + (*rec).d2211*math.Sin((*rec).xli-g22) + (*rec).d3210*math.Sin(xomi+(*rec).xli-g32) + (*rec).d3222*math.Sin(-xomi+(*rec).xli-g32) + (*rec).d4410*math.Sin(x2omi+x2li-g44) + (*rec).d4422*math.Sin(x2li-g44) + (*rec).d5220*math.Sin(xomi+(*rec).xli-g52) + (*rec).d5232*math.Sin(-xomi+(*rec).xli-g52) + (*rec).d5421*math.Sin(xomi+x2li-g54) + (*rec).d5433*math.Sin(-xomi+x2li-g54)
				xldot = (*rec).xni + (*rec).xfact
				xnddt = (*rec).d2201*math.Cos(x2omi+(*rec).xli-g22) + (*rec).d2211*math.Cos((*rec).xli-g22) + (*rec).d3210*math.Cos(xomi+(*rec).xli-g32) + (*rec).d3222*math.Cos(-xomi+(*rec).xli-g32) + (*rec).d5220*math.Cos(xomi+(*rec).xli-g52) + (*rec).d5232*math.Cos(-xomi+(*rec).xli-g52) + 2*((*rec).d4410*math.Cos(x2omi+x2li-g44)+(*rec).d4422*math.Cos(x2li-g44)+(*rec).d5421*math.Cos(xomi+x2li-g54)+(*rec).d5433*math.Cos(-xomi+x2li-g54))
				xnddt = xnddt * xldot
			}
			if math.Abs((*rec).t-(*rec).atime) >= stepp {
				iretn = int64(381)
			} else {
				ft = (*rec).t - (*rec).atime
				iretn = int64(0)
			}
			if iretn == int64(381) {
				(*rec).xli = (*rec).xli + xldot*delt + xndt*step2
				(*rec).xni = (*rec).xni + xndt*delt + xnddt*step2
				(*rec).atime = (*rec).atime + delt
			}
		}
		(*rec).nm = (*rec).xni + xndt*ft + xnddt*ft*ft*0.5
		xl = (*rec).xli + xldot*ft + xndt*ft*ft*0.5
		if (*rec).irez != int64(1) {
			(*rec).mm = xl - 2*(*rec).nodem + 2*theta
			(*rec).dndt = (*rec).nm - (*rec).no_unkozai
		} else {
			(*rec).mm = xl - (*rec).nodem - (*rec).argpm + theta
			(*rec).dndt = (*rec).nm - (*rec).no_unkozai
		}
		(*rec).nm = (*rec).no_unkozai + (*rec).dndt
	}
}

// initl - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:1278
/*-----------------------------------------------------------------------------
 *
 *                           procedure initl
 *
 *  this procedure initializes the spg4 propagator. all the initialization is
 *    consolidated here instead of having multiple loops inside other routines.
 *
 *  author        : david vallado                  719-573-2600   28 jun 2005
 *
 *  inputs        :
 *    satn        - satellite number - not needed, placed in satrec
 *    xke         - reciprocal of tumin
 *    j2          - j2 zonal harmonic
 *    ecco        - eccentricity                           0.0 - 1.0
 *    epoch       - epoch time in days from jan 0, 1950. 0 hr
 *    inclo       - inclination of satellite
 *    no          - mean motion of satellite
 *
 *  outputs       :
 *    ainv        - 1.0 / a
 *    ao          - semi major axis
 *    con41       -
 *    con42       - 1.0 - 5.0 cos(i)
 *    cosio       - cosine of inclination
 *    cosio2      - cosio squared
 *    eccsq       - eccentricity squared
 *    method      - flag for deep space                    'd', 'n'
 *    omeosq      - 1.0 - ecco * ecco
 *    posq        - semi-parameter squared
 *    rp          - radius of perigee
 *    rteosq      - square root of (1.0 - ecco*ecco)
 *    sinio       - sine of inclination
 *    gsto        - gst at time of observation               rad
 *    no          - mean motion of satellite
 *
 *  locals        :
 *    ak          -
 *    d1          -
 *    del         -
 *    adel        -
 *    po          -
 *
 *  coupling      :
 *    getgravconst- no longer used
 *    gstime      - find greenwich sidereal time from the julian date
 *
 *  references    :
 *    hoots, roehrich, norad spacetrack report #3 1980
 *    hoots, norad spacetrack report #6 1986
 *    hoots, schumacher and glover 2004
 *    vallado, crawford, hujsak, kelso  2006
 ----------------------------------------------------------------------------*/ //
/* --------------------- local variables ------------------------ */ //
// sgp4fix use old way of finding gst
/* ----------------------- earth constants ---------------------- */ //
// sgp4fix identify constants and allow alternate values
// only xke and j2 are used here so pass them in directly
// getgravconst( whichconst, tumin, mu, radiusearthkm, xke, j2, j3, j4, j3oj2 );
/* ------------- calculate auxillary epoch quantities ---------- */ //
/* ------------------ un-kozai the mean motion ----------------- */ //
// sgp4fix modern approach to finding sidereal time
//   if (opsmode == 'a')
//      {
// sgp4fix use old way of finding gst
// count integer number of days from 0 jan 1970
// find greenwich location at epoch
//    }
//    else
// initl
//
func initl(epoch float64, rec *ElsetRec) {
	var ak float64
	var d1 float64
	var del float64
	var adel float64
	var po float64
	var x2o3 float64
	var ds70 float64
	var ts70 float64
	var tfrac float64
	var c1 float64
	var thgr70 float64
	var fk5r float64
	var c1p2p float64
	x2o3 = 2.0 / 3.0
	(*rec).eccsq = (*rec).ecco * (*rec).ecco
	(*rec).omeosq = 1 - (*rec).eccsq
	(*rec).rteosq = math.Sqrt((*rec).omeosq)
	(*rec).cosio = math.Cos((*rec).inclo)
	(*rec).cosio2 = (*rec).cosio * (*rec).cosio
	ak = math.Pow((*rec).xke/(*rec).no_kozai, x2o3)
	d1 = 0.75 * (*rec).j2 * (3.0*(*rec).cosio2 - 1.0) / ((*rec).rteosq * (*rec).omeosq)
	del = d1 / (ak * ak)
	adel = ak * (1 - del*del - del*(1.0/3.0+134*del*del/81))
	del = d1 / (adel * adel)
	(*rec).no_unkozai = (*rec).no_kozai / (1 + del)
	(*rec).ao = math.Pow((*rec).xke/(*rec).no_unkozai, x2o3)
	(*rec).sinio = math.Sin((*rec).inclo)
	po = (*rec).ao * (*rec).omeosq
	(*rec).con42 = 1.0 - 5.0*(*rec).cosio2
	(*rec).con41 = -(*rec).con42 - (*rec).cosio2 - (*rec).cosio2
	(*rec).ainv = 1.0 / (*rec).ao
	(*rec).posq = po * po
	(*rec).rp = (*rec).ao * (1.0 - (*rec).ecco)
	(*rec).method = 'n'
	ts70 = epoch - 7305.0
	ds70 = math.Floor(ts70 + 1e-08)
	tfrac = ts70 - ds70
	c1 = 0.017202791694070362
	thgr70 = 1.7321343856509375
	fk5r = 5.075514194322695e-15
	c1p2p = c1 + 2*3.141592653589793
	var gsto1 float64 = math.Mod(thgr70+c1*ds70+c1p2p*tfrac+ts70*ts70*fk5r, (2 * 3.141592653589793))
	if gsto1 < 0 {
		gsto1 = gsto1 + 2*3.141592653589793
	}
	(*rec).gsto = gstime(epoch + 2.4332815e+06)
}

// sgp4init - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:1423
/*-----------------------------------------------------------------------------
 *
 *                             procedure sgp4init
 *
 *  this procedure initializes variables for sgp4.
 *
 *  author        : david vallado                  719-573-2600   28 jun 2005
 *
 *  inputs        :
 *    opsmode     - mode of operation afspc or improved 'a', 'i'
 *    whichconst  - which set of constants to use  72, 84
 *    satn        - satellite number
 *    bstar       - sgp4 type drag coefficient              kg/m2er
 *    ecco        - eccentricity
 *    epoch       - epoch time in days from jan 0, 1950. 0 hr
 *    argpo       - argument of perigee (output if ds)
 *    inclo       - inclination
 *    mo          - mean anomaly (output if ds)
 *    no          - mean motion
 *    nodeo       - right ascension of ascending node
 *
 *  outputs       :
 *    satrec      - common values for subsequent calls
 *    return code - non-zero on error.
 *                   1 - mean elements, ecc >= 1.0 or ecc < -0.001 or a < 0.95 er
 *                   2 - mean motion less than 0.0
 *                   3 - pert elements, ecc < 0.0  or  ecc > 1.0
 *                   4 - semi-latus rectum < 0.0
 *                   5 - epoch elements are sub-orbital
 *                   6 - satellite has decayed
 *
 *  locals        :
 *    cnodm  , snodm  , cosim  , sinim  , cosomm , sinomm
 *    cc1sq  , cc2    , cc3
 *    coef   , coef1
 *    cosio4      -
 *    day         -
 *    dndt        -
 *    em          - eccentricity
 *    emsq        - eccentricity squared
 *    eeta        -
 *    etasq       -
 *    gam         -
 *    argpm       - argument of perigee
 *    nodem       -
 *    inclm       - inclination
 *    mm          - mean anomaly
 *    nm          - mean motion
 *    perige      - perigee
 *    pinvsq      -
 *    psisq       -
 *    qzms24      -
 *    rtemsq      -
 *    s1, s2, s3, s4, s5, s6, s7          -
 *    sfour       -
 *    ss1, ss2, ss3, ss4, ss5, ss6, ss7         -
 *    sz1, sz2, sz3
 *    sz11, sz12, sz13, sz21, sz22, sz23, sz31, sz32, sz33        -
 *    tc          -
 *    temp        -
 *    temp1, temp2, temp3       -
 *    tsi         -
 *    xpidot      -
 *    xhdot1      -
 *    z1, z2, z3          -
 *    z11, z12, z13, z21, z22, z23, z31, z32, z33         -
 *
 *  coupling      :
 *    getgravconst-
 *    initl       -
 *    dscom       -
 *    dpper       -
 *    dsinit      -
 *    sgp4        -
 *
 *  references    :
 *    hoots, roehrich, norad spacetrack report #3 1980
 *    hoots, norad spacetrack report #6 1986
 *    hoots, schumacher and glover 2004
 *    vallado, crawford, hujsak, kelso  2006
 ----------------------------------------------------------------------------*/ //
/* --------------------- local variables ------------------------ */ //
/* ------------------------ initialization --------------------- */ //
// sgp4fix divisor for divide by zero check on inclination
// the old check used 1.0 + cos(pi-1.0e-9), but then compared it to
// 1.5 e-12, so the threshold was changed to 1.5e-12 for consistency
/* ----------- set all near earth variables to zero ------------ */ //
/* ----------- set all deep space variables to zero ------------ */ //
/* ------------------------ earth constants ----------------------- */ //
// sgp4fix identify constants and allow alternate values
// this is now the only call for the constants
//-------------------------------------------------------------------------
// single averaged mean elements
/* ------------------------ earth constants ----------------------- */ //
// sgp4fix identify constants and allow alternate values no longer needed
// getgravconst( whichconst, tumin, mu, radiusearthkm, xke, j2, j3, j4, j3oj2 );
// sgp4fix use multiply for speed instead of pow
// sgp4fix remove satn as it is not needed in initl
// sgp4fix remove this check as it is unnecessary
// the mrt check in sgp4 handles decaying satellite cases even if the starting
// condition is below the surface of te earth
//     if (rp < 1.0)
//       {
//         satrec->error = 5;
//       }
/* - for perigees below 156 km, s and qoms2t are altered - */ //
// sgp4fix use multiply for speed instead of pow
// sgp4fix for divide by zero with xinco = 180 deg
// sgp4fix use multiply for speed instead of pow
/* --------------- deep space initialization ------------- */ //
/* ----------- set variables if not deep space ----------- */ //
// if omeosq = 0 ...
/* finally propogate to zero epoch to initialize all others. */ //
// sgp4fix take out check to let satellites process until they are actually below earth surface
//       if(satrec->error == 0)
//sgp4fix return bool_ean. satrec->error contains any error codes
// sgp4init
//
func sgp4init(opsmode byte, satrec *ElsetRec) bool_ {
	var cc1sq float64
	var cc2 float64
	var cc3 float64
	var coef float64
	var coef1 float64
	var cosio4 float64
	var eeta float64
	var etasq float64
	var perige float64
	var pinvsq float64
	var psisq float64
	var qzms24 float64
	var sfour float64
	var tc float64
	var temp float64
	var temp1 float64
	var temp2 float64
	var temp3 float64
	var tsi float64
	var xpidot float64
	var xhdot1 float64
	var qzms2t float64
	var ss float64
	var x2o3 float64
	var r []float64 = make([]float64, 3, 3)
	var v []float64 = make([]float64, 3, 3)
	var delmotemp float64
	var qzms2ttemp float64
	var qzms24temp float64
	var epoch float64 = (*satrec).jdsatepoch + (*satrec).jdsatepochF - 2.4332815e+06
	var temp4 float64 = 1.5e-12
	(*satrec).isimp = int64(0)
	(*satrec).method = 'n'
	(*satrec).aycof = 0
	(*satrec).con41 = 0
	(*satrec).cc1 = 0
	(*satrec).cc4 = 0
	(*satrec).cc5 = 0
	(*satrec).d2 = 0
	(*satrec).d3 = 0
	(*satrec).d4 = 0
	(*satrec).delmo = 0
	(*satrec).eta = 0
	(*satrec).argpdot = 0
	(*satrec).omgcof = 0
	(*satrec).sinmao = 0
	(*satrec).t = 0
	(*satrec).t2cof = 0
	(*satrec).t3cof = 0
	(*satrec).t4cof = 0
	(*satrec).t5cof = 0
	(*satrec).x1mth2 = 0
	(*satrec).x7thm1 = 0
	(*satrec).mdot = 0
	(*satrec).nodedot = 0
	(*satrec).xlcof = 0
	(*satrec).xmcof = 0
	(*satrec).nodecf = 0
	(*satrec).irez = int64(0)
	(*satrec).d2201 = 0
	(*satrec).d2211 = 0
	(*satrec).d3210 = 0
	(*satrec).d3222 = 0
	(*satrec).d4410 = 0
	(*satrec).d4422 = 0
	(*satrec).d5220 = 0
	(*satrec).d5232 = 0
	(*satrec).d5421 = 0
	(*satrec).d5433 = 0
	(*satrec).dedt = 0
	(*satrec).del1 = 0
	(*satrec).del2 = 0
	(*satrec).del3 = 0
	(*satrec).didt = 0
	(*satrec).dmdt = 0
	(*satrec).dnodt = 0
	(*satrec).domdt = 0
	(*satrec).e3 = 0
	(*satrec).ee2 = 0
	(*satrec).peo = 0
	(*satrec).pgho = 0
	(*satrec).pho = 0
	(*satrec).pinco = 0
	(*satrec).plo = 0
	(*satrec).se2 = 0
	(*satrec).se3 = 0
	(*satrec).sgh2 = 0
	(*satrec).sgh3 = 0
	(*satrec).sgh4 = 0
	(*satrec).sh2 = 0
	(*satrec).sh3 = 0
	(*satrec).si2 = 0
	(*satrec).si3 = 0
	(*satrec).sl2 = 0
	(*satrec).sl3 = 0
	(*satrec).sl4 = 0
	(*satrec).gsto = 0
	(*satrec).xfact = 0
	(*satrec).xgh2 = 0
	(*satrec).xgh3 = 0
	(*satrec).xgh4 = 0
	(*satrec).xh2 = 0
	(*satrec).xh3 = 0
	(*satrec).xi2 = 0
	(*satrec).xi3 = 0
	(*satrec).xl2 = 0
	(*satrec).xl3 = 0
	(*satrec).xl4 = 0
	(*satrec).xlamo = 0
	(*satrec).zmol = 0
	(*satrec).zmos = 0
	(*satrec).atime = 0
	(*satrec).xli = 0
	(*satrec).xni = 0
	getgravconst((*satrec).whichconst, satrec)
	(*satrec).error = int64(0)
	(*satrec).operationmode = opsmode
	(*satrec).nm = 0
	(*satrec).mm = (*satrec).nm
	(*satrec).Om = (*satrec).mm
	(*satrec).im = (*satrec).Om
	(*satrec).em = (*satrec).im
	(*satrec).am = (*satrec).em
	ss = 78.0/(*satrec).radiusearthkm + 1
	qzms2ttemp = (120.0 - 78.0) / (*satrec).radiusearthkm
	qzms2t = qzms2ttemp * qzms2ttemp * qzms2ttemp * qzms2ttemp
	x2o3 = 2.0 / 3.0
	(*satrec).init_ = 'y'
	(*satrec).t = 0
	initl(epoch, satrec)

	(*satrec).a = math.Pow((*satrec).no_unkozai*(*satrec).tumin, (-2.0 / 3.0))

	(*satrec).alta = (*satrec).a*(1+(*satrec).ecco) - 1
	(*satrec).altp = (*satrec).a*(1-(*satrec).ecco) - 1
	(*satrec).error = int64(0)
	if (*satrec).omeosq >= 0 || (*satrec).no_unkozai >= 0 {
		(*satrec).isimp = int64(0)
		if (*satrec).rp < 220.0/(*satrec).radiusearthkm+1 {
			(*satrec).isimp = int64(1)
		}
		sfour = ss
		qzms24 = qzms2t
		perige = ((*satrec).rp - 1) * (*satrec).radiusearthkm
		if perige < 156 {
			sfour = perige - 78
			if perige < 98 {
				sfour = 20
			}
			qzms24temp = (120.0 - sfour) / (*satrec).radiusearthkm
			qzms24 = qzms24temp * qzms24temp * qzms24temp * qzms24temp
			sfour = sfour/(*satrec).radiusearthkm + 1
		}
		pinvsq = 1.0 / (*satrec).posq
		tsi = 1.0 / ((*satrec).ao - sfour)
		(*satrec).eta = (*satrec).ao * (*satrec).ecco * tsi
		etasq = (*satrec).eta * (*satrec).eta
		eeta = (*satrec).ecco * (*satrec).eta
		psisq = math.Abs(1 - etasq)
		coef = qzms24 * math.Pow(tsi, 4)
		coef1 = coef / math.Pow(psisq, 3.5)
		cc2 = coef1 * (*satrec).no_unkozai * ((*satrec).ao*(1+1.5*etasq+eeta*(4+etasq)) + 0.375*(*satrec).j2*tsi/psisq*(*satrec).con41*(8+3*etasq*(8+etasq)))
		(*satrec).cc1 = (*satrec).bstar * cc2
		cc3 = 0
		if (*satrec).ecco > 0.0001 {
			cc3 = -2 * coef * tsi * (*satrec).j3oj2 * (*satrec).no_unkozai * (*satrec).sinio / (*satrec).ecco
		}
		(*satrec).x1mth2 = 1 - (*satrec).cosio2
		(*satrec).cc4 = 2 * (*satrec).no_unkozai * coef1 * (*satrec).ao * (*satrec).omeosq * ((*satrec).eta*(2+0.5*etasq) + (*satrec).ecco*(0.5+2*etasq) - (*satrec).j2*tsi/((*satrec).ao*psisq)*(-3*(*satrec).con41*(1-2*eeta+etasq*(1.5-0.5*eeta))+0.75*(*satrec).x1mth2*(2*etasq-eeta*(1+etasq))*math.Cos(2*(*satrec).argpo)))
		(*satrec).cc5 = 2 * coef1 * (*satrec).ao * (*satrec).omeosq * (1 + 2.75*(etasq+eeta) + eeta*etasq)
		cosio4 = (*satrec).cosio2 * (*satrec).cosio2
		temp1 = 1.5 * (*satrec).j2 * pinvsq * (*satrec).no_unkozai
		temp2 = 0.5 * temp1 * (*satrec).j2 * pinvsq
		temp3 = -0.46875 * (*satrec).j4 * pinvsq * pinvsq * (*satrec).no_unkozai
		(*satrec).mdot = (*satrec).no_unkozai + 0.5*temp1*(*satrec).rteosq*(*satrec).con41 + 0.0625*temp2*(*satrec).rteosq*(13-78*(*satrec).cosio2+137*cosio4)
		(*satrec).argpdot = -0.5*temp1*(*satrec).con42 + 0.0625*temp2*(7-114*(*satrec).cosio2+395*cosio4) + temp3*(3-36*(*satrec).cosio2+49*cosio4)
		xhdot1 = -temp1 * (*satrec).cosio
		(*satrec).nodedot = xhdot1 + (0.5*temp2*(4-19*(*satrec).cosio2)+2*temp3*(3-7*(*satrec).cosio2))*(*satrec).cosio
		xpidot = (*satrec).argpdot + (*satrec).nodedot
		(*satrec).omgcof = (*satrec).bstar * cc3 * math.Cos((*satrec).argpo)
		(*satrec).xmcof = 0
		if (*satrec).ecco > 0.0001 {
			(*satrec).xmcof = -x2o3 * coef * (*satrec).bstar / eeta
		}
		(*satrec).nodecf = 3.5 * (*satrec).omeosq * xhdot1 * (*satrec).cc1
		(*satrec).t2cof = 1.5 * (*satrec).cc1
		if math.Abs((*satrec).cosio+1) > 1.5e-12 {
			(*satrec).xlcof = -0.25 * (*satrec).j3oj2 * (*satrec).sinio * (3.0 + 5.0*(*satrec).cosio) / (1.0 + (*satrec).cosio)
		} else {
			(*satrec).xlcof = -0.25 * (*satrec).j3oj2 * (*satrec).sinio * (3.0 + 5.0*(*satrec).cosio) / temp4
		}
		(*satrec).aycof = -0.5 * (*satrec).j3oj2 * (*satrec).sinio
		delmotemp = 1 + (*satrec).eta*math.Cos((*satrec).mo)
		(*satrec).delmo = delmotemp * delmotemp * delmotemp
		(*satrec).sinmao = math.Sin((*satrec).mo)
		(*satrec).x7thm1 = 7*(*satrec).cosio2 - 1
		if float64(int64(2))*3.141592653589793/(*satrec).no_unkozai >= 225 {
			(*satrec).method = 'd'
			(*satrec).isimp = int64(1)
			tc = 0
			(*satrec).inclm = (*satrec).inclo
			dscom(epoch, (*satrec).ecco, (*satrec).argpo, tc, (*satrec).inclo, (*satrec).nodeo, (*satrec).no_unkozai, satrec)
			(*satrec).ep = (*satrec).ecco
			(*satrec).inclp = (*satrec).inclo
			(*satrec).nodep = (*satrec).nodeo
			(*satrec).argpp = (*satrec).argpo
			(*satrec).mp = (*satrec).mo
			dpper((*satrec).e3, (*satrec).ee2, (*satrec).peo, (*satrec).pgho, (*satrec).pho, (*satrec).pinco, (*satrec).plo, (*satrec).se2, (*satrec).se3, (*satrec).sgh2, (*satrec).sgh3, (*satrec).sgh4, (*satrec).sh2, (*satrec).sh3, (*satrec).si2, (*satrec).si3, (*satrec).sl2, (*satrec).sl3, (*satrec).sl4, (*satrec).t, (*satrec).xgh2, (*satrec).xgh3, (*satrec).xgh4, (*satrec).xh2, (*satrec).xh3, (*satrec).xi2, (*satrec).xi3, (*satrec).xl2, (*satrec).xl3, (*satrec).xl4, (*satrec).zmol, (*satrec).zmos, (*satrec).init_, satrec, (*satrec).operationmode)
			(*satrec).ecco = (*satrec).ep
			(*satrec).inclo = (*satrec).inclp
			(*satrec).nodeo = (*satrec).nodep
			(*satrec).argpo = (*satrec).argpp
			(*satrec).mo = (*satrec).mp
			(*satrec).argpm = 0
			(*satrec).nodem = 0
			(*satrec).mm = 0
			dsinit(tc, xpidot, satrec)
		}
		if (*satrec).isimp != int64(1) {
			cc1sq = (*satrec).cc1 * (*satrec).cc1
			(*satrec).d2 = 4 * (*satrec).ao * tsi * cc1sq
			temp = (*satrec).d2 * tsi * (*satrec).cc1 / 3.0
			(*satrec).d3 = (17*(*satrec).ao + sfour) * temp
			(*satrec).d4 = 0.5 * temp * (*satrec).ao * tsi * (221*(*satrec).ao + 31*sfour) * (*satrec).cc1
			(*satrec).t3cof = (*satrec).d2 + 2*cc1sq
			(*satrec).t4cof = 0.25 * (3*(*satrec).d3 + (*satrec).cc1*(12*(*satrec).d2+10*cc1sq))
			(*satrec).t5cof = 0.2 * (3*(*satrec).d4 + 12*(*satrec).cc1*(*satrec).d3 + 6*(*satrec).d2*(*satrec).d2 + 15*cc1sq*(2*(*satrec).d2+cc1sq))
		}
	}
	sgp4(satrec, 0, &r[0], &v[0])

	(*satrec).init_ = 'n'
	return bool_((int64(1)))
}

// sgp4 - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:1755
/*-----------------------------------------------------------------------------
 *
 *                             procedure sgp4
 *
 *  this procedure is the sgp4 prediction model from space command. this is an
 *    updated and combined version of sgp4 and sdp4, which were originally
 *    published separately in spacetrack report #3. this version follows the
 *    methodology from the aiaa paper (2006) describing the history and
 *    development of the code.
 *
 *  author        : david vallado                  719-573-2600   28 jun 2005
 *
 *  inputs        :
 *    satrec     - initialised structure from sgp4init() call.
 *    tsince     - time since epoch (minutes)
 *
 *  outputs       :
 *    r           - position vector                     km
 *    v           - velocity                            km/sec
 *  return code - non-zero on error.
 *                   1 - mean elements, ecc >= 1.0 or ecc < -0.001 or a < 0.95 er
 *                   2 - mean motion less than 0.0
 *                   3 - pert elements, ecc < 0.0  or  ecc > 1.0
 *                   4 - semi-latus rectum < 0.0
 *                   5 - epoch elements are sub-orbital
 *                   6 - satellite has decayed
 *
 *  locals        :
 *    am          -
 *    axnl, aynl        -
 *    betal       -
 *    cosim   , sinim   , cosomm  , sinomm  , cnod    , snod    , cos2u   ,
 *    sin2u   , coseo1  , sineo1  , cosi    , sini    , cosip   , sinip   ,
 *    cosisq  , cossu   , sinsu   , cosu    , sinu
 *    delm        -
 *    delomg      -
 *    dndt        -
 *    eccm        -
 *    emsq        -
 *    ecose       -
 *    el2         -
 *    eo1         -
 *    eccp        -
 *    esine       -
 *    argpm       -
 *    argpp       -
 *    omgadf      -c
 *    pl          -
 *    r           -
 *    rtemsq      -
 *    rdotl       -
 *    rl          -
 *    rvdot       -
 *    rvdotl      -
 *    su          -
 *    t2  , t3   , t4    , tc
 *    tem5, temp , temp1 , temp2  , tempa  , tempe  , templ
 *    u   , ux   , uy    , uz     , vx     , vy     , vz
 *    inclm       - inclination
 *    mm          - mean anomaly
 *    nm          - mean motion
 *    nodem       - right asc of ascending node
 *    xinc        -
 *    xincp       -
 *    xl          -
 *    xlm         -
 *    mp          -
 *    xmdf        -
 *    xmx         -
 *    xmy         -
 *    nodedf      -
 *    xnode       -
 *    nodep       -
 *    np          -
 *
 *  coupling      :
 *    getgravconst- no longer used. Variables are conatined within satrec
 *    dpper
 *    dpspace
 *
 *  references    :
 *    hoots, roehrich, norad spacetrack report #3 1980
 *    hoots, norad spacetrack report #6 1986
 *    hoots, schumacher and glover 2004
 *    vallado, crawford, hujsak, kelso  2006
 ----------------------------------------------------------------------------*/ //
/* ------------------ set mathematical constants --------------- */ //
// sgp4fix divisor for divide by zero check on inclination
// the old check used 1.0 + cos(pi-1.0e-9), but then compared it to
// 1.5 e-12, so the threshold was changed to 1.5e-12 for consistency
// sgp4fix identify constants and allow alternate values
// getgravconst( whichconst, tumin, mu, radiusearthkm, xke, j2, j3, j4, j3oj2 );
/* --------------------- clear sgp4 error flag ----------------- */ //
/* ------- update for secular gravity and atmospheric drag ----- */ //
// sgp4fix use mutliply for speed instead of pow
// if method = d
// sgp4fix add return
// fix tolerance for error recognition
// sgp4fix am is fixed from the previous nm check
/* || (am < 0.95)*/ //
// sgp4fix to return if there is an error in eccentricity
// sgp4fix fix tolerance to avoid a divide by zero
// sgp4fix recover singly averaged mean elements
/* ----------------- compute extra mean quantities ------------- */ //
/* -------------------- add lunar-solar periodics -------------- */ //
// sgp4fix add return
// if method = d
/* -------------------- long period periodics ------------------ */ //
// sgp4fix for divide by zero for xincp = 180 deg
/* --------------------- solve kepler's equation --------------- */ //
//   sgp4fix for kepler iteration
//   the following iteration needs better limits on corrections
/* ------------- short period preliminary quantities ----------- */ //
// sgp4fix add return
/* -------------- update for short period periodics ------------ */ //
/* --------------------- orientation vectors ------------------- */ //
/* --------- position and velocity (in km and km/sec) ---------- */ //
// if pl > 0
// sgp4fix for decaying satellites
// sgp4
//
func sgp4(satrec *ElsetRec, tsince float64, r *float64, v *float64) bool_ {
	var axnl float64
	var aynl float64
	var betal float64
	var cnod float64
	var cos2u float64
	var coseo1 float64
	var cosi float64
	var cosip float64
	var cosisq float64
	var cossu float64
	var cosu float64
	var delm float64
	var delomg float64
	var ecose float64
	var el2 float64
	var eo1 float64
	var esine float64
	var argpdf float64
	var pl float64
	var mrt float64 = 0
	var mvt float64
	var rdotl float64
	var rl float64
	var rvdot float64
	var rvdotl float64
	var sin2u float64
	var sineo1 float64
	var sini float64
	var sinip float64
	var sinsu float64
	var sinu float64
	var snod float64
	var su float64
	var t2 float64
	var t3 float64
	var t4 float64
	var tem5 float64
	var temp float64
	var temp1 float64
	var temp2 float64
	var tempa float64
	var tempe float64
	var templ float64
	var u float64
	var ux float64
	var uy float64
	var uz float64
	var vx float64
	var vy float64
	var vz float64
	var xinc float64
	var xincp float64
	var xl float64
	var xlm float64
	var xmdf float64
	var xmx float64
	var xmy float64
	var nodedf float64
	var xnode float64
	var tc float64
	var x2o3 float64
	var vkmpersec float64
	var delmtemp float64
	var ktr int64
	var temp4 float64 = 1.5e-12
	x2o3 = 2.0 / 3.0
	vkmpersec = (*satrec).radiusearthkm * (*satrec).xke / 60.0
	(*satrec).t = tsince
	(*satrec).error = int64(0)
	xmdf = (*satrec).mo + (*satrec).mdot*(*satrec).t
	argpdf = (*satrec).argpo + (*satrec).argpdot*(*satrec).t
	nodedf = (*satrec).nodeo + (*satrec).nodedot*(*satrec).t
	(*satrec).argpm = argpdf
	(*satrec).mm = xmdf
	t2 = (*satrec).t * (*satrec).t
	(*satrec).nodem = nodedf + (*satrec).nodecf*t2
	tempa = 1 - (*satrec).cc1*(*satrec).t
	tempe = (*satrec).bstar * (*satrec).cc4 * (*satrec).t
	templ = (*satrec).t2cof * t2
	delomg = float64(int64(0))
	delmtemp = float64(int64(0))
	delm = float64(int64(0))
	temp = float64(int64(0))
	t3 = float64(int64(0))
	t4 = float64(int64(0))
	mrt = float64(int64(0))
	if (*satrec).isimp != int64(1) {
		delomg = (*satrec).omgcof * (*satrec).t
		delmtemp = 1 + (*satrec).eta*math.Cos(xmdf)
		delm = (*satrec).xmcof * (delmtemp*delmtemp*delmtemp - (*satrec).delmo)
		temp = delomg + delm
		(*satrec).mm = xmdf + temp
		(*satrec).argpm = argpdf - temp
		t3 = t2 * (*satrec).t
		t4 = t3 * (*satrec).t
		tempa = tempa - (*satrec).d2*t2 - (*satrec).d3*t3 - (*satrec).d4*t4
		tempe = tempe + (*satrec).bstar*(*satrec).cc5*(math.Sin((*satrec).mm)-(*satrec).sinmao)
		templ = templ + (*satrec).t3cof*t3 + t4*((*satrec).t4cof+(*satrec).t*(*satrec).t5cof)
	}
	tc = float64(int64(0))
	(*satrec).nm = (*satrec).no_unkozai
	(*satrec).em = (*satrec).ecco
	(*satrec).inclm = (*satrec).inclo
	if int64((*satrec).method) == int64('d') {
		tc = (*satrec).t
		dspace(tc, satrec)
	}
	if (*satrec).nm <= 0 {
		(*satrec).error = int64(2)
		return bool_((int64(0)))
	}
	(*satrec).am = math.Pow(((*satrec).xke/(*satrec).nm), x2o3) * tempa * tempa
	(*satrec).nm = (*satrec).xke / math.Pow((*satrec).am, 1.5)
	(*satrec).em = (*satrec).em - tempe
	if (*satrec).em >= 1 || (*satrec).em < -0.001 {
		(*satrec).error = int64(1)
		return bool_((int64(0)))
	}
	if (*satrec).em < 1e-06 {
		(*satrec).em = 1e-06
	}
	(*satrec).mm = (*satrec).mm + (*satrec).no_unkozai*templ
	xlm = (*satrec).mm + (*satrec).argpm + (*satrec).nodem
	(*satrec).emsq = (*satrec).em * (*satrec).em
	temp = 1 - (*satrec).emsq
	(*satrec).nodem = math.Mod((*satrec).nodem, (2 * 3.141592653589793))
	(*satrec).argpm = math.Mod((*satrec).argpm, (2 * 3.141592653589793))
	xlm = math.Mod(xlm, (2 * 3.141592653589793))
	(*satrec).mm = math.Mod(xlm-(*satrec).argpm-(*satrec).nodem, (2 * 3.141592653589793))
	(*satrec).am = (*satrec).am
	(*satrec).em = (*satrec).em
	(*satrec).im = (*satrec).inclm
	(*satrec).Om = (*satrec).nodem
	(*satrec).om = (*satrec).argpm
	(*satrec).mm = (*satrec).mm
	(*satrec).nm = (*satrec).nm
	(*satrec).sinim = math.Sin((*satrec).inclm)
	(*satrec).cosim = math.Cos((*satrec).inclm)
	(*satrec).ep = (*satrec).em
	xincp = (*satrec).inclm
	(*satrec).inclp = (*satrec).inclm
	(*satrec).argpp = (*satrec).argpm
	(*satrec).nodep = (*satrec).nodem
	(*satrec).mp = (*satrec).mm
	sinip = (*satrec).sinim
	cosip = (*satrec).cosim
	if int64((*satrec).method) == int64('d') {
		dpper((*satrec).e3, (*satrec).ee2, (*satrec).peo, (*satrec).pgho, (*satrec).pho, (*satrec).pinco, (*satrec).plo, (*satrec).se2, (*satrec).se3, (*satrec).sgh2, (*satrec).sgh3, (*satrec).sgh4, (*satrec).sh2, (*satrec).sh3, (*satrec).si2, (*satrec).si3, (*satrec).sl2, (*satrec).sl3, (*satrec).sl4, (*satrec).t, (*satrec).xgh2, (*satrec).xgh3, (*satrec).xgh4, (*satrec).xh2, (*satrec).xh3, (*satrec).xi2, (*satrec).xi3, (*satrec).xl2, (*satrec).xl3, (*satrec).xl4, (*satrec).zmol, (*satrec).zmos, 'n', satrec, (*satrec).operationmode)
		xincp = (*satrec).inclp
		if xincp < 0 {
			xincp = -xincp
			(*satrec).nodep = (*satrec).nodep + 3.141592653589793
			(*satrec).argpp = (*satrec).argpp - 3.141592653589793
		}
		if (*satrec).ep < 0 || (*satrec).ep > 1 {
			(*satrec).error = int64(3)
			return bool_((int64(0)))
		}
	}
	if int64((*satrec).method) == int64('d') {
		sinip = math.Sin(xincp)
		cosip = math.Cos(xincp)
		(*satrec).aycof = -0.5 * (*satrec).j3oj2 * sinip
		if math.Abs(cosip+1) > 1.5e-12 {
			(*satrec).xlcof = -0.25 * (*satrec).j3oj2 * sinip * (3.0 + 5.0*cosip) / (1.0 + cosip)
		} else {
			(*satrec).xlcof = -0.25 * (*satrec).j3oj2 * sinip * (3.0 + 5.0*cosip) / temp4
		}
	}
	axnl = (*satrec).ep * math.Cos((*satrec).argpp)
	temp = 1.0 / ((*satrec).am * (1 - (*satrec).ep*(*satrec).ep))
	aynl = (*satrec).ep*math.Sin((*satrec).argpp) + temp*(*satrec).aycof
	xl = (*satrec).mp + (*satrec).argpp + (*satrec).nodep + temp*(*satrec).xlcof*axnl
	u = math.Mod(xl-(*satrec).nodep, (2 * 3.141592653589793))
	eo1 = u
	tem5 = 9999.9
	ktr = int64(1)
	sineo1 = float64(int64(0))
	coseo1 = float64(int64(0))
	for math.Abs(tem5) >= 1e-12 && ktr <= int64(10) {
		sineo1 = math.Sin(eo1)
		coseo1 = math.Cos(eo1)
		tem5 = 1 - coseo1*axnl - sineo1*aynl
		tem5 = (u - aynl*coseo1 + axnl*sineo1 - eo1) / tem5
		if math.Abs(tem5) >= 0.95 {
			tem5 = func() float64 {
				if tem5 > 0 {
					return 0.95
				} else {
					return -0.95
				}
			}()
		}
		eo1 = eo1 + tem5
		ktr = ktr + int64(1)
	}
	ecose = axnl*coseo1 + aynl*sineo1
	esine = axnl*sineo1 - aynl*coseo1
	el2 = axnl*axnl + aynl*aynl
	pl = (*satrec).am * (1 - el2)
	if pl < 0 {
		(*satrec).error = int64(4)
		return bool_((int64(0)))
	} else {
		rl = (*satrec).am * (1 - ecose)
		rdotl = math.Sqrt((*satrec).am) * esine / rl
		rvdotl = math.Sqrt(pl) / rl
		betal = math.Sqrt(1 - el2)
		temp = esine / (1 + betal)
		sinu = (*satrec).am / rl * (sineo1 - aynl - axnl*temp)
		cosu = (*satrec).am / rl * (coseo1 - axnl + aynl*temp)
		su = math.Atan2(sinu, cosu)
		sin2u = (cosu + cosu) * sinu
		cos2u = 1 - 2*sinu*sinu
		temp = 1.0 / pl
		temp1 = 0.5 * (*satrec).j2 * temp
		temp2 = temp1 * temp
		if int64((*satrec).method) == int64('d') {
			cosisq = cosip * cosip
			(*satrec).con41 = 3*cosisq - 1
			(*satrec).x1mth2 = 1 - cosisq
			(*satrec).x7thm1 = 7*cosisq - 1
		}
		mrt = rl*(1-1.5*temp2*betal*(*satrec).con41) + 0.5*temp1*(*satrec).x1mth2*cos2u
		su = su - 0.25*temp2*(*satrec).x7thm1*sin2u
		xnode = (*satrec).nodep + 1.5*temp2*cosip*sin2u
		xinc = xincp + 1.5*temp2*cosip*sinip*cos2u
		mvt = rdotl - (*satrec).nm*temp1*(*satrec).x1mth2*sin2u/(*satrec).xke
		rvdot = rvdotl + (*satrec).nm*temp1*((*satrec).x1mth2*cos2u+1.5*(*satrec).con41)/(*satrec).xke
		sinsu = math.Sin(su)
		cossu = math.Cos(su)
		snod = math.Sin(xnode)
		cnod = math.Cos(xnode)
		sini = math.Sin(xinc)
		cosi = math.Cos(xinc)
		xmx = -snod * cosi
		xmy = cnod * cosi
		ux = xmx*sinsu + cnod*cossu
		uy = xmy*sinsu + snod*cossu
		uz = sini * sinsu
		vx = xmx*cossu - cnod*sinsu
		vy = xmy*cossu - snod*sinsu
		vz = sini * cossu
		*r = mrt * ux * (*satrec).radiusearthkm
		*((*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + (uintptr)(int64(1))*unsafe.Sizeof(*r)))) = mrt * uy * (*satrec).radiusearthkm
		*((*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(r)) + (uintptr)(int64(2))*unsafe.Sizeof(*r)))) = mrt * uz * (*satrec).radiusearthkm
		*v = (mvt*ux + rvdot*vx) * vkmpersec
		*((*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + (uintptr)(int64(1))*unsafe.Sizeof(*v)))) = (mvt*uy + rvdot*vy) * vkmpersec
		*((*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + (uintptr)(int64(2))*unsafe.Sizeof(*v)))) = (mvt*uz + rvdot*vz) * vkmpersec
	}
	if mrt < 1 {
		(*satrec).error = int64(6)
		return bool_((int64(0)))
	}
	return bool_((int64(1)))
}

// getgravconst - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:2071
/* -----------------------------------------------------------------------------
 *
 *                           function getgravconst
 *
 *  this function gets constants for the propagator. note that mu is identified to
 *    facilitiate comparisons with newer models. the common useage is wgs72.
 *
 *  author        : david vallado                  719-573-2600   21 jul 2006
 *
 *  inputs        :
 *    whichconst  - which set of constants to use  wgs72old, wgs72, wgs84
 *
 *  outputs       :
 *    tumin       - minutes in one time unit
 *    mu          - earth gravitational parameter
 *    radiusearthkm - radius of the earth in km
 *    xke         - reciprocal of tumin
 *    j2, j3, j4  - un-normalized zonal harmonic values
 *    j3oj2       - j3 divided by j2
 *
 *  locals        :
 *
 *  coupling      :
 *    none
 *
 *  references    :
 *    norad spacetrack report #3
 *    vallado, crawford, hujsak, kelso  2006
 --------------------------------------------------------------------------- */ //
// -- wgs-72 low precision str#3 constants --
// in km3 / s2
// km
// reciprocal of tumin
// ------------ wgs-72 constants ------------
// in km3 / s2
// km
// ------------ wgs-84 constants ------------
// in km3 / s2
// km
// getgravconst
//
func getgravconst(whichconst int64, rec *ElsetRec) {
	(*rec).whichconst = whichconst
	switch whichconst {
	case int64(1):
		{
			(*rec).mu = 398600.79964
			(*rec).radiusearthkm = 6378.135
			(*rec).xke = 0.0743669161
			(*rec).tumin = 1.0 / (*rec).xke
			(*rec).j2 = 0.001082616
			(*rec).j3 = -2.53881e-06
			(*rec).j4 = -1.65597e-06
			(*rec).j3oj2 = (*rec).j3 / (*rec).j2
		}
	case int64(2):
		{
			(*rec).mu = 398600.8
			(*rec).radiusearthkm = 6378.135
			(*rec).xke = 60.0 / math.Sqrt((*rec).radiusearthkm*(*rec).radiusearthkm*(*rec).radiusearthkm/(*rec).mu)
			(*rec).tumin = 1.0 / (*rec).xke
			(*rec).j2 = 0.001082616
			(*rec).j3 = -2.53881e-06
			(*rec).j4 = -1.65597e-06
			(*rec).j3oj2 = (*rec).j3 / (*rec).j2
		}
	default:
		fallthrough
	case int64(3):
		{
			(*rec).mu = 398600.5
			(*rec).radiusearthkm = 6378.137
			(*rec).xke = 60 / math.Sqrt((*rec).radiusearthkm*(*rec).radiusearthkm*(*rec).radiusearthkm/(*rec).mu)
			(*rec).tumin = 1.0 / (*rec).xke
			(*rec).j2 = 0.00108262998905
			(*rec).j3 = -2.53215306e-06
			(*rec).j4 = -1.61098761e-06
			(*rec).j3oj2 = (*rec).j3 / (*rec).j2
			break
		}
	}
}

// gstime - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:2147
/* -----------------------------------------------------------------------------
 *
 *                           function gstime
 *
 *  this function finds the greenwich sidereal time.
 *
 *  author        : david vallado                  719-573-2600    1 mar 2001
 *
 *  inputs          description                    range / units
 *    jdut1       - julian date in ut1             days from 4713 bc
 *
 *  outputs       :
 *    gstime      - greenwich sidereal time        0 to 2pi rad
 *
 *  locals        :
 *    temp        - temporary variable for doubles   rad
 *    tut1        - julian centuries from the
 *                  jan 1, 2000 12 h epoch (ut1)
 *
 *  coupling      :
 *    none
 *
 *  references    :
 *    vallado       2013, 187, eq 3-45
 * --------------------------------------------------------------------------- */ //
// sec
//360/86400 = 1/240, to deg, to rad
// ------------------------ check quadrants ---------------------
// gstime
//
func gstime(jdut1 float64) float64 {
	var temp float64
	var tut1 float64
	tut1 = (jdut1 - 2.451545e+06) / 36525.0
	temp = -6.2e-06*tut1*tut1*tut1 + 0.093104*tut1*tut1 + (876600*float64(int64(3600))+8.640184812866e+06)*tut1 + 67310.54841
	temp = math.Mod(temp*(3.141592653589793/180)/240, (2 * 3.141592653589793))
	if temp < 0 {
		temp += 2 * 3.141592653589793
	}
	return temp
}

// jday - transpiled function from  /home/somebody/aholinch/sgp4/src/c/all.c:2196
/* -----------------------------------------------------------------------------
 *
 *                           procedure jday
 *
 *  this procedure finds the julian date given the year, month, day, and time.
 *    the julian date is defined by each elapsed day since noon, jan 1, 4713 bc.
 *
 *  algorithm     : calculate the answer in one step for efficiency
 *
 *  author        : david vallado                  719-573-2600    1 mar 2001
 *
 *  inputs          description                    range / units
 *    year        - year                           1900 .. 2100
 *    mon         - month                          1 .. 12
 *    day         - day                            1 .. 28,29,30,31
 *    hr          - universal time hour            0 .. 23
 *    min         - universal time min             0 .. 59
 *    sec         - universal time sec             0.0 .. 59.999
 *
 *  outputs       :
 *    jd          - julian date                    days from 4713 bc
 *    jdfrac      - julian date fraction into day  days from 4713 bc
 *
 *  locals        :
 *    none.
 *
 *  coupling      :
 *    none.
 *
 *  references    :
 *    vallado       2013, 183, alg 14, ex 3-4
 * --------------------------------------------------------------------------- */ //
// use - 678987.0 to go to mjd directly
// check that the day and fractional day are correct
// jday
//
func jday(year int64, mon int64, day int64, hr int64, minute int64, sec float64, jd *float64, jdfrac *float64) {
	*jd = 367*float64(year) - math.Floor(float64(int64(7))*(float64(year)+math.Floor(float64(mon+int64(9))/12))*0.25) + math.Floor(float64(int64(275)*mon)/9) + float64(day) + 1.7210135e+06
	*jdfrac = (sec + float64(minute)*60 + float64(hr)*3600) / 86400
	if math.Abs(*jdfrac) > 1 {
		var dtt float64 = math.Floor(*jdfrac)
		*jd = *jd + dtt
		*jdfrac = *jdfrac - dtt
	}
	return
}

type Error int

const decayError = "satellite has decayed"

func HasDecayed(e error) bool {
	// Too much trouble to try to check type (in the face of
	// fmt.wrapError, etc).
	return strings.Contains(e.Error(), decayError)
}

func (e Error) Error() string {
	var msg string
	switch int(e) {
	case 1:
		msg = "mean elements, ecc >= 1.0 or ecc < -0.001 or a < 0.95 er"
	case 2:
		msg = "mean motion less than 0.0"
	case 3:
		msg = "pert elements, ecc < 0.0  or  ecc > 1.0"
	case 4:
		msg = "semi-latus rectum < 0.0"
	case 5:
		msg = "epoch elements are sub-orbital"
	case 6:
		// See HashDecayed.
		msg = decayError
	default:
		msg = "NA"
	}
	return fmt.Sprintf("code=%d: %s", e, msg)
}

func (tle *TLE) PropUnixMillis(ms int64) ([]float64, []float64, error) {
	var (
		r = make([]float64, 3)
		v = make([]float64, 3)
	)

	tle.Lock()
	tle.Rec.error = 0
	getRVForDate(tle, ms, (*float64)(&r[0]), (*float64)(&v[0]))
	e := tle.sgp4Error
	tle.sgp4Error = 0
	tle.Unlock()

	if e != 0 {
		return nil, nil, fmt.Errorf("SGP4 error at ms=%d: %w", ms, Error(e))
	}
	return r, v, nil
}

func (tle *TLE) PropForMins(mins float64) ([]float64, []float64, error) {
	var (
		r = make([]float64, 3)
		v = make([]float64, 3)
	)

	tle.Lock()
	tle.Rec.error = 0
	getRV(tle, mins, (*float64)(&r[0]), (*float64)(&v[0]))
	e := tle.sgp4Error
	tle.Rec.error = 0
	tle.Unlock()

	if e != 0 {
		return nil, nil, fmt.Errorf("SGP4 error at mins=%f: %w", mins, Error(e))
	}
	return r, v, nil
}

func ParseLines(line1, line2 string) (*TLE, error) {
	tle := &TLE{}
	bs1 := []byte(line1)
	bs2 := []byte(line2)
	parseLines(tle, (*byte)(&bs1[0]), (*byte)(&bs2[0]))
	return tle, nil
}
