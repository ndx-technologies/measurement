package measurement

import "strconv"

type Area struct {
	Amount float64  `json:"amount"`
	Unit   UnitArea `json:"unit"`
}

func NewAreaFromString(s string) (*Area, error) {
	amount, unit, err := amountUnitFromString(s, MeasureTypeArea, UnitAreaAll[:], UnitAreaUnknown, UnitArea.String)
	if err != nil {
		return nil, err
	}
	return &Area{Amount: amount, Unit: unit}, nil
}

func (s Area) String() string { return strconv.FormatFloat(s.Amount, 'f', -1, 64) + s.Unit.String() }

func (s *Area) IsZero() bool { return s == nil || s.Amount == 0 }

func (s Area) Convert(unit UnitArea) Area {
	if v, ok := TryConvertExactArea(s.Amount, s.Unit, unit); ok {
		return Area{Amount: v, Unit: unit}
	}
	return Area{Amount: convertAreaApproxFromSquareMeter(convertAreaApproxToSquareMeter(s.Amount, s.Unit), unit), Unit: unit}
}

func TryConvertExactArea[T int32 | int64 | float32 | float64](amount T, from, to UnitArea) (v T, ok bool) {
	for _, q := range unitAreaLadders {
		if v, ok := convertByLadder(amount, from, to, q); ok {
			return v, true
		}
	}
	return 0, false
}

type UnitArea uint8

//go:generate go-enum-encoding -type=UnitArea -string
const (
	UnitAreaUnknown       UnitArea = iota // json:""
	UnitSquareMilliMeters                 // json:"mm2"
	UnitSquareCentiMeters                 // json:"cm2"
	UnitSquareMeters                      // json:"m2"
	UnitHectares                          // json:"ha"
	UnitSquareKiloMeters                  // json:"km2"
	UnitSquareInches                      // json:"in2"
	UnitSquareFeet                        // json:"ft2"
	UnitSquareYards                       // json:"yd2"
	UnitAcres                             // json:"acre"
	UnitSquareMiles                       // json:"mi2"
)

var UnitAreaAll = [...]UnitArea{
	UnitSquareMilliMeters,
	UnitSquareCentiMeters,
	UnitSquareMeters,
	UnitHectares,
	UnitSquareKiloMeters,
	UnitSquareInches,
	UnitSquareFeet,
	UnitSquareYards,
	UnitAcres,
	UnitSquareMiles,
}

var unitAreaMetricLadder = ladder[UnitArea]{
	{UnitSquareMilliMeters, 1},
	{UnitSquareCentiMeters, 100},
	{UnitSquareMeters, 100 * 100},
	{UnitHectares, 100 * 100},
	{UnitSquareKiloMeters, 100},
}

var unitAreaImperialLadder = ladder[UnitArea]{
	{UnitSquareInches, 1},
	{UnitSquareFeet, 12 * 12},
	{UnitSquareYards, 3 * 3},
	{UnitAcres, 4840},
	{UnitSquareMiles, 640},
}

var unitAreaLadders = [...]ladder[UnitArea]{
	unitAreaMetricLadder,
	unitAreaImperialLadder,
}

const (
	squareMeterMulApproxSquareInches = 0.00064516
	squareMeterMulApproxSquareFeet   = 0.09290304
	squareMeterMulApproxSquareYards  = 0.83612736
	squareMeterMulApproxAcres        = 4046.8564224
	squareMeterMulApproxSquareMiles  = 2589988.110336
)

func convertAreaApproxToSquareMeter(amount float64, unit UnitArea) float64 {
	switch unit {
	case UnitSquareMilliMeters:
		return amount / (1000 * 1000)
	case UnitSquareCentiMeters:
		return amount / (100 * 100)
	case UnitHectares:
		return amount * (100 * 100)
	case UnitSquareKiloMeters:
		return amount * (1000 * 1000)
	case UnitSquareInches:
		return amount * squareMeterMulApproxSquareInches
	case UnitSquareFeet:
		return amount * squareMeterMulApproxSquareFeet
	case UnitSquareYards:
		return amount * squareMeterMulApproxSquareYards
	case UnitAcres:
		return amount * squareMeterMulApproxAcres
	case UnitSquareMiles:
		return amount * squareMeterMulApproxSquareMiles
	default:
		return amount
	}
}

func convertAreaApproxFromSquareMeter(amount float64, unit UnitArea) float64 {
	switch unit {
	case UnitSquareMilliMeters:
		return amount * (1000 * 1000)
	case UnitSquareCentiMeters:
		return amount * (100 * 100)
	case UnitHectares:
		return amount / (100 * 100)
	case UnitSquareKiloMeters:
		return amount / (1000 * 1000)
	case UnitSquareInches:
		return amount / squareMeterMulApproxSquareInches
	case UnitSquareFeet:
		return amount / squareMeterMulApproxSquareFeet
	case UnitSquareYards:
		return amount / squareMeterMulApproxSquareYards
	case UnitAcres:
		return amount / squareMeterMulApproxAcres
	case UnitSquareMiles:
		return amount / squareMeterMulApproxSquareMiles
	default:
		return amount
	}
}
