package measurement

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidVolumeAmount = errors.New("invalid volume amount")
	ErrInvalidVolumeUnit   = errors.New("invalid volume unit")
)

type Volume struct {
	Amount float64    `json:"amount"`
	Unit   UnitVolume `json:"unit"`
}

func NewVolumeFromString(s string) (*Volume, error) {
	var unit UnitVolume
	var maxl int
	for _, u := range UnitVolumeAll {
		if strings.HasSuffix(s, u.String()) && len(u.String()) > maxl {
			unit = u
			maxl = len(u.String())
		}
	}

	if unit == UnitVolumeUnknown {
		return nil, ErrInvalidVolumeUnit
	}

	lenAmount := len(s) - len(unit.String())
	if lenAmount <= 0 {
		return nil, ErrInvalidVolumeAmount
	}

	amount, err := strconv.ParseFloat(s[:lenAmount], 64)
	if err != nil {
		return nil, err
	}

	return &Volume{Amount: amount, Unit: unit}, nil
}

func (s Volume) String() string { return strconv.FormatFloat(s.Amount, 'f', -1, 64) + s.Unit.String() }

func (s *Volume) IsZero() bool {
	if s == nil {
		return true
	}
	if s.Amount == 0 {
		return true
	}
	return false
}

func (s Volume) Convert(unit UnitVolume) Volume {
	if v, ok := TryConvertExactVolume(s.Amount, s.Unit, unit); ok {
		return Volume{Amount: v, Unit: unit}
	}

	litters, factor := convertVolumeToLiters(s.Amount, s.Unit)
	amount, factor := convertVolumeFromLiters(litters*factor, unit)
	return Volume{Amount: amount / factor, Unit: unit}
}

// TryConvertExactVolume using only integer factors
func TryConvertExactVolume[T int32 | int64 | float32 | float64](amount T, from, to UnitVolume) (v T, ok bool) {
	for _, q := range unitVolumeLadders {
		if v, ok := convertByLadder(amount, from, to, q.ladder); ok {
			return v, true
		}
	}
	return 0, false
}

// skipping `metric cup` and `acre-feet`, they are not in any ladder.

type UnitVolume uint8

//go:generate go-enum-encoding -type=UnitVolume -string
const (
	UnitVolumeUnknown       UnitVolume = iota // json:""
	UnitMilliLiters                           // json:"ml"
	UnitCentiLiters                           // json:"cl"
	UnitDeciLiters                            // json:"dl"
	UnitLiters                                // json:"l"
	UnitKiloLiters                            // json:"kl"
	UnitMegaLiters                            // json:"Ml"
	UnitCubicMilliMeters                      // json:"mm3"
	UnitCubicCentiMeters                      // json:"cm3"
	UnitCubicDeciMeters                       // json:"dm3"
	UnitCubicFeet                             // json:"ft3"
	UnitCubicInches                           // json:"in3"
	UnitCubicMeters                           // json:"m3"
	UnitCubicKiloMeters                       // json:"km3"
	UnitCubicMiles                            // json:"mi3"
	UnitCubicYards                            // json:"yd3"
	UnitBushels                               // json:"bu"
	UnitCups                                  // json:"cup"
	UnitFluidOunces                           // json:"floz"
	UnitGallons                               // json:"gal"
	UnitPints                                 // json:"pt"
	UnitQuarts                                // json:"qt"
	UnitTablespoons                           // json:"tbsp"
	UnitTeaspoons                             // json:"tsp"
	UnitImperialFluidOunces                   // json:"impfloz"
	UnitImperialGallons                       // json:"impgal"
	UnitImperialGills                         // json:"impgil"
	UnitImperialPints                         // json:"imppt"
	UnitImperialQuarts                        // json:"impqt"
	UnitImperialTablespoons                   // json:"imptbsp"
	UnitImperialTeaspoons                     // json:"imptsp"
)

var UnitVolumeAll = [...]UnitVolume{
	UnitMilliLiters,
	UnitCentiLiters,
	UnitDeciLiters,
	UnitLiters,
	UnitKiloLiters,
	UnitMegaLiters,
	UnitCubicMilliMeters,
	UnitCubicCentiMeters,
	UnitCubicDeciMeters,
	UnitCubicFeet,
	UnitCubicInches,
	UnitCubicMeters,
	UnitCubicKiloMeters,
	UnitCubicMiles,
	UnitCubicYards,
	UnitBushels,
	UnitCups,
	UnitFluidOunces,
	UnitGallons,
	UnitPints,
	UnitQuarts,
	UnitTablespoons,
	UnitTeaspoons,
	UnitImperialFluidOunces,
	UnitImperialGallons,
	UnitImperialPints,
	UnitImperialQuarts,
	UnitImperialTablespoons,
	UnitImperialTeaspoons,
}

var unitVolumeLiterLadder = ladder[UnitVolume]{
	{UnitMilliLiters, 1},
	{UnitCentiLiters, 10},
	{UnitDeciLiters, 10},
	{UnitLiters, 10},
	{UnitKiloLiters, 1000},
	{UnitMegaLiters, 1000},
}

var unitVolumeMeterLadder = ladder[UnitVolume]{
	{UnitCubicMilliMeters, 1},
	{UnitCubicCentiMeters, 10 * 10 * 10},
	{UnitCubicDeciMeters, 10 * 10 * 10},
	{UnitCubicMeters, 10 * 10 * 10},
	{UnitCubicKiloMeters, 1000 * 1000 * 1000},
}

var unitVolumeInchLadder = ladder[UnitVolume]{
	{UnitCubicInches, 1},
	{UnitCubicFeet, 12 * 12 * 12},
	{UnitCubicYards, 3 * 3 * 3},
	{UnitCubicMiles, 1760 * 1760 * 1760},
}

var unitVolumeImperialLadder = ladder[UnitVolume]{
	{UnitImperialTeaspoons, 1},
	{UnitImperialTablespoons, 3},
	{UnitImperialFluidOunces, 2},
	{UnitImperialGills, 5},
	{UnitImperialPints, 4},
	{UnitImperialQuarts, 2},
	{UnitImperialGallons, 4},
	{UnitBushels, 8},
}

var unitVolumeUSALadder = ladder[UnitVolume]{
	{UnitTeaspoons, 1},
	{UnitTablespoons, 3},
	{UnitFluidOunces, 2},
	{UnitCups, 8},
	{UnitPints, 2},
	{UnitQuarts, 2},
	{UnitGallons, 4},
}

const (
	litersMulApproxCubicFeet = 28.3168
	litersMulImperialPint    = 0.568261
	litersMulPints           = 0.473176
)

var unitVolumeLadders = [...]struct {
	ladder ladder[UnitVolume]
	factor float64
	target UnitVolume
}{
	{ladder: unitVolumeMeterLadder, target: UnitCubicDeciMeters, factor: 1},
	{ladder: unitVolumeInchLadder, target: UnitCubicFeet, factor: litersMulApproxCubicFeet},
	{ladder: unitVolumeImperialLadder, target: UnitImperialPints, factor: litersMulImperialPint},
	{ladder: unitVolumeUSALadder, target: UnitPints, factor: litersMulPints},
	{ladder: unitVolumeLiterLadder, target: UnitLiters, factor: 1},
}

// move from ladder to conversion point, and convert.
// if factor is 1, no conversion to float is needed, otherwise cast to float and multiply.
func convertVolumeToLiters[T int64 | float32 | float64](amount T, unit UnitVolume) (T, float64) {
	for _, q := range unitVolumeLadders {
		if v, ok := TryConvertExactVolume(amount, unit, q.target); ok {
			return v, q.factor
		}
	}
	return amount, 1
}

// find which ladder we need to use and calculate movement of conversion point to target.
// if factor is not 1, then cast to float and divide.
func convertVolumeFromLiters[T int64 | float32 | float64](amount T, unit UnitVolume) (T, float64) {
	for _, q := range unitVolumeLadders {
		if v, ok := TryConvertExactVolume(amount, q.target, unit); ok {
			return v, q.factor
		}
	}
	return amount, 1
}
