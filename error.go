package measurement

import "strconv"

// ErrInvalidUnit is returned when the trailing unit of a measurement string is not known.
type ErrInvalidUnit struct {
	Type  MeasureType
	Value string
}

func (e ErrInvalidUnit) Error() string {
	return "invalid " + e.Type.String() + " unit: " + strconv.Quote(e.Value)
}

// ErrInvalidAmount is returned when the amount in front of the unit is not a positive number.
type ErrInvalidAmount struct {
	Type  MeasureType
	Value string
}

func (e ErrInvalidAmount) Error() string {
	return "invalid " + e.Type.String() + " amount: " + strconv.Quote(e.Value)
}
