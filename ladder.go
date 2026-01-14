package measurement

type ladderItem[U comparable] struct {
	unit     U
	fromPrev int
}

// ladder is a sequence of units with conversion factors between adjacent units.
// Unit systems typically define units in a sequence with whole multipliers between units.
type ladder[U comparable] []ladderItem[U]

func (l ladder[U]) indexOf(unit U) int {
	for i, item := range l {
		if item.unit == unit {
			return i
		}
	}
	return -1
}

func convertByLadder[U comparable, T int32 | int64 | float32 | float64](amount T, from, to U, ladder ladder[U]) (T, bool) {
	if from == to || amount == 0 {
		return amount, true
	}

	idxFrom, idxTo := ladder.indexOf(from), ladder.indexOf(to)
	if idxFrom == -1 || idxTo == -1 || idxFrom == idxTo {
		return amount, false
	}

	var f int = 1

	for idx := idxFrom; (idx + 1) <= idxTo; idx++ {
		f *= ladder[idx+1].fromPrev
	}

	for idx := idxFrom; idx > idxTo; idx-- {
		f *= ladder[idx].fromPrev
	}

	ft := T(f)

	if idxFrom < idxTo {
		if (amount/ft)*ft != amount {
			return 0, false // loss of precision without fractions
		}
		amount /= ft
	} else {
		amount *= ft
	}

	return amount, true
}
