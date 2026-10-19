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
			Measurements{Length: &Length{Amount: 1, Unit: UnitMeters}},
			Measurements{Length: &Length{Amount: 1, Unit: UnitMeters}},
		},
		{
			Measurements{Area: &Area{Amount: 1, Unit: UnitHectares}},
			Measurements{Area: &Area{Amount: 1, Unit: UnitHectares}},
		},
		{
			Measurements{Energy: &Energy{Amount: 100, Unit: UnitKiloCalories}},
			Measurements{Energy: &Energy{Amount: 100, Unit: UnitKiloCalories}},
		},
		{
			Measurements{Duration: &Duration{Amount: 30, Unit: UnitMinutes}},
			Measurements{Duration: &Duration{Amount: 30, Unit: UnitMinutes}},
		},
		{
			Measurements{
				Quantity: 2,
				Mass:     &Mass{Amount: 500, Unit: UnitMilligrams},
				Volume:   &Volume{Amount: 1, Unit: UnitLiters},
				Length:   &Length{Amount: 30, Unit: UnitCentiMeters},
				Area:     &Area{Amount: 2, Unit: UnitSquareMeters},
				Energy:   &Energy{Amount: 250, Unit: UnitKiloCalories},
				Duration: &Duration{Amount: 45, Unit: UnitMinutes},
			},
			Measurements{
				Quantity: 2,
				Mass:     &Mass{Amount: 500, Unit: UnitMilligrams},
				Volume:   &Volume{Amount: 1, Unit: UnitLiters},
				Length:   &Length{Amount: 30, Unit: UnitCentiMeters},
				Area:     &Area{Amount: 2, Unit: UnitSquareMeters},
				Energy:   &Energy{Amount: 250, Unit: UnitKiloCalories},
				Duration: &Duration{Amount: 45, Unit: UnitMinutes},
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
			Measurements{},
			Measurements{Length: &Length{Amount: 1, Unit: UnitMeters}},
		},
		{
			Measurements{},
			Measurements{Area: &Area{Amount: 1, Unit: UnitSquareMeters}},
		},
		{
			Measurements{},
			Measurements{Energy: &Energy{Amount: 1, Unit: UnitJoules}},
		},
		{
			Measurements{},
			Measurements{Duration: &Duration{Amount: 1, Unit: UnitSeconds}},
		},
		{
			Measurements{Length: &Length{Amount: 1, Unit: UnitMeters}},
			Measurements{Length: &Length{Amount: 1, Unit: UnitMilliMeters}},
		},
		{
			Measurements{Area: &Area{Amount: 1, Unit: UnitSquareMeters}},
			Measurements{Area: &Area{Amount: 1, Unit: UnitSquareCentiMeters}},
		},
		{
			Measurements{Energy: &Energy{Amount: 1, Unit: UnitJoules}},
			Measurements{Energy: &Energy{Amount: 1, Unit: UnitKiloJoules}},
		},
		{
			Measurements{Duration: &Duration{Amount: 1, Unit: UnitSeconds}},
			Measurements{Duration: &Duration{Amount: 1, Unit: UnitMinutes}},
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
