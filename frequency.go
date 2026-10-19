package measurement

import "strconv"

type Frequency struct {
	Amount float64       `json:"amount"`
	Unit   UnitFrequency `json:"unit"`
}

func NewFrequencyFromString(s string) (*Frequency, error) {
	amount, unit, err := amountUnitFromString(s, MeasureTypeFrequency, UnitFrequencyAll[:], UnitFrequencyUnknown, UnitFrequency.String)
	if err != nil {
		return nil, err
	}
	return &Frequency{Amount: amount, Unit: unit}, nil
}

func (s Frequency) String() string {
	return strconv.FormatFloat(s.Amount, 'f', -1, 64) + s.Unit.String()
}

func (s *Frequency) IsZero() bool { return s == nil || s.Amount == 0 }

func (s Frequency) Convert(unit UnitFrequency) Frequency {
	if v, ok := TryConvertExactFrequency(s.Amount, s.Unit, unit); ok {
		return Frequency{Amount: v, Unit: unit}
	}
	return Frequency{Amount: s.Amount, Unit: unit}
}

func TryConvertExactFrequency[T int32 | int64 | float32 | float64](amount T, from, to UnitFrequency) (v T, ok bool) {
	return convertByLadder(amount, from, to, unitFrequencyLadder)
}

type UnitFrequency uint8

//go:generate go-enum-encoding -type=UnitFrequency -string
const (
	UnitFrequencyUnknown UnitFrequency = iota // json:""
	UnitHertz                                 // json:"Hz"
	UnitKiloHertz                             // json:"kHz"
	UnitMegaHertz                             // json:"MHz"
	UnitGigaHertz                             // json:"GHz"
	UnitTeraHertz                             // json:"THz"
)

var UnitFrequencyAll = [...]UnitFrequency{
	UnitHertz,
	UnitKiloHertz,
	UnitMegaHertz,
	UnitGigaHertz,
	UnitTeraHertz,
}

var unitFrequencyLadder = ladder[UnitFrequency]{
	{UnitHertz, 1},
	{UnitKiloHertz, 1000},
	{UnitMegaHertz, 1000},
	{UnitGigaHertz, 1000},
	{UnitTeraHertz, 1000},
}
