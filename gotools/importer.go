package gotools

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
)

type srcImporter struct {
	fileSet *token.FileSet
}

func NewSrcImporter(fileSet *token.FileSet) *srcImporter {
	if fileSet == nil {
		fileSet = token.NewFileSet()
	}
	return &srcImporter{
		fileSet: fileSet,
	}
}

func (imp *srcImporter) Import(path string) (*types.Package, error) {
	pkgPath, fileList, err := ListDir(path)
	if err != nil {
		return nil, err
	}

	return imp.importPackage(pkgPath, fileList, nil)
}

func (imp *srcImporter) ImportFrom(path, dir string, _ types.ImportMode) (*types.Package, error) {
	pkgPath, fileList, err := ListPackage(path, dir)
	if err != nil {
		return nil, err
	}

	return imp.importPackage(pkgPath, fileList, nil)
}

func (imp *srcImporter) ImportWithInfo(path string, info *types.Info) (*types.Package, error) {
	pkgPath, fileList, err := ListDir(path)
	if err != nil {
		return nil, err
	}

	return imp.importPackage(pkgPath, fileList, info)
}

func (imp *srcImporter) importPackage(pkgPath string, fileList []string, info *types.Info) (*types.Package, error) {
	var astFiles []*ast.File
	for _, filePath := range fileList {
		// log.Printf("Parsing %q", filePath)
		astFile, err := parser.ParseFile(imp.fileSet, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		astFiles = append(astFiles, astFile)
	}

	conf := types.Config{Importer: imp}
	return conf.Check(pkgPath, imp.fileSet, astFiles, info)
}
