package measurement

type MeasureType uint8

//go:generate go-enum-encoding -type=MeasureType -string
const (
	MeasureTypeUndefined   MeasureType = iota // json:""
	MeasureTypeMass                           // json:"mass"
	MeasureTypeVolume                         // json:"volume"
	MeasureTypeLength                         // json:"length"
	MeasureTypeArea                           // json:"area"
	MeasureTypeEnergy                         // json:"energy"
	MeasureTypeDuration                       // json:"duration"
	MeasureTypePower                          // json:"power"
	MeasureTypeDataStorage                    // json:"data_storage"
	MeasureTypeFrequency                      // json:"frequency"
	MeasureTypeCharge                         // json:"charge"
	MeasureTypeVoltage                        // json:"voltage"
)
