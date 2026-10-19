package measurement

import (
	"testing"
)

func TestAllUnitUnique(t *testing.T) {
	var units []string

	for _, q := range UnitMassAll {
		units = append(units, q.String())
	}
	for _, q := range UnitVolumeAll {
		units = append(units, q.String())
	}
	for _, q := range UnitLengthAll {
		units = append(units, q.String())
	}
	for _, q := range UnitAreaAll {
		units = append(units, q.String())
	}
	for _, q := range UnitEnergyAll {
		units = append(units, q.String())
	}
	for _, q := range UnitDurationAll {
		units = append(units, q.String())
	}

	all := make(map[string]bool, len(units))
	for _, q := range units {
		all[q] = true
	}

	if len(all) != len(units) {
		t.Error("duplicates found")
	}
}
