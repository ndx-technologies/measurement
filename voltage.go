package measurement

import "strconv"

type Voltage struct {
	Amount float64     `json:"amount"`
	Unit   UnitVoltage `json:"unit"`
}

func NewVoltageFromString(s string) (*Voltage, error) {
	amount, unit, err := amountUnitFromString(s, MeasureTypeVoltage, UnitVoltageAll[:], UnitVoltageUnknown, UnitVoltage.String)
	if err != nil {
		return nil, err
	}
	return &Voltage{Amount: amount, Unit: unit}, nil
}

func (s Voltage) String() string { return strconv.FormatFloat(s.Amount, 'f', -1, 64) + s.Unit.String() }

func (s *Voltage) IsZero() bool { return s == nil || s.Amount == 0 }

func (s Voltage) Convert(unit UnitVoltage) Voltage {
	if v, ok := TryConvertExactVoltage(s.Amount, s.Unit, unit); ok {
		return Voltage{Amount: v, Unit: unit}
	}
	return Voltage{Amount: s.Amount, Unit: unit}
}

func TryConvertExactVoltage[T int32 | int64 | float32 | float64](amount T, from, to UnitVoltage) (v T, ok bool) {
	return convertByLadder(amount, from, to, unitVoltageLadder)
}

type UnitVoltage uint8

//go:generate go-enum-encoding -type=UnitVoltage -string
const (
	UnitVoltageUnknown UnitVoltage = iota // json:""
	UnitMilliVolts                        // json:"mV"
	UnitVolts                             // json:"V"
	UnitKiloVolts                         // json:"kV"
)

var UnitVoltageAll = [...]UnitVoltage{
	UnitMilliVolts,
	UnitVolts,
	UnitKiloVolts,
}

var unitVoltageLadder = ladder[UnitVoltage]{
	{UnitMilliVolts, 1},
	{UnitVolts, 1000},
	{UnitKiloVolts, 1000},
}
