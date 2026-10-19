package measurement

import (
	"fmt"
	"testing"
)

func TestMeasurementsIsEqual(t *testing.T) {
	tests := []struct{ a, b Measurements }{
		{Measurements{}, Measurements{}},
		{Measurements{Quantity: 3}, Measurements{Quantity: 3}},
		{
			Measurements{Mass: &Mass{Amount: 100, Unit: UnitGrams}},
			Measurements{Mass: &Mass{Amount: 100, Unit: UnitGrams}},
		},
		{
			Measurements{Volume: &Volume{Amount: 1, Unit: UnitLiters}},
			Measurements{Volume: &Volume{Amount: 1, Unit: UnitLiters}},
		},
		{
			Measurements{
				Quantity: 2,
				Mass:     &Mass{Amount: 500, Unit: UnitMilligrams},
				Volume:   &Volume{Amount: 1, Unit: UnitLiters},
			},
			Measurements{
				Quantity: 2,
				Mass:     &Mass{Amount: 500, Unit: UnitMilligrams},
				Volume:   &Volume{Amount: 1, Unit: UnitLiters},
			},
		},
	}

	for _, tc := range tests {
		if !tc.a.IsEqual(tc.b) {
			t.Error(tc)
		}
	}
}

func ExampleMeasurements_IsZero() {
	fmt.Println(Measurements{}.IsZero())
	fmt.Println(Measurements{Quantity: 1}.IsZero())

	// Output:
	// true
	// false
}

func ExampleMeasurements_IsEqual() {
	a := Measurements{Mass: &Mass{Amount: 100, Unit: UnitGrams}}
	b := Measurements{Mass: &Mass{Amount: 100, Unit: UnitGrams}}
	c := Measurements{Mass: &Mass{Amount: 200, Unit: UnitGrams}}

	fmt.Println(a.IsEqual(b))
	fmt.Println(a.IsEqual(c))

	// Output:
	// true
	// false
}

func TestMeasurementsIsNotEqual(t *testing.T) {
	tests := []struct{ a, b Measurements }{
		{Measurements{Quantity: 3}, Measurements{Quantity: 4}},
		{
			Measurements{Mass: &Mass{Amount: 100, Unit: UnitGrams}},
			Measurements{Mass: &Mass{Amount: 200, Unit: UnitGrams}},
		},
		{
			Measurements{Mass: &Mass{Amount: 100, Unit: UnitGrams}},
			Measurements{Mass: &Mass{Amount: 100, Unit: UnitKilograms}},
		},
		{
			Measurements{},
			Measurements{Mass: &Mass{Amount: 100, Unit: UnitGrams}},
		},
		{
			Measurements{Volume: &Volume{Amount: 1, Unit: UnitLiters}},
			Measurements{Volume: &Volume{Amount: 2, Unit: UnitLiters}},
		},
		{
			Measurements{Volume: &Volume{Amount: 1, Unit: UnitLiters}},
			Measurements{Volume: &Volume{Amount: 1, Unit: UnitMilliLiters}},
		},
		{
			Measurements{},
			Measurements{Volume: &Volume{Amount: 1, Unit: UnitLiters}},
		},
		{
			Measurements{
				Quantity: 2,
				Mass:     &Mass{Amount: 500, Unit: UnitMilligrams},
				Volume:   &Volume{Amount: 1, Unit: UnitLiters},
			},
			Measurements{
				Quantity: 2,
				Mass:     &Mass{Amount: 500, Unit: UnitMilligrams},
				Volume:   &Volume{Amount: 2, Unit: UnitLiters},
			},
		},
	}

	for _, tc := range tests {
		if tc.a.IsEqual(tc.b) {
			t.Error(tc)
		}
	}
}
