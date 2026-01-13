package measurement

// Measurements is a container describing single instance in measurement dimensions.
// Not all dimensions may be present.
type Measurements struct {
	Quantity float32 `json:"quantity,omitzero"`
	Mass     *Mass   `json:"mass,omitzero"`
	Volume   *Volume `json:"volume,omitzero"`
}

func (s *Measurements) IsZero() bool {
	if s == nil {
		return true
	}
	return s.Quantity == 0 && s.Mass.IsZero() && s.Volume.IsZero()
}
