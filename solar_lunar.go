package calendar

import (
	"fmt"
	"math"
)

func jdFromDate(dd, mm, yy int) int {
	a := (14 - mm) / 12
	y := yy + 4800 - a
	m := mm + 12*a - 3

	jd := dd + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
	if jd < 2299161 {
		jd = dd + (153*m+2)/5 + 365*y + y/4 - 32083
	}
	return jd
}

func jdToDate(jd int) (dd, mm, yy int) {
	var a, b, c int
	if jd > 2299160 { // After 5/10/1582, Gregorian calendar
		a = jd + 32044
		b = int((4*a + 3) / 146097.)
		c = a - int((b*146097)/4.)
	} else {
		b = 0
		c = jd + 32082
	}
	d := int((4*c + 3) / 1461.)
	e := c - int((1461*d)/4.)
	m := int((5*e + 2) / 153.)
	dd = e - int((153*m+2)/5.) + 1
	mm = m + 3 - 12*int(m/10.)
	yy = b*100 + d - 4800 + int(m/10.)
	return dd, mm, yy
}

// NewMoon(k int): Compute the time of the k-th new moon after the new moon of 1/1/1900 13:52 UCT
// (measured as the number of days since 1/1/4713 BC noon UCT, e.g., 2451545.125 is 1/1/2000 15:00 UTC.
// Returns a floating number, e.g., 2415079.9758617813 for k=2 or 2414961.935157746 for k=-2.
// Time in Julian centuries from 1900 January 0.5
func NewMoon(k int) float64 {
	kf := float64(k)

	T := kf / 1236.85
	T2 := T * T
	T3 := T2 * T
	dr := math.Pi / 180.
	Jd1 := 2415020.75933 + 29.53058868*kf + 0.0001178*T2 - 0.000000155*T3
	Jd1 = Jd1 + 0.00033*math.Sin((166.56+132.87*T-0.009173*T2)*dr)

	//   Mean new moon
	M := 359.2242 + 29.10535608*kf - 0.0000333*T2 - 0.00000347*T3

	//   Sun's mean anomaly
	Mpr := 306.0253 + 385.81691806*kf + 0.0107306*T2 + 0.00001236*T3

	// Moon's mean anomaly
	F := 21.2964 + 390.67050646*kf - 0.0016528*T2 - 0.00000239*T3

	//   Moon's argument of latitude
	var C1 float64
	C1 = (0.1734-0.000393*T)*math.Sin(M*dr) + 0.0021*math.Sin(2*dr*M)
	C1 = C1 - 0.4068*math.Sin(Mpr*dr) + 0.0161*math.Sin(dr*2*Mpr)
	C1 = C1 - 0.0004*math.Sin(dr*3*Mpr)
	C1 = C1 + 0.0104*math.Sin(dr*2*F) - 0.0051*math.Sin(dr*(M+Mpr))
	C1 = C1 - 0.0074*math.Sin(dr*(M-Mpr)) + 0.0004*math.Sin(dr*(2*F+M))
	C1 = C1 - 0.0004*math.Sin(dr*(2*F-M)) - 0.0006*math.Sin(dr*(2*F+Mpr))
	C1 = C1 + 0.0010*math.Sin(dr*(2*F-Mpr)) + 0.0005*math.Sin(dr*(2*Mpr+M))

	var deltat float64
	if T < -11 {
		deltat = 0.001 + 0.000839*T + 0.0002261*T2 - 0.00000845*T3 - 0.000000081*T*T3
	} else {
		deltat = -0.000278 + 0.000265*T + 0.000262*T2
	}

	JdNew := Jd1 + C1 - deltat
	return JdNew
}

// Compute the day of the k-th new moon in the given time zone.
// The time zone if the time difference between local time and UTC: 7.0 for UTC+7:00.
func GetNewMoonDay(k, timezone int) int {
	return int(NewMoon(k) + 0.5 + (float64(timezone) / 24.))
}

// Compute the longitude of the sun at any time.
// Parameter: floating number jdn, the number of days since 1/1/4713 BC noon.
func SunLongitude(jdn float64) float64 {
	T := (jdn - 2451545.0) / 36525.

	// Time in Julian centuries
	// from 2000-01-01 12:00:00 GMT
	T2 := T * T
	dr := math.Pi / 180. // degree to radian
	M := 357.52910 + 35999.05030*T - 0.0001559*T2 - 0.00000048*T*T2

	// mean anomaly, degree
	L0 := 280.46645 + 36000.76983*T + 0.0003032*T2
	// mean longitude, degree
	DL := (1.914600 - 0.004817*T - 0.000014*T2) * math.Sin(dr*M)
	DL += (0.019993-0.000101*T)*math.Sin(dr*2*M) + 0.000290*math.Sin(dr*3*M)
	L := L0 + DL // true longitude, degree
	L = L * dr
	L = L - math.Pi*2*float64(int(L/(math.Pi*2)))
	//  Normalize to (0, 2*math.pi)
	return L
}

// Compute sun position at midnight of the day with the given Julian day number.
// The time zone if the time difference between local time and UTC: 7.0 for UTC+7:00.
// The function returns a number between 0 and 11.
// From the day after March equinox and the 1st major term after March equinox, 0 is returned.
// After that, return 1, 2, 3 ...
func GetSunLongitude(dayNumber, timezone int) int {
	return int(SunLongitude(float64(dayNumber)-0.5-float64(timezone)/24.) / math.Pi * 6)
}

// Find the day that starts the luner month 11 of the given year for the given time zone.
func GetLunarMonth11(yy, timezone int) int {
	off := float64(jdFromDate(31, 12, yy)) - 2415021.076998695
	k := int(off / 29.530588853)
	nm := GetNewMoonDay(k, timezone)
	sunLong := GetSunLongitude(nm, timezone)
	fmt.Println(nm)
	// sun longitude at local midnight
	if sunLong >= 9 {
		nm = GetNewMoonDay(k-1, timezone)
	}

	return nm
}

// Find the index of the leap month after the month starting on the day a11.
func GetLeapMonthOffset(a11, timezone int) int {
	k := int((float64(a11)-2415021.076998695)/29.530588853 + 0.5)
	last := 0
	i := 1 // start with month following lunar month 11
	arc := GetSunLongitude(GetNewMoonDay(k+i, timezone), timezone)
	for {
		last = arc
		i += 1
		arc = GetSunLongitude(GetNewMoonDay(k+i, timezone), timezone)
		if !(arc != last && i < 14) {
			break
		}
	}
	return i - 1
}

// S2L Convert solar date dd/mm/yyyy to the corresponding lunar date.
func S2L(dd, mm, yy, timezone int) (lunarDay, lunarMonth, lunarYear, lunarLeap int) {
	dayNumber := jdFromDate(dd, mm, yy)
	k := int((float64(dayNumber) - 2415021.076998695) / 29.530588853)
	monthStart := GetNewMoonDay(k+1, timezone)
	if monthStart > dayNumber {
		monthStart = GetNewMoonDay(k, timezone)
	}
	a11 := GetLunarMonth11(yy, timezone)
	b11 := a11
	if a11 >= monthStart {
		lunarYear = yy
		a11 = GetLunarMonth11(yy-1, timezone)
	} else {
		lunarYear = yy + 1
		b11 = GetLunarMonth11(yy+1, timezone)
	}
	lunarDay = dayNumber - monthStart + 1
	diff := int(float64(monthStart-a11) / 29.)
	lunarLeap = 0
	lunarMonth = diff + 11
	if b11-a11 > 365 {
		leapMonthDiff := GetLeapMonthOffset(a11, timezone)
		if diff >= leapMonthDiff {
			lunarMonth = diff + 10
		}
		if diff == leapMonthDiff {
			lunarLeap = 1
		}
	}
	if lunarMonth > 12 {
		lunarMonth = lunarMonth - 12
	}
	if lunarMonth >= 11 && diff < 4 {
		lunarYear -= 1
	}
	return lunarDay, lunarMonth, lunarYear, lunarLeap
}

// Convert a lunar date to the corresponding solar date.
func L2S(lunarD, lunarM, lunarY, lunarLeap, timezone int) (dd, mm, yy int) {
	var a11, b11 int
	if lunarM < 11 {
		a11 = GetLunarMonth11(lunarY-1, timezone)
		b11 = GetLunarMonth11(lunarY, timezone)
	} else {
		a11 = GetLunarMonth11(lunarY, timezone)
		b11 = GetLunarMonth11(lunarY+1, timezone)
	}
	k := int(0.5 + (float64(a11)-2415021.076998695)/29.530588853)
	off := lunarM - 11
	if off < 0 {
		off += 12
	}
	if b11-a11 > 365 {
		leapOff := GetLeapMonthOffset(a11, timezone)
		leapM := leapOff - 2
		if leapM < 0 {
			leapM += 12
		}
		if lunarLeap != 0 && lunarM != leapM {
			return 0, 0, 0
		}
		if lunarLeap != 0 || off >= leapOff {
			off += 1
		}
	}
	monthStart := GetNewMoonDay(k+off, timezone)
	return jdToDate(monthStart + lunarD - 1)
}
