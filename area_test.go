package measurement

import (
	"fmt"
	"math"
	"testing"
)

func ExampleNewAreaFromString() {
	v, _ := NewAreaFromString("40m2")
	fmt.Println(v.Amount, v.Unit)
	// Output: 40 m2
}

func TestArea(t *testing.T) {
	tests := map[string]Area{
		"1m2":   {Amount: 1, Unit: UnitSquareMeters},
		"1ha":   {Amount: 1, Unit: UnitHectares},
		"1km2":  {Amount: 1, Unit: UnitSquareKiloMeters},
		"1in2":  {Amount: 1, Unit: UnitSquareInches},
		"1ft2":  {Amount: 1, Unit: UnitSquareFeet},
		"1acre": {Amount: 1, Unit: UnitAcres},
	}
	for s, v := range tests {
		t.Run(s, func(t *testing.T) {
			u, err := NewAreaFromString(s)
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

func TestAreaConversion_Exact(t *testing.T) {
	tests := [][2]Area{
		{{0, UnitSquareMeters}, {0, UnitAcres}},
		{{1, UnitSquareMeters}, {1, UnitSquareMeters}},
		{{100, UnitSquareMilliMeters}, {1, UnitSquareCentiMeters}},
		{{10_000, UnitSquareCentiMeters}, {1, UnitSquareMeters}},
		{{10_000, UnitSquareMeters}, {1, UnitHectares}},
		{{100, UnitHectares}, {1, UnitSquareKiloMeters}},
		{{144, UnitSquareInches}, {1, UnitSquareFeet}},
		{{9, UnitSquareFeet}, {1, UnitSquareYards}},
		{{4840, UnitSquareYards}, {1, UnitAcres}},
		{{640, UnitAcres}, {1, UnitSquareMiles}},
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

func TestAreaConversion_UnknownUnit(t *testing.T) {
	a := Area{Amount: 5, Unit: UnitSquareMeters}
	if got := a.Convert(UnitAreaUnknown); got != (Area{Amount: 5, Unit: UnitAreaUnknown}) {
		t.Error(got)
	}
}

const areaConversionPrecision = 0.000001

func TestAreaConversion_Approx(t *testing.T) {
	tests := [][2]Area{
		{{1, UnitSquareInches}, {6.4516, UnitSquareCentiMeters}},
		{{1, UnitSquareMilliMeters}, {0.0015500031, UnitSquareInches}},
		{{1, UnitSquareCentiMeters}, {0.15500031, UnitSquareInches}},
		{{1, UnitSquareFeet}, {929.0304, UnitSquareCentiMeters}},
		{{1, UnitSquareYards}, {0.83612736, UnitSquareMeters}},
		{{1, UnitHectares}, {2.471053814671653, UnitAcres}},
		{{1, UnitAcres}, {4046.8564224, UnitSquareMeters}},
		{{1, UnitAcres}, {0.40468564224, UnitHectares}},
		{{1, UnitSquareKiloMeters}, {247.1053814671653, UnitAcres}},
		{{1, UnitSquareMiles}, {2.589988110336, UnitSquareKiloMeters}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		if c := a.Convert(b.Unit); math.Abs(c.Amount-b.Amount)/c.Amount > areaConversionPrecision {
			t.Error(c, b)
		}
		if c := b.Convert(a.Unit); math.Abs(c.Amount-a.Amount)/c.Amount > areaConversionPrecision {
			t.Error(c, a)
		}
	}
}
