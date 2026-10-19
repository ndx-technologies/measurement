package measurement

import (
	"fmt"
	"testing"
)

func ExampleNewVoltageFromString() {
	v, _ := NewVoltageFromString("220V")
	fmt.Println(v.Amount, v.Unit)
	// Output: 220 V
}

func TestVoltage(t *testing.T) {
	tests := map[string]Voltage{
		"12mV": {Amount: 12, Unit: UnitMilliVolts},
		"5V":   {Amount: 5, Unit: UnitVolts},
		"1kV":  {Amount: 1, Unit: UnitKiloVolts},
	}
	for s, v := range tests {
		t.Run(s, func(t *testing.T) {
			u, err := NewVoltageFromString(s)
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

func TestVoltageConversion_UnknownUnit(t *testing.T) {
	v := Voltage{Amount: 5, Unit: UnitVolts}
	if got := v.Convert(UnitVoltageUnknown); got != (Voltage{Amount: 5, Unit: UnitVoltageUnknown}) {
		t.Error(got)
	}
}

func TestVoltageConversion_Exact(t *testing.T) {
	tests := [][2]Voltage{
		{{0, UnitMilliVolts}, {0, UnitKiloVolts}},
		{{1, UnitMilliVolts}, {1, UnitMilliVolts}},
		{{1000, UnitMilliVolts}, {1, UnitVolts}},
		{{1000, UnitVolts}, {1, UnitKiloVolts}},
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
