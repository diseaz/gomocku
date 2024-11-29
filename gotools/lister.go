package gotools

import (
	"fmt"
	"log"

	"github.com/diseaz/gomocku/debug"
	"golang.org/x/tools/go/packages"
)

func ListDir(srcPath string) (string, []string, error) {
	cfg := &packages.Config{
		Dir:  srcPath,
		Mode: packages.LoadImports | packages.NeedModule,
	}

	pkgs, err := packages.Load(cfg)
	if err != nil {
		return "", nil, err
	}

	if len(pkgs) > 1 {
		log.Printf("ListDir(%q) ->\n%s", srcPath, debug.CompactDump.Sdump(pkgs))
		return "", nil, fmt.Errorf("too many packages in result: %d", len(pkgs))
	} else if len(pkgs) == 0 {
		return "", nil, fmt.Errorf("no packages in result")
	}

	pkg := pkgs[0]
	// log.Printf("ListDir(%q) ->\n%s", srcPath, debug.CompactDump.Sdump(pkg))
	return pkg.PkgPath, pkg.CompiledGoFiles, nil
}

func ListPackage(pkgPath string, srcPath string) (string, []string, error) {
	cfg := &packages.Config{
		Dir:  srcPath,
		Mode: packages.LoadImports | packages.NeedModule,
	}

	pkgs, err := packages.Load(cfg, "pattern="+pkgPath)
	if err != nil {
		return "", nil, err
	}

	if len(pkgs) > 1 {
		log.Printf("ListPackage(%q, %q) ->\n%s", pkgPath, srcPath, debug.CompactDump.Sdump(pkgs))
		return "", nil, fmt.Errorf("too many packages in result: %d", len(pkgs))
	} else if len(pkgs) == 0 {
		return "", nil, fmt.Errorf("no packages in result")
	}

	pkg := pkgs[0]
	// log.Printf("ListPackage(%q, %q) ->\n%s", pkgPath, srcPath, debug.CompactDump.Sdump(pkg))
	return pkg.PkgPath, pkg.CompiledGoFiles, nil
}
