package measurement

import "strconv"

type DataStorage struct {
	Amount float64         `json:"amount"`
	Unit   UnitDataStorage `json:"unit"`
}

func NewDataStorageFromString(s string) (*DataStorage, error) {
	amount, unit, err := amountUnitFromString(s, MeasureTypeDataStorage, UnitDataStorageAll[:], UnitDataStorageUnknown, UnitDataStorage.String)
	if err != nil {
		return nil, err
	}
	return &DataStorage{Amount: amount, Unit: unit}, nil
}

func (s DataStorage) String() string {
	return strconv.FormatFloat(s.Amount, 'f', -1, 64) + s.Unit.String()
}

func (s *DataStorage) IsZero() bool { return s == nil || s.Amount == 0 }

func (s DataStorage) Convert(unit UnitDataStorage) DataStorage {
	if v, ok := TryConvertExactDataStorage(s.Amount, s.Unit, unit); ok {
		return DataStorage{Amount: v, Unit: unit}
	}
	return DataStorage{Amount: convertDataStorageApproxFromByte(convertDataStorageApproxToByte(s.Amount, s.Unit), unit), Unit: unit}
}

func TryConvertExactDataStorage[T int32 | int64 | float32 | float64](amount T, from, to UnitDataStorage) (v T, ok bool) {
	for _, q := range unitDataStorageLadders {
		if v, ok := convertByLadder(amount, from, to, q); ok {
			return v, true
		}
	}
	return 0, false
}

type UnitDataStorage uint8

//go:generate go-enum-encoding -type=UnitDataStorage -string
const (
	UnitDataStorageUnknown UnitDataStorage = iota // json:""
	UnitBytes                                     // json:"B"
	UnitKiloBytes                                 // json:"kB"
	UnitMegaBytes                                 // json:"MB"
	UnitGigaBytes                                 // json:"GB"
	UnitTeraBytes                                 // json:"TB"
	UnitPetaBytes                                 // json:"PB"
	UnitKibiBytes                                 // json:"KiB"
	UnitMebiBytes                                 // json:"MiB"
	UnitGibiBytes                                 // json:"GiB"
	UnitTebiBytes                                 // json:"TiB"
	UnitPebiBytes                                 // json:"PiB"
)

var UnitDataStorageAll = [...]UnitDataStorage{
	UnitBytes,
	UnitKiloBytes,
	UnitMegaBytes,
	UnitGigaBytes,
	UnitTeraBytes,
	UnitPetaBytes,
	UnitKibiBytes,
	UnitMebiBytes,
	UnitGibiBytes,
	UnitTebiBytes,
	UnitPebiBytes,
}

// decimal prefixes, 1 kB = 1000 B.
var unitDataStorageDecimalLadder = ladder[UnitDataStorage]{
	{UnitBytes, 1},
	{UnitKiloBytes, 1000},
	{UnitMegaBytes, 1000},
	{UnitGigaBytes, 1000},
	{UnitTeraBytes, 1000},
	{UnitPetaBytes, 1000},
}

// binary prefixes, 1 KiB = 1024 B.
var unitDataStorageBinaryLadder = ladder[UnitDataStorage]{
	{UnitBytes, 1},
	{UnitKibiBytes, 1024},
	{UnitMebiBytes, 1024},
	{UnitGibiBytes, 1024},
	{UnitTebiBytes, 1024},
	{UnitPebiBytes, 1024},
}

var unitDataStorageLadders = [...]ladder[UnitDataStorage]{
	unitDataStorageDecimalLadder,
	unitDataStorageBinaryLadder,
}

func convertDataStorageApproxToByte(amount float64, unit UnitDataStorage) float64 {
	switch unit {
	case UnitKiloBytes:
		return amount * 1000
	case UnitMegaBytes:
		return amount * (1000 * 1000)
	case UnitGigaBytes:
		return amount * (1000 * 1000 * 1000)
	case UnitTeraBytes:
		return amount * (1000 * 1000 * 1000 * 1000)
	case UnitPetaBytes:
		return amount * (1000 * 1000 * 1000 * 1000 * 1000)
	case UnitKibiBytes:
		return amount * 1024
	case UnitMebiBytes:
		return amount * (1024 * 1024)
	case UnitGibiBytes:
		return amount * (1024 * 1024 * 1024)
	case UnitTebiBytes:
		return amount * (1024 * 1024 * 1024 * 1024)
	case UnitPebiBytes:
		return amount * (1024 * 1024 * 1024 * 1024 * 1024)
	default:
		return amount
	}
}

func convertDataStorageApproxFromByte(amount float64, unit UnitDataStorage) float64 {
	switch unit {
	case UnitKiloBytes:
		return amount / 1000
	case UnitMegaBytes:
		return amount / (1000 * 1000)
	case UnitGigaBytes:
		return amount / (1000 * 1000 * 1000)
	case UnitTeraBytes:
		return amount / (1000 * 1000 * 1000 * 1000)
	case UnitPetaBytes:
		return amount / (1000 * 1000 * 1000 * 1000 * 1000)
	case UnitKibiBytes:
		return amount / 1024
	case UnitMebiBytes:
		return amount / (1024 * 1024)
	case UnitGibiBytes:
		return amount / (1024 * 1024 * 1024)
	case UnitTebiBytes:
		return amount / (1024 * 1024 * 1024 * 1024)
	case UnitPebiBytes:
		return amount / (1024 * 1024 * 1024 * 1024 * 1024)
	default:
		return amount
	}
}
