package measurement

import "strconv"

type Charge struct {
	Amount float64    `json:"amount"`
	Unit   UnitCharge `json:"unit"`
}

func NewChargeFromString(s string) (*Charge, error) {
	amount, unit, err := amountUnitFromString(s, MeasureTypeCharge, UnitChargeAll[:], UnitChargeUnknown, UnitCharge.String)
	if err != nil {
		return nil, err
	}
	return &Charge{Amount: amount, Unit: unit}, nil
}

func (s Charge) String() string { return strconv.FormatFloat(s.Amount, 'f', -1, 64) + s.Unit.String() }

func (s *Charge) IsZero() bool { return s == nil || s.Amount == 0 }

func (s Charge) Convert(unit UnitCharge) Charge {
	if v, ok := TryConvertExactCharge(s.Amount, s.Unit, unit); ok {
		return Charge{Amount: v, Unit: unit}
	}
	return Charge{Amount: convertChargeApproxFromMilliAmpereHour(convertChargeApproxToMilliAmpereHour(s.Amount, s.Unit), unit), Unit: unit}
}

func TryConvertExactCharge[T int32 | int64 | float32 | float64](amount T, from, to UnitCharge) (v T, ok bool) {
	return convertByLadder(amount, from, to, unitChargeLadder)
}

// coulombs are not in the ladder.
type UnitCharge uint8

//go:generate go-enum-encoding -type=UnitCharge -string
const (
	UnitChargeUnknown    UnitCharge = iota // json:""
	UnitMilliAmpereHours                   // json:"mAh"
	UnitAmpereHours                        // json:"Ah"
	UnitKiloAmpereHours                    // json:"kAh"
	UnitCoulombs                           // json:"C"
)

var UnitChargeAll = [...]UnitCharge{
	UnitMilliAmpereHours,
	UnitAmpereHours,
	UnitKiloAmpereHours,
	UnitCoulombs,
}

var unitChargeLadder = ladder[UnitCharge]{
	{UnitMilliAmpereHours, 1},
	{UnitAmpereHours, 1000},
	{UnitKiloAmpereHours, 1000},
}

const milliAmpereHourMulApproxCoulomb = 3.6

func convertChargeApproxToMilliAmpereHour(amount float64, unit UnitCharge) float64 {
	switch unit {
	case UnitAmpereHours:
		return amount * 1000
	case UnitKiloAmpereHours:
		return amount * (1000 * 1000)
	case UnitCoulombs:
		return amount / milliAmpereHourMulApproxCoulomb
	default:
		return amount
	}
}

func convertChargeApproxFromMilliAmpereHour(amount float64, unit UnitCharge) float64 {
	switch unit {
	case UnitAmpereHours:
		return amount / 1000
	case UnitKiloAmpereHours:
		return amount / (1000 * 1000)
	case UnitCoulombs:
		return amount * milliAmpereHourMulApproxCoulomb
	default:
		return amount
	}
}
