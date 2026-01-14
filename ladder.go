package measurement

type ladderItem[U comparable] struct {
	unit     U
	fromPrev int64
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

	var f T = 1

	for idx := idxFrom; (idx + 1) <= idxTo; idx++ {
		f *= T(ladder[idx+1].fromPrev)
	}

	for idx := idxFrom; idx > idxTo; idx-- {
		f *= T(ladder[idx].fromPrev)
	}

	if idxFrom < idxTo {
		if (amount/f)*f != amount {
			return 0, false // loss of precision without fractions
		}
		amount /= f
	} else {
		amount *= f
	}

	return amount, true
}
