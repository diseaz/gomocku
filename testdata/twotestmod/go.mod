module github.com/diseaz/gomocku/testdata/twotestmod

go 1.23

require github.com/diseaz/gomocku/fortesting/onetestmod v0.0.0-00010101000000-000000000000

replace github.com/diseaz/gomocku/fortesting/onetestmod => ../onetestmod
