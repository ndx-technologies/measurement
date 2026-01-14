package measurement

import (
	"fmt"
	"math"
	"testing"
)

func ExampleNewVolumeFromString() {
	v, _ := NewVolumeFromString("150ml")
	fmt.Println(v.Amount, v.Unit)
	// Output: 150 ml
}

func TestVolume(t *testing.T) {
	tests := map[string]Volume{
		"1l":     {Amount: 1, Unit: UnitLiters},
		"1ml":    {Amount: 1, Unit: UnitMilliLiters},
		"0.123l": {Amount: 0.123, Unit: UnitLiters},
	}
	for s, v := range tests {
		t.Run(s, func(t *testing.T) {
			u, err := NewVolumeFromString(s)
			if err != nil {
				t.Error(err.Error())
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

func TestVolumeConversion_Exact(t *testing.T) {
	tests := [][2]Volume{
		{{0, UnitLiters}, {0, UnitBushels}},
		{{1, UnitLiters}, {1000, UnitMilliLiters}},
		{{1, UnitPints}, {96, UnitTeaspoons}},
		{{0.01, UnitLiters}, {10, UnitMilliLiters}},
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

const volumeConversionPrecision = 0.000001

func TestVolumeConversion_Approx(t *testing.T) {
	tests := [][2]Volume{
		{{1, UnitPints}, {0.4731762648307425, UnitLiters}},
		{{1, UnitFluidOunces}, {1.80469, UnitCubicInches}},
		{{1, UnitMegaLiters}, {1_000_000_000, UnitMilliLiters}},
		{{1, UnitCubicMiles}, {254_358_061_056_000, UnitCubicInches}},
		// Liter ladder
		{{1, UnitLiters}, {10, UnitDeciLiters}},
		{{1, UnitLiters}, {100, UnitCentiLiters}},
		{{1, UnitKiloLiters}, {1000, UnitLiters}},
		{{1, UnitMegaLiters}, {1000, UnitKiloLiters}},
		// Cubic meter ladder
		{{1, UnitCubicDeciMeters}, {1000, UnitCubicCentiMeters}},
		{{1, UnitCubicMeters}, {1000, UnitCubicDeciMeters}},
		{{1, UnitCubicCentiMeters}, {1000, UnitCubicMilliMeters}},
		// Cubic inch ladder
		{{1, UnitCubicFeet}, {1728, UnitCubicInches}},
		{{1, UnitCubicYards}, {27, UnitCubicFeet}},
		{{27, UnitCubicFeet}, {46656, UnitCubicInches}},
		// US customary ladder
		{{1, UnitTablespoons}, {3, UnitTeaspoons}},
		{{1, UnitFluidOunces}, {2, UnitTablespoons}},
		{{1, UnitCups}, {8, UnitFluidOunces}},
		{{1, UnitPints}, {2, UnitCups}},
		{{1, UnitQuarts}, {2, UnitPints}},
		{{1, UnitGallons}, {4, UnitQuarts}},
		{{1, UnitGallons}, {128, UnitFluidOunces}},
		// Imperial ladder
		{{1, UnitImperialTablespoons}, {3, UnitImperialTeaspoons}},
		{{1, UnitImperialFluidOunces}, {2, UnitImperialTablespoons}},
		{{1, UnitImperialGills}, {5, UnitImperialFluidOunces}},
		{{1, UnitImperialPints}, {4, UnitImperialGills}},
		{{1, UnitImperialQuarts}, {2, UnitImperialPints}},
		{{1, UnitImperialGallons}, {4, UnitImperialQuarts}},
		{{1, UnitBushels}, {8, UnitImperialGallons}},
		// Cross-system
		{{1, UnitGallons}, {3.785408, UnitLiters}},
		{{1, UnitImperialPints}, {0.568261, UnitLiters}},
		{{1, UnitImperialGallons}, {4.546088, UnitLiters}},
		{{1, UnitCubicFeet}, {28.3168, UnitLiters}},
		// dm³ = L
		{{1, UnitCubicDeciMeters}, {1, UnitLiters}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		if c := a.Convert(b.Unit); math.Abs(c.Amount-b.Amount)/c.Amount > volumeConversionPrecision {
			t.Error(c, b)
		}
		if c := b.Convert(a.Unit); math.Abs(c.Amount-a.Amount)/c.Amount > volumeConversionPrecision {
			t.Error(c, a)
		}
	}
}

// TestTryConvertExactVolume_Int64Precision demonstrates that int64 preserves
// exact values where float64 would lose precision (beyond 2^53 ≈ 9e15).
func TestTryConvertExactVolume_Int64Precision(t *testing.T) {
	type V struct {
		amount int64
		unit   UnitVolume
	}

	tests := [][2]V{
		// Large metric conversions that stay exact with int64
		{{1_000_000_000_000, UnitMilliLiters}, {1_000_000, UnitKiloLiters}},
		{{1_000_000_000, UnitMilliLiters}, {1_000_000, UnitLiters}},
		// Cubic volumes: mm³->cm³->dm³->m³ each has factor 1000
		// So mm³ to m³ is 10^9
		{{1_000_000_000_000, UnitCubicMilliMeters}, {1_000, UnitCubicMeters}},
		{{1_000_000_000, UnitCubicCentiMeters}, {1_000, UnitCubicMeters}},
		// US customary: 1 gallon = 768 teaspoons (exact integer)
		{{768_000_000, UnitTeaspoons}, {1_000_000, UnitGallons}},
		// Cubic miles have huge factor: 1760^3 = 5,451,776,000
		{{5_451_776_000, UnitCubicYards}, {1, UnitCubicMiles}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		v, ok := TryConvertExactVolume(a.amount, a.unit, b.unit)
		if !ok {
			t.Error(v, b.amount)
		}
		if v != b.amount {
			t.Error(v, b.amount)
		}

		back, ok := TryConvertExactVolume(v, b.unit, a.unit)
		if !ok || back != a.amount {
			t.Error("round-trip failed:", a.amount, "->", v, "->", back)
		}
	}
}

func TestTryConvertExactVolume_CrossSystemFails(t *testing.T) {
	tests := [][2]UnitVolume{
		{UnitLiters, UnitGallons},
		{UnitMilliLiters, UnitFluidOunces},
		{UnitCubicMeters, UnitCubicFeet},
		{UnitImperialPints, UnitPints},
	}
	for _, tc := range tests {
		if _, ok := TryConvertExactVolume(int64(1), tc[0], tc[1]); ok {
			t.Error(tc)
		}
	}
}

func TestTryConvertExactVolume(t *testing.T) {
	type VolumeInt64 struct {
		Amount int64
		Unit   UnitVolume
	}

	t.Run("exact", func(t *testing.T) {
		tests := [][2]VolumeInt64{
			{{1000, UnitMilliLiters}, {1, UnitLiters}},
			{{1, UnitLiters}, {1000, UnitMilliLiters}},
			{{10, UnitMilliLiters}, {1, UnitCentiLiters}},
		}
		for _, tc := range tests {
			if v, ok := TryConvertExactVolume(tc[0].Amount, tc[0].Unit, tc[1].Unit); !ok || v != tc[1].Amount {
				t.Error(tc, v, ok)
			}
		}
	})

	t.Run("not exact", func(t *testing.T) {
		tests := []VolumeInt64{
			{5, UnitMilliLiters},
			{100, UnitMilliLiters},
		}
		for _, tc := range tests {
			if v, ok := TryConvertExactVolume(tc.Amount, tc.Unit, UnitLiters); ok {
				t.Error(tc, v)
			}
		}
	})
}
