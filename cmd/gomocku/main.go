package main

import (
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/diseaz/gomocku/gotools"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal(fmt.Errorf("path to a source must be specified"))
	}
	srcPath := flag.Arg(0)

	fileSet := token.NewFileSet()
	imp := gotools.NewSrcImporter(fileSet)

	selfPkg, err := imp.ImportWithInfo(srcPath, nil)
	if err != nil {
		log.Fatal(err)
	}

	scope := selfPkg.Scope()
	qualifier := gotools.NewQualifier(selfPkg)

	for _, n := range scope.Names() {
		obj := scope.Lookup(n)

		typInterface, typNamed := gotools.AsNamedInterface(obj.Type())
		if typInterface == nil {
			continue
		}

		generator := gotools.NewMockGenerator(qualifier, typNamed)
		mockFilePath := filepath.Join(srcPath, generator.FileName())
		log.Printf("Generating %q", mockFilePath)
		outFile, err := os.Create(mockFilePath)
		if err != nil {
			log.Fatal(err)
		}
		func() {
			defer outFile.Close()
			generator.GenerateFile(outFile)
		}()
	}
}
