package measurement

import (
	"fmt"
	"math"
	"testing"
)

func ExampleNewEnergyFromString() {
	v, _ := NewEnergyFromString("100kcal")
	fmt.Println(v.Amount, v.Unit)
	// Output: 100 kcal
}

func TestEnergy(t *testing.T) {
	tests := map[string]Energy{
		"100J":    {Amount: 100, Unit: UnitJoules},
		"500kJ":   {Amount: 500, Unit: UnitKiloJoules},
		"2MJ":     {Amount: 2, Unit: UnitMegaJoules},
		"100kcal": {Amount: 100, Unit: UnitKiloCalories},
		"1.5kWh":  {Amount: 1.5, Unit: UnitKiloWattHours},
		"5000btu": {Amount: 5000, Unit: UnitBTU},
	}
	for s, v := range tests {
		t.Run(s, func(t *testing.T) {
			u, err := NewEnergyFromString(s)
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

func TestEnergyConversion_Exact(t *testing.T) {
	tests := [][2]Energy{
		{{0, UnitJoules}, {0, UnitBTU}},
		{{1, UnitJoules}, {1, UnitJoules}},
		{{1000, UnitJoules}, {1, UnitKiloJoules}},
		{{1000, UnitKiloJoules}, {1, UnitMegaJoules}},
		{{1000, UnitCalories}, {1, UnitKiloCalories}},
		{{1000, UnitWattHours}, {1, UnitKiloWattHours}},
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

const energyConversionPrecision = 0.000001

func TestEnergyConversion_Approx(t *testing.T) {
	tests := [][2]Energy{
		{{1, UnitCalories}, {4.184, UnitJoules}},
		{{1, UnitKiloCalories}, {4.184, UnitKiloJoules}},
		{{1, UnitWattHours}, {3600, UnitJoules}},
		{{1, UnitWattHours}, {3.6, UnitKiloJoules}},
		{{1, UnitKiloWattHours}, {3.6, UnitMegaJoules}},
		{{1, UnitBTU}, {1055.05585262, UnitJoules}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		if c := a.Convert(b.Unit); math.Abs(c.Amount-b.Amount)/c.Amount > energyConversionPrecision {
			t.Error(c, b)
		}
		if c := b.Convert(a.Unit); math.Abs(c.Amount-a.Amount)/c.Amount > energyConversionPrecision {
			t.Error(c, a)
		}
	}
}
