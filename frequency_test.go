package measurement

import (
	"fmt"
	"testing"
)

func ExampleNewFrequencyFromString() {
	v, _ := NewFrequencyFromString("2.4GHz")
	fmt.Println(v.Amount, v.Unit)
	// Output: 2.4 GHz
}

func TestFrequency(t *testing.T) {
	tests := map[string]Frequency{
		"50Hz":   {Amount: 50, Unit: UnitHertz},
		"100kHz": {Amount: 100, Unit: UnitKiloHertz},
		"2.4GHz": {Amount: 2.4, Unit: UnitGigaHertz},
		"1MHz":   {Amount: 1, Unit: UnitMegaHertz},
		"1THz":   {Amount: 1, Unit: UnitTeraHertz},
	}
	for s, v := range tests {
		t.Run(s, func(t *testing.T) {
			u, err := NewFrequencyFromString(s)
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

func TestFrequencyConversion_UnknownUnit(t *testing.T) {
	f := Frequency{Amount: 5, Unit: UnitHertz}
	if got := f.Convert(UnitFrequencyUnknown); got != (Frequency{Amount: 5, Unit: UnitFrequencyUnknown}) {
		t.Error(got)
	}
}

func TestFrequencyConversion_Exact(t *testing.T) {
	tests := [][2]Frequency{
		{{0, UnitHertz}, {0, UnitTeraHertz}},
		{{1, UnitHertz}, {1, UnitHertz}},
		{{1000, UnitHertz}, {1, UnitKiloHertz}},
		{{1000, UnitKiloHertz}, {1, UnitMegaHertz}},
		{{1000, UnitMegaHertz}, {1, UnitGigaHertz}},
		{{1000, UnitGigaHertz}, {1, UnitTeraHertz}},
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
