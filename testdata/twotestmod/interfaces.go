package twotestmod

import (
	"github.com/diseaz/gomocku/fortesting/onetestmod"
)

type TwoJoiner interface {
	// JoinTwo joins two strings
	JoinTwo(a, b onetestmod.StrValue) onetestmod.StrValue
}

type ManyJoiner interface {
	// JoinMany joins several strings with separator
	JoinMany(sep onetestmod.StrValue, parts ...onetestmod.StrValue) (out onetestmod.StrValue)
}

// Converter may convert int to string and vice versa
type Converter interface {
	onetestmod.IntToStrConverter
	onetestmod.StrToIntConverter

	ConverterTwo()
}
