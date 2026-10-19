package measurement

import (
	"strconv"
	"strings"
)

// amountUnitFromString splits s into the trailing unit and the amount before it.
// The longest matching unit wins, so "100kcal" is not read as "100k" + "cal".
func amountUnitFromString[U comparable](s string, typ MeasureType, units []U, unknown U, unitText func(U) string) (float64, U, error) {
	var unit U
	var maxlen int
	for _, u := range units {
		t := unitText(u)
		if strings.HasSuffix(s, t) && len(t) > maxlen {
			unit = u
			maxlen = len(t)
		}
	}

	if unit == unknown {
		return 0, unknown, ErrInvalidUnit{Type: typ, Value: s}
	}

	amountLen := len(s) - maxlen
	if amountLen <= 0 {
		return 0, unknown, ErrInvalidAmount{Type: typ, Value: s}
	}

	amount, err := strconv.ParseFloat(s[:amountLen], 64)
	if err != nil {
		return 0, unknown, ErrInvalidAmount{Type: typ, Value: s[:amountLen]}
	}
	return amount, unit, nil
}
