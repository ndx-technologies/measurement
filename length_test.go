package measurement

import (
	"fmt"
	"math"
	"testing"
)

func ExampleNewLengthFromString() {
	v, _ := NewLengthFromString("25.4mm")
	fmt.Println(v.Amount, v.Unit)
	// Output: 25.4 mm
}

func TestLength(t *testing.T) {
	tests := map[string]Length{
		"1m":    {Amount: 1, Unit: UnitMeters},
		"1km":   {Amount: 1, Unit: UnitKiloMeters},
		"1in":   {Amount: 1, Unit: UnitInches},
		"12ft":  {Amount: 12, Unit: UnitFeet},
		"1.5mi": {Amount: 1.5, Unit: UnitMiles},
		"1nmi":  {Amount: 1, Unit: UnitNauticalMiles},
	}
	for s, v := range tests {
		t.Run(s, func(t *testing.T) {
			u, err := NewLengthFromString(s)
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

func TestLengthConversion_Exact(t *testing.T) {
	tests := [][2]Length{
		{{0, UnitMeters}, {0, UnitMiles}},
		{{1, UnitMeters}, {1, UnitMeters}},
		{{1000, UnitMilliMeters}, {1, UnitMeters}},
		{{100, UnitCentiMeters}, {1, UnitMeters}},
		{{10, UnitDeciMeters}, {1, UnitMeters}},
		{{1000, UnitMeters}, {1, UnitKiloMeters}},
		{{12, UnitInches}, {1, UnitFeet}},
		{{3, UnitFeet}, {1, UnitYards}},
		{{1760, UnitYards}, {1, UnitMiles}},
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

func TestLengthConversion_UnknownUnit(t *testing.T) {
	l := Length{Amount: 5, Unit: UnitMeters}
	if got := l.Convert(UnitLengthUnknown); got != (Length{Amount: 5, Unit: UnitLengthUnknown}) {
		t.Error(got)
	}
}

const lengthConversionPrecision = 0.000001

func TestLengthConversion_Approx(t *testing.T) {
	tests := [][2]Length{
		{{1, UnitInches}, {25.4, UnitMilliMeters}},
		{{1, UnitFeet}, {30.48, UnitCentiMeters}},
		{{1, UnitYards}, {0.9144, UnitMeters}},
		{{1, UnitDeciMeters}, {3.937007874015748, UnitInches}},
		{{1, UnitMiles}, {1.609344, UnitKiloMeters}},
		{{1, UnitNauticalMiles}, {1852, UnitMeters}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		if c := a.Convert(b.Unit); math.Abs(c.Amount-b.Amount)/c.Amount > lengthConversionPrecision {
			t.Error(c, b)
		}
		if c := b.Convert(a.Unit); math.Abs(c.Amount-a.Amount)/c.Amount > lengthConversionPrecision {
			t.Error(c, a)
		}
	}
}
