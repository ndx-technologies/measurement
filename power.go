package measurement

import "strconv"

type Power struct {
	Amount float64   `json:"amount"`
	Unit   UnitPower `json:"unit"`
}

func NewPowerFromString(s string) (*Power, error) {
	amount, unit, err := amountUnitFromString(s, MeasureTypePower, UnitPowerAll[:], UnitPowerUnknown, UnitPower.String)
	if err != nil {
		return nil, err
	}
	return &Power{Amount: amount, Unit: unit}, nil
}

func (s Power) String() string { return strconv.FormatFloat(s.Amount, 'f', -1, 64) + s.Unit.String() }

func (s *Power) IsZero() bool { return s == nil || s.Amount == 0 }

func (s Power) Convert(unit UnitPower) Power {
	if v, ok := TryConvertExactPower(s.Amount, s.Unit, unit); ok {
		return Power{Amount: v, Unit: unit}
	}
	return Power{Amount: convertPowerApproxFromWatt(convertPowerApproxToWatt(s.Amount, s.Unit), unit), Unit: unit}
}

func TryConvertExactPower[T int32 | int64 | float32 | float64](amount T, from, to UnitPower) (v T, ok bool) {
	return convertByLadder(amount, from, to, unitPowerLadder)
}

// horsepower is not in the ladder.
type UnitPower uint8

//go:generate go-enum-encoding -type=UnitPower -string
const (
	UnitPowerUnknown UnitPower = iota // json:""
	UnitWatts                         // json:"W"
	UnitKiloWatts                     // json:"kW"
	UnitMegaWatts                     // json:"MW"
	UnitGigaWatts                     // json:"GW"
	UnitHorsepower                    // json:"hp"
)

var UnitPowerAll = [...]UnitPower{
	UnitWatts,
	UnitKiloWatts,
	UnitMegaWatts,
	UnitGigaWatts,
	UnitHorsepower,
}

var unitPowerLadder = ladder[UnitPower]{
	{UnitWatts, 1},
	{UnitKiloWatts, 1000},
	{UnitMegaWatts, 1000},
	{UnitGigaWatts, 1000},
}

const wattMulApproxHorsepower = 745.6998715822702

func convertPowerApproxToWatt(amount float64, unit UnitPower) float64 {
	switch unit {
	case UnitKiloWatts:
		return amount * 1000
	case UnitMegaWatts:
		return amount * (1000 * 1000)
	case UnitGigaWatts:
		return amount * (1000 * 1000 * 1000)
	case UnitHorsepower:
		return amount * wattMulApproxHorsepower
	default:
		return amount
	}
}

func convertPowerApproxFromWatt(amount float64, unit UnitPower) float64 {
	switch unit {
	case UnitKiloWatts:
		return amount / 1000
	case UnitMegaWatts:
		return amount / (1000 * 1000)
	case UnitGigaWatts:
		return amount / (1000 * 1000 * 1000)
	case UnitHorsepower:
		return amount / wattMulApproxHorsepower
	default:
		return amount
	}
}
