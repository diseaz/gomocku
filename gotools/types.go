package gotools

import (
	"go/types"
	"iter"
	"log"
)

func IterateNamed(typ types.Type) iter.Seq[*types.Named] {
	return func(yield func(*types.Named) bool) {
		var process func(types.Type) bool
		process = func(typ types.Type) bool {
			switch t := typ.(type) {
			case *types.Named:
				return yield(t)
			case *types.Basic:
				return true
			case *types.Array:
				return process(t.Elem())
			case *types.Slice:
				return process(t.Elem())
			case *types.Map:
				return process(t.Elem()) && process(t.Key())
			case *types.Pointer:
				return process(t.Elem())
			case *types.Signature:
				// TODO: implement
			case *types.Struct:
				// TODO: implement
			default:
				log.Printf("Unhandled type: %T", t)
			}
			return true
		}
		process(typ)
	}
}

func AsNamedInterface(typ types.Type) (*types.Interface, *types.Named) {
	named, isNamed := typ.(*types.Named)
	if !isNamed {
		return nil, nil
	}
	intf, _ := typ.Underlying().(*types.Interface)
	return intf, named
}

type PkgImport struct {
	*types.PkgName
	implicit bool
}

func NewPkgImport(pkgName *types.PkgName, implicit bool) *PkgImport {
	return &PkgImport{
		PkgName:  pkgName,
		implicit: implicit,
	}
}

func (pi *PkgImport) Implicit() bool {
	return pi.implicit
}
