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

// TryConvertExactMass using only integer factors
func TryConvertExactMass[T int32 | int64 | float32 | float64](amount T, from, to UnitMass) (v T, ok bool) {
	if from == to || to == UnitMassUnknown || amount == 0 {
		return amount, true
	}
	idxFrom, idxTo, ladder := tryGetSameMassLadder(from, to)
	if ladder == nil || idxFrom == -1 || idxTo == -1 {
		return 0, false
	}
	return convertMassExact(amount, idxFrom, idxTo, *ladder), true
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

type unitMassLadderItem struct {
	unit     UnitMass
	fromPrev int64
}

type unitMassLadderType []unitMassLadderItem

// this is modified SI conversion ladder
// it utilizes the fact that measurements likely to have at most 3 decimal places or else they can be next unit
// how many m[i] = how many i-1 units needed for this unit
// this constructs 10^3 ladder of unit transformations
var unitMassLadder = [...]unitMassLadderItem{
	{unit: UnitPicograms, fromPrev: 1},
	{unit: UnitNanograms, fromPrev: 1000},
	{unit: UnitMicrograms, fromPrev: 1000},
	{unit: UnitMilligrams, fromPrev: 1000},
	{unit: UnitCentigrams, fromPrev: 10},
	{unit: UnitDecigrams, fromPrev: 10},
	{unit: UnitCarats, fromPrev: 2},
	{unit: UnitGrams, fromPrev: 5},
	{unit: UnitKilograms, fromPrev: 1000},
	{unit: UnitMetricTons, fromPrev: 1000},
}

func idxUnitMassInLadder(unit UnitMass) int {
	for i, u := range unitMassLadder {
		if u.unit == unit {
			return i
		}
	}
	return -1
}

func factorFromMassLadder[T int64 | float64](from, to int) T {
	var f T = 1
	for idx := from; (idx + 1) <= to; idx++ {
		f *= T(unitMassLadder[idx+1].fromPrev)
	}
	for idx := from; idx > to; idx-- {
		f *= T(unitMassLadder[idx].fromPrev)
	}
	return f
}

func convertMassApproxInLadder(amount float64, from, to UnitMass) float64 {
	if from == to {
		return amount
	}

	idxFrom, idxTo := idxUnitMassInLadder(from), idxUnitMassInLadder(to)
	if idxFrom == -1 || idxTo == -1 || idxFrom == idxTo {
		return amount
	}

	f := factorFromMassLadder[float64](idxFrom, idxTo)

	if idxFrom < idxTo {
		amount /= f
	} else {
		amount *= f
	}

	return amount
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

func convertMassApproxToGram(amount float64, unit UnitMass) float64 {
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
		return convertMassApproxInLadder(amount, unit, UnitGrams)
	}
}

func convertMassApproxFromGram(amount float64, unit UnitMass) float64 {
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
		return convertMassApproxInLadder(amount, UnitGrams, unit)
	}
}

func tryGetSameMassLadder(a, b UnitMass) (idxA, idxB int, ladder *unitMassLadderType) {
	idxA = idxUnitMassInLadder(a)
	idxB = idxUnitMassInLadder(b)
	if idxA != -1 && idxB != -1 {
		l := unitMassLadderType(unitMassLadder[:])
		return idxA, idxB, &l
	}
	return -1, -1, nil
}

func convertMassExact[T int32 | int64 | float32 | float64](amount T, idxFrom, idxTo int, ladder unitMassLadderType) T {
	if idxFrom == idxTo {
		return amount
	}

	var f T = 1

	for idx := idxFrom; (idx + 1) <= idxTo; idx++ {
		f *= T(ladder[idx+1].fromPrev)
	}

	for idx := idxFrom; idx > idxTo; idx-- {
		f *= T(ladder[idx].fromPrev)
	}

	if idxFrom < idxTo {
		amount /= f
	} else {
		amount *= f
	}

	return amount
}
