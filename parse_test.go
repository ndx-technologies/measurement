package measurement

import (
	"errors"
	"testing"
)

func TestAmountUnitFromString(t *testing.T) {
	tests := map[string]struct {
		amount float64
		unit   UnitDuration
		err    error
	}{
		"30s":   {amount: 30, unit: UnitSeconds},
		"90min": {amount: 90, unit: UnitMinutes},
		"":      {err: ErrInvalidUnit{Type: MeasureTypeDuration}},
		"s":     {err: ErrInvalidAmount{Type: MeasureTypeDuration, Value: "s"}},
		"1.2.3s": {
			err: ErrInvalidAmount{Type: MeasureTypeDuration, Value: "1.2.3"},
		},
	}

	for s, tc := range tests {
		t.Run(s, func(t *testing.T) {
			amount, unit, err := amountUnitFromString(s, MeasureTypeDuration, UnitDurationAll[:], UnitDurationUnknown, UnitDuration.String)
			if tc.err != nil {
				if !errors.Is(err, tc.err) {
					t.Error(err)
				}
				return
			}
			if err != nil {
				t.Error(err)
			}
			if amount != tc.amount || unit != tc.unit {
				t.Error(amount, unit)
			}
		})
	}
}

func TestErrorStrings(t *testing.T) {
	tests := map[error]string{
		ErrInvalidUnit{Type: MeasureTypeMass, Value: "100"}:       `invalid mass unit: "100"`,
		ErrInvalidAmount{Type: MeasureTypeDuration, Value: "s"}:   `invalid duration amount: "s"`,
		ErrInvalidAmount{Type: MeasureTypeEnergy, Value: "1.2.3"}: `invalid energy amount: "1.2.3"`,
	}

	for err, s := range tests {
		if err.Error() != s {
			t.Error(err.Error(), s)
		}
	}
}

func TestNewFromStringInvalidUnit(t *testing.T) {
	if _, err := NewMassFromString("100"); err == nil {
		t.Error("mass")
	}
	if _, err := NewVolumeFromString("100"); err == nil {
		t.Error("volume")
	}
	if _, err := NewLengthFromString("100"); err == nil {
		t.Error("length")
	}
	if _, err := NewAreaFromString("100"); err == nil {
		t.Error("area")
	}
	if _, err := NewEnergyFromString("100"); err == nil {
		t.Error("energy")
	}
	if _, err := NewDurationFromString("100"); err == nil {
		t.Error("duration")
	}
	if _, err := NewPowerFromString("100"); err == nil {
		t.Error("power")
	}
	if _, err := NewDataStorageFromString("100"); err == nil {
		t.Error("data storage")
	}
	if _, err := NewFrequencyFromString("100"); err == nil {
		t.Error("frequency")
	}
	if _, err := NewChargeFromString("100"); err == nil {
		t.Error("charge")
	}
	if _, err := NewVoltageFromString("100"); err == nil {
		t.Error("voltage")
	}
}
