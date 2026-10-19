package measurement

import (
	"fmt"
	"math"
	"testing"
)

func ExampleNewPowerFromString() {
	v, _ := NewPowerFromString("1.5kW")
	fmt.Println(v.Amount, v.Unit)
	// Output: 1.5 kW
}

func TestPower(t *testing.T) {
	tests := map[string]Power{
		"100W":  {Amount: 100, Unit: UnitWatts},
		"1.5kW": {Amount: 1.5, Unit: UnitKiloWatts},
		"2MW":   {Amount: 2, Unit: UnitMegaWatts},
		"1GW":   {Amount: 1, Unit: UnitGigaWatts},
		"1hp":   {Amount: 1, Unit: UnitHorsepower},
	}
	for s, v := range tests {
		t.Run(s, func(t *testing.T) {
			u, err := NewPowerFromString(s)
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

func TestPowerConversion_Exact(t *testing.T) {
	tests := [][2]Power{
		{{0, UnitWatts}, {0, UnitHorsepower}},
		{{1, UnitWatts}, {1, UnitWatts}},
		{{1000, UnitWatts}, {1, UnitKiloWatts}},
		{{1000, UnitKiloWatts}, {1, UnitMegaWatts}},
		{{1000, UnitMegaWatts}, {1, UnitGigaWatts}},
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

const powerConversionPrecision = 0.000001

func TestPowerConversion_Approx(t *testing.T) {
	tests := [][2]Power{
		{{1, UnitHorsepower}, {745.6998715822702, UnitWatts}},
		{{1, UnitHorsepower}, {0.7456998715822702, UnitKiloWatts}},
		{{1, UnitHorsepower}, {0.0007456998715822702, UnitMegaWatts}},
		{{1, UnitHorsepower}, {0.0000007456998715822702, UnitGigaWatts}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		if c := a.Convert(b.Unit); math.Abs(c.Amount-b.Amount)/c.Amount > powerConversionPrecision {
			t.Error(c, b)
		}
		if c := b.Convert(a.Unit); math.Abs(c.Amount-a.Amount)/c.Amount > powerConversionPrecision {
			t.Error(c, a)
		}
	}
}
