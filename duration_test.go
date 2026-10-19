package measurement

import (
	"fmt"
	"math"
	"testing"
)

func ExampleNewDurationFromString() {
	v, _ := NewDurationFromString("90min")
	fmt.Println(v.Amount, v.Unit)
	// Output: 90 min
}

func TestDuration(t *testing.T) {
	tests := map[string]Duration{
		"5ns":   {Amount: 5, Unit: UnitNanoSeconds},
		"500µs": {Amount: 500, Unit: UnitMicroSeconds},
		"250ms": {Amount: 250, Unit: UnitMilliSeconds},
		"30s":   {Amount: 30, Unit: UnitSeconds},
		"90min": {Amount: 90, Unit: UnitMinutes},
		"2h":    {Amount: 2, Unit: UnitHours},
		"3d":    {Amount: 3, Unit: UnitDays},
		"2wk":   {Amount: 2, Unit: UnitWeeks},
		"3mo":   {Amount: 3, Unit: UnitMonths},
		"2y":    {Amount: 2, Unit: UnitYears},
	}
	for s, v := range tests {
		t.Run(s, func(t *testing.T) {
			u, err := NewDurationFromString(s)
			if err != nil {
				t.Error(err)
			}
			if u == nil || *u != v {
				t.Error(u, v)
			}
			if v.String() != s {
				t.Error(v.String(), s)
			}
		})
	}
}

func TestDurationConversion_UnknownUnit(t *testing.T) {
	d := Duration{Amount: 5, Unit: UnitSeconds}
	if got := d.Convert(UnitDurationUnknown); got != d {
		t.Error(got)
	}
}

func TestDurationConversion_Exact(t *testing.T) {
	tests := [][2]Duration{
		{{0, UnitSeconds}, {0, UnitWeeks}},
		{{1, UnitSeconds}, {1, UnitSeconds}},
		{{1000, UnitNanoSeconds}, {1, UnitMicroSeconds}},
		{{1000, UnitMicroSeconds}, {1, UnitMilliSeconds}},
		{{1000, UnitMilliSeconds}, {1, UnitSeconds}},
		{{60, UnitSeconds}, {1, UnitMinutes}},
		{{60, UnitMinutes}, {1, UnitHours}},
		{{24, UnitHours}, {1, UnitDays}},
		{{7, UnitDays}, {1, UnitWeeks}},
		{{12, UnitMonths}, {1, UnitYears}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		if c := a.Convert(b.Unit); c != b {
			t.Error(c, b)
		}
		if c := b.Convert(a.Unit); c != a {
			t.Error(c, a)
		}
	}
}

const durationConversionPrecision = 0.000001

func TestDurationConversion_Approx(t *testing.T) {
	tests := [][2]Duration{
		{{1, UnitMonths}, {30.436875, UnitDays}},
		{{1, UnitMonths}, {730.485, UnitHours}},
		{{1, UnitYears}, {365.2425, UnitDays}},
		{{1, UnitYears}, {52.1775, UnitWeeks}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		if c := a.Convert(b.Unit); math.Abs(c.Amount-b.Amount)/c.Amount > durationConversionPrecision {
			t.Error(c, b)
		}
		if c := b.Convert(a.Unit); math.Abs(c.Amount-a.Amount)/c.Amount > durationConversionPrecision {
			t.Error(c, a)
		}
	}
}
