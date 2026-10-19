package measurement

import (
	"fmt"
	"math"
	"testing"
)

func ExampleNewChargeFromString() {
	v, _ := NewChargeFromString("5000mAh")
	fmt.Println(v.Amount, v.Unit)
	// Output: 5000 mAh
}

func TestCharge(t *testing.T) {
	tests := map[string]Charge{
		"5000mAh": {Amount: 5000, Unit: UnitMilliAmpereHours},
		"5Ah":     {Amount: 5, Unit: UnitAmpereHours},
		"1kAh":    {Amount: 1, Unit: UnitKiloAmpereHours},
		"3600C":   {Amount: 3600, Unit: UnitCoulombs},
	}
	for s, v := range tests {
		t.Run(s, func(t *testing.T) {
			u, err := NewChargeFromString(s)
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

func TestChargeConversion_Exact(t *testing.T) {
	tests := [][2]Charge{
		{{0, UnitMilliAmpereHours}, {0, UnitCoulombs}},
		{{1, UnitMilliAmpereHours}, {1, UnitMilliAmpereHours}},
		{{1000, UnitMilliAmpereHours}, {1, UnitAmpereHours}},
		{{1000, UnitAmpereHours}, {1, UnitKiloAmpereHours}},
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

const chargeConversionPrecision = 0.000001

func TestChargeConversion_Approx(t *testing.T) {
	tests := [][2]Charge{
		{{1, UnitMilliAmpereHours}, {3.6, UnitCoulombs}},
		{{1, UnitAmpereHours}, {3600, UnitCoulombs}},
		{{1, UnitKiloAmpereHours}, {3_600_000, UnitCoulombs}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		if c := a.Convert(b.Unit); math.Abs(c.Amount-b.Amount)/c.Amount > chargeConversionPrecision {
			t.Error(c, b)
		}
		if c := b.Convert(a.Unit); math.Abs(c.Amount-a.Amount)/c.Amount > chargeConversionPrecision {
			t.Error(c, a)
		}
	}
}
