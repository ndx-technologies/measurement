package measurement

import (
	"fmt"
	"math"
	"testing"
)

func ExampleNewDataStorageFromString() {
	v, _ := NewDataStorageFromString("512GB")
	fmt.Println(v.Amount, v.Unit)
	// Output: 512 GB
}

func TestDataStorage(t *testing.T) {
	tests := map[string]DataStorage{
		"500B":    {Amount: 500, Unit: UnitBytes},
		"1kB":     {Amount: 1, Unit: UnitKiloBytes},
		"16GB":    {Amount: 16, Unit: UnitGigaBytes},
		"1TB":     {Amount: 1, Unit: UnitTeraBytes},
		"1024KiB": {Amount: 1024, Unit: UnitKibiBytes},
		"8GiB":    {Amount: 8, Unit: UnitGibiBytes},
	}
	for s, v := range tests {
		t.Run(s, func(t *testing.T) {
			u, err := NewDataStorageFromString(s)
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

func TestDataStorageConversion_Exact(t *testing.T) {
	tests := [][2]DataStorage{
		{{0, UnitBytes}, {0, UnitGibiBytes}},
		{{1, UnitBytes}, {1, UnitBytes}},
		{{1000, UnitBytes}, {1, UnitKiloBytes}},
		{{1000, UnitKiloBytes}, {1, UnitMegaBytes}},
		{{1000, UnitMegaBytes}, {1, UnitGigaBytes}},
		{{1000, UnitGigaBytes}, {1, UnitTeraBytes}},
		{{1000, UnitTeraBytes}, {1, UnitPetaBytes}},
		{{1024, UnitBytes}, {1, UnitKibiBytes}},
		{{1024, UnitKibiBytes}, {1, UnitMebiBytes}},
		{{1024, UnitMebiBytes}, {1, UnitGibiBytes}},
		{{1024, UnitGibiBytes}, {1, UnitTebiBytes}},
		{{1024, UnitTebiBytes}, {1, UnitPebiBytes}},
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

const dataStorageConversionPrecision = 0.000001

func TestDataStorageConversion_Approx(t *testing.T) {
	tests := [][2]DataStorage{
		{{1, UnitGigaBytes}, {0.9313225746154785, UnitGibiBytes}},
		{{1, UnitGibiBytes}, {1.073741824, UnitGigaBytes}},
		{{1, UnitTeraBytes}, {0.9094947017729282, UnitTebiBytes}},
		{{1, UnitKiloBytes}, {0.9765625, UnitKibiBytes}},
		{{1, UnitMegaBytes}, {0.95367431640625, UnitMebiBytes}},
		{{1, UnitPetaBytes}, {0.8881784197001252, UnitPebiBytes}},
	}
	for _, tc := range tests {
		a, b := tc[0], tc[1]
		if c := a.Convert(b.Unit); math.Abs(c.Amount-b.Amount)/c.Amount > dataStorageConversionPrecision {
			t.Error(c, b)
		}
		if c := b.Convert(a.Unit); math.Abs(c.Amount-a.Amount)/c.Amount > dataStorageConversionPrecision {
			t.Error(c, a)
		}
	}
}
