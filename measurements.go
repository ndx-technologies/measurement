package measurement

// Measurements is a container describing single instance in measurement dimensions.
// Not all dimensions may be present.
type Measurements struct {
	Quantity float32 `json:"quantity,omitzero"`
	Mass     *Mass   `json:"mass,omitzero"`
	Volume   *Volume `json:"volume,omitzero"`
}

func (s Measurements) IsZero() bool { return s.Quantity == 0 && s.Mass.IsZero() && s.Volume.IsZero() }

// IsEqual without units conversion
func (s Measurements) IsEqual(v Measurements) bool {
	if s.Quantity != v.Quantity {
		return false
	}
	if (s.Mass == nil) != (v.Mass == nil) {
		return false
	}
	if s.Mass != nil && *s.Mass != *v.Mass {
		return false
	}
	if (s.Volume == nil) != (v.Volume == nil) {
		return false
	}
	return s.Volume == nil || *s.Volume == *v.Volume
}
