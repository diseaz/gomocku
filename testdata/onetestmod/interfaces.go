package onetestmod

type IntToStrConverter interface {
	// ConvertIntToStr converts int value to string value somehow
	ConvertIntToStr(IntValue) (StrValue, error)
}

type StrToIntConverter interface {
	// ConvertStrToInt converts string value to int value somehow
	ConvertStrToInt(in StrValue) (out IntValue, err error)
}

// Converter may convert int to string and vice versa
type Converter interface {
	IntToStrConverter
	StrToIntConverter
}
