package measurement

import "strconv"

type Length struct {
	Amount float64    `json:"amount"`
	Unit   UnitLength `json:"unit"`
}

func NewLengthFromString(s string) (*Length, error) {
	amount, unit, err := amountUnitFromString(s, MeasureTypeLength, UnitLengthAll[:], UnitLengthUnknown, UnitLength.String)
	if err != nil {
		return nil, err
	}
	return &Length{Amount: amount, Unit: unit}, nil
}

func (s Length) String() string { return strconv.FormatFloat(s.Amount, 'f', -1, 64) + s.Unit.String() }

func (s *Length) IsZero() bool { return s == nil || s.Amount == 0 }

func (s Length) Convert(unit UnitLength) Length {
	if v, ok := TryConvertExactLength(s.Amount, s.Unit, unit); ok {
		return Length{Amount: v, Unit: unit}
	}
	return Length{Amount: convertLengthApproxFromMeter(convertLengthApproxToMeter(s.Amount, s.Unit), unit), Unit: unit}
}

func TryConvertExactLength[T int32 | int64 | float32 | float64](amount T, from, to UnitLength) (v T, ok bool) {
	for _, q := range unitLengthLadders {
		if v, ok := convertByLadder(amount, from, to, q); ok {
			return v, true
		}
	}
	return 0, false
}

// skipping nautical miles, they are not in any ladder.
type UnitLength uint8

//go:generate go-enum-encoding -type=UnitLength -string
const (
	UnitLengthUnknown UnitLength = iota // json:""
	UnitMilliMeters                     // json:"mm"
	UnitCentiMeters                     // json:"cm"
	UnitDeciMeters                      // json:"dm"
	UnitMeters                          // json:"m"
	UnitKiloMeters                      // json:"km"
	UnitInches                          // json:"in"
	UnitFeet                            // json:"ft"
	UnitYards                           // json:"yd"
	UnitMiles                           // json:"mi"
	UnitNauticalMiles                   // json:"nmi"
)

var UnitLengthAll = [...]UnitLength{
	UnitMilliMeters,
	UnitCentiMeters,
	UnitDeciMeters,
	UnitMeters,
	UnitKiloMeters,
	UnitInches,
	UnitFeet,
	UnitYards,
	UnitMiles,
	UnitNauticalMiles,
}

var unitLengthMetricLadder = ladder[UnitLength]{
	{UnitMilliMeters, 1},
	{UnitCentiMeters, 10},
	{UnitDeciMeters, 10},
	{UnitMeters, 10},
	{UnitKiloMeters, 1000},
}

var unitLengthImperialLadder = ladder[UnitLength]{
	{UnitInches, 1},
	{UnitFeet, 12},
	{UnitYards, 3},
	{UnitMiles, 1760},
}

var unitLengthLadders = [...]ladder[UnitLength]{
	unitLengthMetricLadder,
	unitLengthImperialLadder,
}

// nautical miles are not in the imperial ladder.
const (
	meterMulApproxInches        = 0.0254
	meterMulApproxFeet          = 0.3048
	meterMulApproxYards         = 0.9144
	meterMulApproxMiles         = 1609.344
	meterMulApproxNauticalMiles = 1852
)

func convertLengthApproxToMeter(amount float64, unit UnitLength) float64 {
	switch unit {
	case UnitMilliMeters:
		return amount / 1000
	case UnitCentiMeters:
		return amount / 100
	case UnitDeciMeters:
		return amount / 10
	case UnitKiloMeters:
		return amount * 1000
	case UnitInches:
		return amount * meterMulApproxInches
	case UnitFeet:
		return amount * meterMulApproxFeet
	case UnitYards:
		return amount * meterMulApproxYards
	case UnitMiles:
		return amount * meterMulApproxMiles
	case UnitNauticalMiles:
		return amount * meterMulApproxNauticalMiles
	default:
		return amount
	}
}

func convertLengthApproxFromMeter(amount float64, unit UnitLength) float64 {
	switch unit {
	case UnitMilliMeters:
		return amount * 1000
	case UnitCentiMeters:
		return amount * 100
	case UnitDeciMeters:
		return amount * 10
	case UnitKiloMeters:
		return amount / 1000
	case UnitInches:
		return amount / meterMulApproxInches
	case UnitFeet:
		return amount / meterMulApproxFeet
	case UnitYards:
		return amount / meterMulApproxYards
	case UnitMiles:
		return amount / meterMulApproxMiles
	case UnitNauticalMiles:
		return amount / meterMulApproxNauticalMiles
	default:
		return amount
	}
}
