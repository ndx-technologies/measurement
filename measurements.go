package measurement

// Measurements is a container describing single instance in measurement dimensions.
// Not all dimensions may be present.
type Measurements struct {
	Quantity float32   `json:"quantity,omitzero"`
	Mass     *Mass     `json:"mass,omitzero"`
	Volume   *Volume   `json:"volume,omitzero"`
	Length   *Length   `json:"length,omitzero"`
	Area     *Area     `json:"area,omitzero"`
	Energy   *Energy   `json:"energy,omitzero"`
	Duration *Duration `json:"duration,omitzero"`
}

func (s Measurements) IsZero() bool {
	return s.Quantity == 0 &&
		s.Mass.IsZero() &&
		s.Volume.IsZero() &&
		s.Length.IsZero() &&
		s.Area.IsZero() &&
		s.Energy.IsZero() &&
		s.Duration.IsZero()
}

// IsEqual without units conversion
func (s Measurements) IsEqual(v Measurements) bool {
	return s.Quantity == v.Quantity &&
		equalPtr(s.Mass, v.Mass) &&
		equalPtr(s.Volume, v.Volume) &&
		equalPtr(s.Length, v.Length) &&
		equalPtr(s.Area, v.Area) &&
		equalPtr(s.Energy, v.Energy) &&
		equalPtr(s.Duration, v.Duration)
}

// equalPtr treats two nils as equal, nil and non-nil as different, and non-nil by value.
func equalPtr[T comparable](a, b *T) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	return a == nil || *a == *b
}
