package measurement

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidMassAmount = errors.New("invalid mass amount")
	ErrInvalidMassUnit   = errors.New("invalid mass unit")
)

type Mass struct {
	Amount float64  `json:"amount"`
	Unit   UnitMass `json:"unit"`
}

func NewMassFromString(s string) (*Mass, error) {
	var unit UnitMass
	var maxl int
	for _, u := range UnitMassAll {
		if strings.HasSuffix(s, u.String()) && len(u.String()) > maxl {
			unit = u
			maxl = len(u.String())
		}
	}

	if unit == UnitMassUnknown {
		return nil, ErrInvalidMassUnit
	}

	lenAmount := len(s) - len(unit.String())
	if lenAmount <= 0 {
		return nil, ErrInvalidMassAmount
	}

	amount, err := strconv.ParseFloat(s[:lenAmount], 64)
	if err != nil {
		return nil, err
	}

	return &Mass{Amount: amount, Unit: unit}, nil
}

func (s Mass) String() string { return strconv.FormatFloat(s.Amount, 'f', -1, 64) + s.Unit.String() }

func (s *Mass) IsZero() bool {
	if s == nil {
		return true
	}
	if s.Amount == 0 {
		return true
	}
	return false
}

func (s Mass) Convert(unit UnitMass) Mass {
	if v, ok := TryConvertExactMass(s.Amount, s.Unit, unit); ok {
		return Mass{Amount: v, Unit: unit}
	}
	return Mass{Amount: convertMassApproxFromGram(convertMassApproxToGram(s.Amount, s.Unit), unit), Unit: unit}
}

func TryConvertExactMass[T int32 | int64 | float32 | float64](amount T, from, to UnitMass) (v T, ok bool) {
	return convertByLadder(amount, from, to, unitMassLadder)
}

type UnitMass uint8

//go:generate go-enum-encoding -type=UnitMass -string
const (
	UnitMassUnknown UnitMass = iota // json:""
	UnitPicograms                   // json:"pg"
	UnitNanograms                   // json:"ng"
	UnitMicrograms                  // json:"mcg"
	UnitMilligrams                  // json:"mg"
	UnitCentigrams                  // json:"cg"
	UnitDecigrams                   // json:"dg"
	UnitGrams                       // json:"g"
	UnitKilograms                   // json:"kg"
	UnitOunces                      // json:"oz"
	UnitPounds                      // json:"lb"
	UnitStones                      // json:"st"
	UnitMetricTons                  // json:"ton"
	UnitShortTons                   // json:"sst"
	UnitCarats                      // json:"ct"
	UnitOuncesTroy                  // json:"ozt"
	UnitSlugs                       // json:"slug"
)

var UnitMassAll = [...]UnitMass{
	UnitPicograms,
	UnitNanograms,
	UnitMicrograms,
	UnitMilligrams,
	UnitCentigrams,
	UnitDecigrams,
	UnitGrams,
	UnitKilograms,
	UnitOunces,
	UnitPounds,
	UnitStones,
	UnitMetricTons,
	UnitShortTons,
	UnitCarats,
	UnitOuncesTroy,
	UnitSlugs,
}

// this is modified SI conversion ladder
// it utilizes the fact that measurements likely to have at most 3 decimal places or else they can be next unit
// how many m[i] = how many i-1 units needed for this unit
// this constructs 10^3 ladder of unit transformations
var unitMassLadder = ladder[UnitMass]{
	{UnitPicograms, 1},
	{UnitNanograms, 1000},
	{UnitMicrograms, 1000},
	{UnitMilligrams, 1000},
	{UnitCentigrams, 10},
	{UnitDecigrams, 10},
	{UnitCarats, 2},
	{UnitGrams, 5},
	{UnitKilograms, 1000},
	{UnitMetricTons, 1000},
}

// masses are a bit not linear. slugs, ounces troy, and short tones are not in same ladder.
const (
	gramMulApproxUnitOunces     = 28.349523125
	gramMulApproxUnitOuncesTroy = 31.1035
	gramMulApproxUnitPounds     = 453.59237
	gramMulApproxUnitStones     = 6350.29318
	gramMulApproxUnitSlugs      = 14593.9029
	gramMulApproxUnitShortTons  = 907184.74
)

func convertMassApproxToGram[T float32 | float64](amount T, unit UnitMass) T {
	switch unit {
	case UnitOunces:
		return amount * gramMulApproxUnitOunces
	case UnitPounds:
		return amount * gramMulApproxUnitPounds
	case UnitStones:
		return amount * gramMulApproxUnitStones
	case UnitShortTons:
		return amount * gramMulApproxUnitShortTons
	case UnitOuncesTroy:
		return amount * gramMulApproxUnitOuncesTroy
	case UnitSlugs:
		return amount * gramMulApproxUnitSlugs
	default:
		v, _ := convertByLadder(amount, unit, UnitGrams, unitMassLadder)
		return v
	}
}

func convertMassApproxFromGram[T float32 | float64](amount T, unit UnitMass) T {
	switch unit {
	case UnitOunces:
		return amount / gramMulApproxUnitOunces
	case UnitPounds:
		return amount / gramMulApproxUnitPounds
	case UnitStones:
		return amount / gramMulApproxUnitStones
	case UnitShortTons:
		return amount / gramMulApproxUnitShortTons
	case UnitOuncesTroy:
		return amount / gramMulApproxUnitOuncesTroy
	case UnitSlugs:
		return amount / gramMulApproxUnitSlugs
	default:
		v, _ := convertByLadder(amount, UnitGrams, unit, unitMassLadder)
		return v
	}
}
