package gotools

import (
	"encoding/hex"
	"fmt"
	"go/types"
	"hash/fnv"
	"log"
)

type PackagesWithAliases map[*types.Package]*PackageAliases

func NewPackagesWithAliases() PackagesWithAliases {
	return make(PackagesWithAliases)
}

func (pwa PackagesWithAliases) AddExplicit(pkg *types.Package, name string) {
	if name == pkg.Name() {
		pwa.getAliases(pkg).AddImplicit()
	} else {
		pwa.getAliases(pkg).AddExplicit(name)
	}
}

func (pwa PackagesWithAliases) AddImplicit(pkg *types.Package) {
	pwa.getAliases(pkg).AddImplicit()
}

func (pwa PackagesWithAliases) getAliases(pkg *types.Package) *PackageAliases {
	pa, hasAliases := pwa[pkg]
	if !hasAliases {
		pa = NewPackageAliases()
		pwa[pkg] = pa
	}
	return pa
}

type PackageAliases struct {
	Explicit      map[string]int // Count of each alias
	ImplicitCount int
}

func NewPackageAliases() *PackageAliases {
	return &PackageAliases{}
}

func (pa *PackageAliases) AddExplicit(name string) {
	if pa.Explicit == nil {
		pa.Explicit = make(map[string]int)
	}
	pa.Explicit[name]++
}

func (pa *PackageAliases) AddImplicit() {
	pa.ImplicitCount++
}

type Qualifier struct {
	src      map[*types.Package]struct{}
	name2pkg map[string]*types.Package
	pkg2name map[*types.Package]string
}

func NewQualifier(srcPkg *types.Package) *Qualifier {
	return &Qualifier{
		src: map[*types.Package]struct{}{
			srcPkg: struct{}{},
		},
		name2pkg: make(map[string]*types.Package),
		pkg2name: make(map[*types.Package]string),
	}
}

func (q *Qualifier) Qualify(pkg *types.Package) string {
	if _, isSrc := q.src[pkg]; isSrc {
		return ""
	}

	if name, hasName := q.pkg2name[pkg]; hasName {
		return name
	}

	name := pkg.Name()
	if _, conflict := q.name2pkg[name]; !conflict {
		q.name2pkg[name] = pkg
		q.pkg2name[pkg] = name
		return name
	}

	hasher := fnv.New64a()
	name = fmt.Sprintf("pkg_%s", hex.EncodeToString(hasher.Sum([]byte(pkg.Path()))))
	if _, conflict := q.name2pkg[name]; !conflict {
		q.name2pkg[name] = pkg
		q.pkg2name[pkg] = name
		return name
	}

	log.Fatal(fmt.Errorf("can't generate unique alias for package %q", pkg.Path()))
	return pkg.Path()
}

func (q *Qualifier) Imports() []string {
	if len(q.pkg2name) == 0 {
		return nil
	}
	result := make([]string, 0, len(q.pkg2name))
	for pkg, name := range q.pkg2name {
		var imp string
		if name == pkg.Name() {
			imp = fmt.Sprintf("\t%q", pkg.Path())
		} else {
			imp = fmt.Sprintf("\t%s %q", name, pkg.Path())
		}
		result = append(result, imp)
	}
	return result
}
