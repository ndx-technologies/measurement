package measurement

import (
	"testing"
)

func TestAllUnitUnique(t *testing.T) {
	var unitsMass []string

	for _, q := range UnitMassAll {
		unitsMass = append(unitsMass, q.String())
	}

	var unitsVolume []string
	for _, q := range UnitVolumeAll {
		unitsVolume = append(unitsVolume, q.String())
	}

	all := make(map[string]bool, len(unitsMass)+len(unitsVolume))
	for _, q := range unitsMass {
		all[q] = true
	}
	for _, q := range unitsVolume {
		all[q] = true
	}

	if len(all) != len(unitsMass)+len(unitsVolume) {
		t.Error("duplicates found")
	}
}
