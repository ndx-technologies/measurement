package measurement

import "strconv"

type Energy struct {
	Amount float64    `json:"amount"`
	Unit   UnitEnergy `json:"unit"`
}

func NewEnergyFromString(s string) (*Energy, error) {
	amount, unit, err := amountUnitFromString(s, MeasureTypeEnergy, UnitEnergyAll[:], UnitEnergyUnknown, UnitEnergy.String)
	if err != nil {
		return nil, err
	}
	return &Energy{Amount: amount, Unit: unit}, nil
}

func (s Energy) String() string { return strconv.FormatFloat(s.Amount, 'f', -1, 64) + s.Unit.String() }

func (s *Energy) IsZero() bool { return s == nil || s.Amount == 0 }

func (s Energy) Convert(unit UnitEnergy) Energy {
	if v, ok := TryConvertExactEnergy(s.Amount, s.Unit, unit); ok {
		return Energy{Amount: v, Unit: unit}
	}
	return Energy{Amount: convertEnergyApproxFromJoule(convertEnergyApproxToJoule(s.Amount, s.Unit), unit), Unit: unit}
}

func TryConvertExactEnergy[T int32 | int64 | float32 | float64](amount T, from, to UnitEnergy) (v T, ok bool) {
	for _, q := range unitEnergyLadders {
		if v, ok := convertByLadder(amount, from, to, q); ok {
			return v, true
		}
	}
	return 0, false
}

type UnitEnergy uint8

//go:generate go-enum-encoding -type=UnitEnergy -string
const (
	UnitEnergyUnknown UnitEnergy = iota // json:""
	UnitJoules                          // json:"J"
	UnitKiloJoules                      // json:"kJ"
	UnitMegaJoules                      // json:"MJ"
	UnitCalories                        // json:"cal"
	UnitKiloCalories                    // json:"kcal"
	UnitWattHours                       // json:"Wh"
	UnitKiloWattHours                   // json:"kWh"
	UnitBTU                             // json:"btu"
)

var UnitEnergyAll = [...]UnitEnergy{
	UnitJoules,
	UnitKiloJoules,
	UnitMegaJoules,
	UnitCalories,
	UnitKiloCalories,
	UnitWattHours,
	UnitKiloWattHours,
	UnitBTU,
}

var unitEnergyJouleLadder = ladder[UnitEnergy]{
	{UnitJoules, 1},
	{UnitKiloJoules, 1000},
	{UnitMegaJoules, 1000},
}

var unitEnergyCalorieLadder = ladder[UnitEnergy]{
	{UnitCalories, 1},
	{UnitKiloCalories, 1000},
}

var unitEnergyWattHourLadder = ladder[UnitEnergy]{
	{UnitWattHours, 1},
	{UnitKiloWattHours, 1000},
}

var unitEnergyLadders = [...]ladder[UnitEnergy]{
	unitEnergyJouleLadder,
	unitEnergyCalorieLadder,
	unitEnergyWattHourLadder,
}

// BTU is not in any ladder.
const (
	jouleMulApproxCalorie  = 4.184
	jouleMulApproxWattHour = 3600
	jouleMulApproxBTU      = 1055.05585262
)

func convertEnergyApproxToJoule(amount float64, unit UnitEnergy) float64 {
	switch unit {
	case UnitKiloJoules:
		return amount * 1000
	case UnitMegaJoules:
		return amount * (1000 * 1000)
	case UnitCalories:
		return amount * jouleMulApproxCalorie
	case UnitKiloCalories:
		return amount * 1000 * jouleMulApproxCalorie
	case UnitWattHours:
		return amount * jouleMulApproxWattHour
	case UnitKiloWattHours:
		return amount * 1000 * jouleMulApproxWattHour
	case UnitBTU:
		return amount * jouleMulApproxBTU
	default:
		return amount
	}
}

func convertEnergyApproxFromJoule(amount float64, unit UnitEnergy) float64 {
	switch unit {
	case UnitKiloJoules:
		return amount / 1000
	case UnitMegaJoules:
		return amount / (1000 * 1000)
	case UnitCalories:
		return amount / jouleMulApproxCalorie
	case UnitKiloCalories:
		return amount / (1000 * jouleMulApproxCalorie)
	case UnitWattHours:
		return amount / jouleMulApproxWattHour
	case UnitKiloWattHours:
		return amount / (1000 * jouleMulApproxWattHour)
	case UnitBTU:
		return amount / jouleMulApproxBTU
	default:
		return amount
	}
}
