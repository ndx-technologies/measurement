package measurement

import "strconv"

type Duration struct {
	Amount float64      `json:"amount"`
	Unit   UnitDuration `json:"unit"`
}

func NewDurationFromString(s string) (*Duration, error) {
	amount, unit, err := amountUnitFromString(s, MeasureTypeDuration, UnitDurationAll[:], UnitDurationUnknown, UnitDuration.String)
	if err != nil {
		return nil, err
	}
	return &Duration{Amount: amount, Unit: unit}, nil
}

func (s Duration) String() string {
	return strconv.FormatFloat(s.Amount, 'f', -1, 64) + s.Unit.String()
}

func (s *Duration) IsZero() bool { return s == nil || s.Amount == 0 }

func (s Duration) Convert(unit UnitDuration) Duration {
	if v, ok := TryConvertExactDuration(s.Amount, s.Unit, unit); ok {
		return Duration{Amount: v, Unit: unit}
	}
	if unit == UnitDurationUnknown {
		return s
	}
	return Duration{Amount: convertDurationApproxFromDay(convertDurationApproxToDay(s.Amount, s.Unit), unit), Unit: unit}
}

func TryConvertExactDuration[T int32 | int64 | float32 | float64](amount T, from, to UnitDuration) (v T, ok bool) {
	for _, q := range unitDurationLadders {
		if v, ok := convertByLadder(amount, from, to, q); ok {
			return v, true
		}
	}
	return 0, false
}

type UnitDuration uint8

//go:generate go-enum-encoding -type=UnitDuration -string
const (
	UnitDurationUnknown UnitDuration = iota // json:""
	UnitNanoSeconds                         // json:"ns"
	UnitMicroSeconds                        // json:"µs"
	UnitMilliSeconds                        // json:"ms"
	UnitSeconds                             // json:"s"
	UnitMinutes                             // json:"min"
	UnitHours                               // json:"h"
	UnitDays                                // json:"d"
	UnitWeeks                               // json:"week"
	UnitMonths                              // json:"month"
	UnitYears                               // json:"y"
)

var UnitDurationAll = [...]UnitDuration{
	UnitNanoSeconds,
	UnitMicroSeconds,
	UnitMilliSeconds,
	UnitSeconds,
	UnitMinutes,
	UnitHours,
	UnitDays,
	UnitWeeks,
	UnitMonths,
	UnitYears,
}

var unitDurationLadder = ladder[UnitDuration]{
	{UnitNanoSeconds, 1},
	{UnitMicroSeconds, 1000},
	{UnitMilliSeconds, 1000},
	{UnitSeconds, 1000},
	{UnitMinutes, 60},
	{UnitHours, 60},
	{UnitDays, 24},
	{UnitWeeks, 7},
}

// months and years form their own ladder, they do not have a whole number of days.
var unitDurationCalendarLadder = ladder[UnitDuration]{
	{UnitMonths, 1},
	{UnitYears, 12},
}

var unitDurationLadders = [...]ladder[UnitDuration]{
	unitDurationLadder,
	unitDurationCalendarLadder,
}

const (
	dayMulApproxMonth = 30.436875 // average Gregorian month, 365.2425 / 12
	dayMulApproxYear  = 365.2425
)

func convertDurationApproxToDay(amount float64, unit UnitDuration) float64 {
	switch unit {
	case UnitNanoSeconds:
		return amount / (1000 * 1000 * 1000 * 60 * 60 * 24)
	case UnitMicroSeconds:
		return amount / (1000 * 1000 * 60 * 60 * 24)
	case UnitMilliSeconds:
		return amount / (1000 * 60 * 60 * 24)
	case UnitSeconds:
		return amount / (60 * 60 * 24)
	case UnitMinutes:
		return amount / (60 * 24)
	case UnitHours:
		return amount / 24
	case UnitWeeks:
		return amount * 7
	case UnitMonths:
		return amount * dayMulApproxMonth
	case UnitYears:
		return amount * dayMulApproxYear
	default:
		return amount
	}
}

func convertDurationApproxFromDay(amount float64, unit UnitDuration) float64 {
	switch unit {
	case UnitNanoSeconds:
		return amount * (1000 * 1000 * 1000 * 60 * 60 * 24)
	case UnitMicroSeconds:
		return amount * (1000 * 1000 * 60 * 60 * 24)
	case UnitMilliSeconds:
		return amount * (1000 * 60 * 60 * 24)
	case UnitSeconds:
		return amount * (60 * 60 * 24)
	case UnitMinutes:
		return amount * (60 * 24)
	case UnitHours:
		return amount * 24
	case UnitWeeks:
		return amount / 7
	case UnitMonths:
		return amount / dayMulApproxMonth
	case UnitYears:
		return amount / dayMulApproxYear
	default:
		return amount
	}
}
