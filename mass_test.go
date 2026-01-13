package measurement

import (
	"fmt"
	"math"
	"testing"
)

func ExampleNewMassFromString() {
	m, _ := NewMassFromString("420g")
	fmt.Println(m.Amount, m.Unit)
	// Output: 420 g
}

func TestMassEncoding(t *testing.T) {
	tests := []struct {
		s string
		v Mass
	}{
		{"1kg", Mass{Amount: 1, Unit: UnitKilograms}},
		{"1g", Mass{Amount: 1, Unit: UnitGrams}},
		{"0.123kg", Mass{Amount: 0.123, Unit: UnitKilograms}},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			v, err := NewMassFromString(tt.s)
			if err != nil {
				t.Error(err)
			}
			if v == nil {
				t.Fail()
			}
			if *v != tt.v {
				t.Error(v, tt.v)
			}
			if v.String() != tt.s {
				t.Error(v.String(), tt.s)
			}
		})
	}
}

func TestMassConversion_Exact(t *testing.T) {
	tests := [][2]Mass{
		{{0, UnitKilograms}, {0, UnitOunces}},
		{{1, UnitKilograms}, {1, UnitKilograms}},
		{{1000, UnitGrams}, {1, UnitKilograms}},
		{{200, UnitGrams}, {0.2, UnitKilograms}},
		{{10, UnitMilligrams}, {0.00001, UnitKilograms}},
		{{333, UnitGrams}, {0.333, UnitKilograms}},
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

const massConversionPrecision = 0.0001

func TestMassConversion_Approx(t *testing.T) {
	tests := [][2]Mass{
		{{35.274, UnitOunces}, {1, UnitKilograms}},
		{{1, UnitOunces}, {28.349, UnitGrams}},
		{{1, UnitOunces}, {0.0283495, UnitKilograms}},
		{{1, UnitMetricTons}, {1_000_000, UnitGrams}},
		{{1, UnitMetricTons}, {1_000_000_000, UnitMilligrams}},
		{{1, UnitMetricTons}, {1_000_000_000_000, UnitMicrograms}},
		{{1, UnitMetricTons}, {1_000_000_000_000_000, UnitNanograms}},
		// Carat conversions (1 ct = 200 mg = 0.2 g)
		{{1, UnitCarats}, {0.2, UnitGrams}},
		{{5, UnitCarats}, {1, UnitGrams}},
		{{1, UnitCarats}, {200, UnitMilligrams}},
		// Pounds (1 lb = 453.59237 g)
		{{1, UnitPounds}, {453.59237, UnitGrams}},
		{{2.20462, UnitPounds}, {1, UnitKilograms}},
		// Stones (1 st = 6350.29318 g)
		{{1, UnitStones}, {6350.29318, UnitGrams}},
		// Troy ounces (1 ozt = 31.1035 g)
		{{1, UnitOuncesTroy}, {31.1035, UnitGrams}},
		// Slugs (1 slug = 14593.9029 g)
		{{1, UnitSlugs}, {14593.9029, UnitGrams}},
		// Short tons (1 sst = 907184.74 g)
		{{1, UnitShortTons}, {907184.74, UnitGrams}},
		// Picograms ladder
		{{1, UnitGrams}, {1_000_000_000_000, UnitPicograms}},
		// Centigrams and decigrams
		{{1, UnitGrams}, {100, UnitCentigrams}},
		{{1, UnitGrams}, {10, UnitDecigrams}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		if c := a.Convert(b.Unit); math.Abs(c.Amount-b.Amount)/c.Amount > massConversionPrecision {
			t.Error(c, b)
		}
		if c := b.Convert(a.Unit); math.Abs(c.Amount-a.Amount)/c.Amount > massConversionPrecision {
			t.Error(c, a)
		}
	}
}

// TestTryConvertExactMass_Int64Precision demonstrates that int64 preserves
// exact values where float64 would lose precision (beyond 2^53 ≈ 9e15).
func TestTryConvertExactMass_Int64Precision(t *testing.T) {
	type M struct {
		amount int64
		unit   UnitMass
	}

	tests := [][2]M{
		// Large metric conversions that stay exact with int64
		{{1_000_000_000_000_000, UnitPicograms}, {1_000, UnitGrams}},
		{{1_000_000_000_000, UnitNanograms}, {1_000, UnitGrams}},
		{{1_000_000_000, UnitMicrograms}, {1_000, UnitGrams}},
		{{1_000_000, UnitMilligrams}, {1, UnitKilograms}},
		// Carat ladder: 1g = 5ct, 1ct = 200mg
		{{5_000_000, UnitCarats}, {1, UnitMetricTons}},
		{{1_000_000_000, UnitMilligrams}, {1, UnitMetricTons}},
		// Picogram to ton spans 10^18 - needs int64
		{{1_000_000_000_000_000_000, UnitPicograms}, {1, UnitMetricTons}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		v, ok := TryConvertExactMass(a.amount, a.unit, b.unit)
		if !ok {
			t.Error(v, b.amount)
		}
		if v != b.amount {
			t.Error(v, b.amount)
		}

		back, ok := TryConvertExactMass(v, b.unit, a.unit)
		if !ok || back != a.amount {
			t.Error("round-trip failed:", a.amount, "->", v, "->", back)
		}
	}
}

func TestTryConvertExactMass_CrossSystemFails(t *testing.T) {
	tests := [][2]UnitMass{
		{UnitGrams, UnitOunces},
		{UnitKilograms, UnitPounds},
		{UnitGrams, UnitOuncesTroy},
		{UnitKilograms, UnitSlugs},
		{UnitMetricTons, UnitShortTons},
	}
	for _, tc := range tests {
		if _, ok := TryConvertExactMass(int64(1), tc[0], tc[1]); ok {
			t.Error(tc)
		}
	}
}
