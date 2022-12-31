package calendar

import (
	"testing"
)

func TestJdFromDate(t *testing.T) {
	jd := jdFromDate(1, 1, 2000)
	if jd != 2451545 {
		t.Errorf("jd of 1/1/2000 is %d not equal 2451545", jd)
	}
}

func TestJdToDate(t *testing.T) {
	dd, mm, yy := jdToDate(2451545)
	if dd != 1 || mm != 1 || yy != 2000 {
		t.Errorf("2451545 is %d/%d/%d not equal 1/1/2000", dd, mm, yy)
	}
}

func TestNewMoon(t *testing.T) {
	tcs := map[int]float64{
		2:    2415079.976104907,
		-2:   2414961.93439546,
		1533: 2460291.480190389,
	}
	for k, v := range tcs {
		r := newMoon(k)
		if r != v {
			t.Errorf("newMoon(%d)=%f but expected %f", k, r, v)
		}
	}
}

func TestGetNewMoonDay(t *testing.T) {
	type c struct {
		K        int
		Timezone int
		Expected int
	}
	tcs := []c{
		{K: 2, Timezone: 7, Expected: 2415080},
		{K: -2, Timezone: 7, Expected: 2414962},
		{K: 1533, Timezone: 7, Expected: 2460292},
	}
	for _, tc := range tcs {
		r := getNewMoonDay(tc.K, tc.Timezone)
		if r != tc.Expected {
			t.Errorf("getNewMoonDay(%d, %d)=%d but expected %d", tc.K, tc.Timezone, r, tc.Expected)
		}
	}
}

func TestSunLongitude(t *testing.T) {
	type c struct {
		Jdn      float64
		Expected float64
	}
	tcs := []c{
		{Jdn: 120300, Expected: -6.013479151413776},
		{Jdn: 100000, Expected: -3.362306698647444},
	}
	for _, tc := range tcs {
		r := sunLongitude(tc.Jdn)
		if r != tc.Expected {
			t.Errorf("sunLongtitude(%f)=%f but expected %f", tc.Jdn, r, tc.Expected)
		}
	}
}

func TestGetSunLongitude(t *testing.T) {
	type c struct {
		DayNumber int
		Timezone  int
		Expected  int
	}
	tcs := []c{
		{DayNumber: 1000, Timezone: 7, Expected: -7},
		{DayNumber: 12, Timezone: 7, Expected: -3},
	}
	for _, tc := range tcs {
		r := getSunLongitude(tc.DayNumber, tc.Timezone)
		if r != tc.Expected {
			t.Errorf("getSunLongitude(%d, %d)=%d but expected %d", tc.DayNumber, tc.Timezone, r, tc.Expected)
		}
	}
}

func TestGetLunarMonth11(t *testing.T) {
	type c struct {
		YY       int
		Timezone int
		Expected int
	}
	tcs := []c{
		{YY: 2023, Timezone: 7, Expected: 2460292},
		{YY: 1985, Timezone: 7, Expected: 2446412},
		{YY: 2000, Timezone: 8, Expected: 2451875},
	}
	for _, tc := range tcs {
		r := getLunarMonth11(tc.YY, tc.Timezone)
		if r != tc.Expected {
			t.Errorf("getLunarMonth11(%d, %d)=%d but expected %d", tc.YY, tc.Timezone, r, tc.Expected)
		}
	}
}

func TestGetLeapMonthOffset(t *testing.T) {
	type c struct {
		A11      int
		Timezone int
		Expected int
	}
	tcs := []c{
		{A11: 1, Timezone: 7, Expected: 12},
		{A11: 5, Timezone: 8, Expected: 12},
		{A11: 31, Timezone: 7, Expected: 11},
	}
	for _, tc := range tcs {
		r := getLeapMonthOffset(tc.A11, tc.Timezone)
		if r != tc.Expected {
			t.Errorf("getLeapMonthOffset(%d, %d)=%d but expected %d", tc.A11, tc.Timezone, r, tc.Expected)
		}
	}
}

func TestS2L(t *testing.T) {
	type c struct {
		DD           int
		MM           int
		YY           int
		Timezone     int
		ExpectedDD   int
		ExpectedMM   int
		ExpectedYY   int
		ExpectedLeap int
	}
	tcs := []c{
		{DD: 11, MM: 2, YY: 1985, Timezone: 7, ExpectedDD: 22, ExpectedMM: 12, ExpectedYY: 1984, ExpectedLeap: 0},
		{DD: 2, MM: 12, YY: 2022, Timezone: 7, ExpectedDD: 9, ExpectedMM: 11, ExpectedYY: 2022, ExpectedLeap: 0},
		{DD: 12, MM: 3, YY: 1979, Timezone: 7, ExpectedDD: 15, ExpectedMM: 1, ExpectedYY: 1979, ExpectedLeap: 0},
		{DD: 1, MM: 1, YY: 2023, Timezone: 7, ExpectedDD: 10, ExpectedMM: 12, ExpectedYY: 2022, ExpectedLeap: 0},
	}
	for _, tc := range tcs {
		dd, mm, yy, leap := S2L(tc.DD, tc.MM, tc.YY, tc.Timezone)
		if (dd != tc.ExpectedDD) || (mm != tc.ExpectedMM) || (yy != tc.ExpectedYY) || (leap != tc.ExpectedLeap) {
			t.Errorf("S2L(%d, %d, %d, %d)=(%d, %d, %d, %d) but expected (%d, %d, %d, %d)", tc.DD, tc.MM, tc.YY, tc.Timezone, dd, mm, yy, leap, tc.ExpectedDD, tc.ExpectedMM, tc.ExpectedYY, tc.ExpectedLeap)
		}
	}
}

func TestL2S(t *testing.T) {
	type c struct {
		DD         int
		MM         int
		YY         int
		Leap       int
		Timezone   int
		ExpectedDD int
		ExpectedMM int
		ExpectedYY int
	}
	tcs := []c{
		{DD: 15, MM: 1, YY: 1979, Leap: 0, Timezone: 7, ExpectedDD: 12, ExpectedMM: 3, ExpectedYY: 1979},
		{DD: 22, MM: 12, YY: 1984, Leap: 0, Timezone: 7, ExpectedDD: 11, ExpectedMM: 2, ExpectedYY: 1985},
		{DD: 10, MM: 12, YY: 2022, Leap: 0, Timezone: 7, ExpectedDD: 1, ExpectedMM: 1, ExpectedYY: 2023},
	}
	for _, tc := range tcs {
		dd, mm, yy := L2S(tc.DD, tc.MM, tc.YY, tc.Leap, tc.Timezone)
		if (dd != tc.ExpectedDD) || (mm != tc.ExpectedMM) || (yy != tc.ExpectedYY) {
			t.Errorf("L2S(%d, %d, %d, %d, %d)=(%d, %d, %d) but expected (%d, %d, %d)", tc.DD, tc.MM, tc.YY, tc.Leap, tc.Timezone, dd, mm, yy, tc.ExpectedDD, tc.ExpectedMM, tc.ExpectedYY)
		}
	}
}
