package measurement

type MeasureType uint8

//go:generate go-enum-encoding -type=MeasureType -string
const (
	MeasureTypeUndefined MeasureType = iota // json:""
	MeasureTypeMass                         // json:"mass"
	MeasureTypeVolume                       // json:"volume"
)
