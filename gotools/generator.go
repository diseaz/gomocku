package gotools

import (
	"bytes"
	"fmt"
	"go/types"
	"io"

	"github.com/huandu/xstrings"
)

type MockGenerator struct {
	qualifier *Qualifier
	pkg       *types.Package
	name      string
	baseName  string
	intf      *types.Interface
}

func NewMockGenerator(
	qualifier *Qualifier,
	namedIntf *types.Named,
) *MockGenerator {
	intf := namedIntf.Underlying().(*types.Interface)
	name := namedIntf.Obj().Name()
	return &MockGenerator{
		qualifier: qualifier,
		name:      name,
		baseName:  xstrings.ToCamelCase(name),
		pkg:       namedIntf.Obj().Pkg(),
		intf:      intf,
	}
}

func (g *MockGenerator) FileName() string {
	return xstrings.ToKebabCase(g.name) + "_mock_gen.go"
}

func (g *MockGenerator) GenerateFile(w io.Writer) {
	fmt.Fprintf(w, "package %s\n", g.pkg.Name())

	var mockBuf bytes.Buffer
	g.GeneratePrologue(&mockBuf)
	g.GenerateM(&mockBuf)
	g.GenerateCtl(&mockBuf)
	g.GenerateMock(&mockBuf)

	imps := g.qualifier.Imports()

	if len(imps) > 0 {
		fmt.Fprintf(w, "\nimport (\n")
		for _, imp := range imps {
			fmt.Fprintln(w, imp)
		}
		fmt.Fprintln(w, ")")
	}
	mockBuf.WriteTo(w)
}

func (g *MockGenerator) GeneratePrologue(w io.Writer) {
	var prefix string

	if pkgName := g.qualifier.Qualify(g.pkg); pkgName != "" {
		prefix = pkgName + "."
	}
	intfName := prefix + g.name

	fmt.Fprintf(w, `
var _ %[2]s = (*%[1]sMock)(nil)

// See docs at https://www.notion.so/joomteam/Mocks-14aa9966b0b280cbaa64ed9dcd28f6f6
func New%[3]sMock() (*%[1]sMock, *%[1]sCtl) {
	ctl := &%[1]sCtl{}
	return ctl.Mock(), ctl
}
`, g.baseName, intfName, xstrings.ToPascalCase(g.name))

}

func (g *MockGenerator) GenerateM(w io.Writer) {
	fmt.Fprintf(w, "\ntype %[1]sM struct {\n", g.baseName)

	var methods bytes.Buffer
	numMethods := g.intf.NumMethods()
	for mIdx := 0; mIdx < numMethods; mIdx++ {
		method := g.intf.Method(mIdx)
		sig := method.Signature()

		fmt.Fprintf(&methods, "\t%s func", method.Name())
		types.WriteSignature(&methods, AnonymSignature(sig), g.qualifier.Qualify)
		fmt.Fprintln(&methods)
	}
	methods.WriteTo(w)

	fmt.Fprintln(w, "}")
}

func (g *MockGenerator) GenerateCtl(w io.Writer) {
	fmt.Fprintf(w, `
type %[1]sCtl struct {
	M %[1]sM
	stack []%[1]sM
}

func (ctl *%[1]sCtl) Mock() *%[1]sMock {
	return &%[1]sMock{
		ctl: ctl,
	}
}

func (ctl *%[1]sCtl) Push() *%[1]sCtl {
	ctl.stack = append(ctl.stack, ctl.M)
	return ctl
}

func (ctl *%[1]sCtl) Pop() {
	i := len(ctl.stack) - 1
	ctl.M = ctl.stack[i]
	ctl.stack = ctl.stack[:i]
}

func (ctl *%[1]sCtl) PopAll() {
	if len(ctl.stack) == 0 {
		return
	}
	ctl.M = ctl.stack[0]
	ctl.stack = ctl.stack[:0]
}
`, g.baseName)
}

func (g *MockGenerator) GenerateMock(w io.Writer) {
	// Mock type
	fmt.Fprintf(w, `
type %[1]sMock struct {
	ctl *%[1]sCtl
}
`, g.baseName)

	// Mock methods implementation
	numMethods := g.intf.NumMethods()
	for mIdx := 0; mIdx < numMethods; mIdx++ {
		method := g.intf.Method(mIdx)
		sig := SafeSignature(method.Signature())

		// Method signature
		fmt.Fprintf(w, "\nfunc (m *%[1]sMock) %[2]s", g.baseName, method.Name())
		var sigBuf bytes.Buffer
		types.WriteSignature(&sigBuf, sig, g.qualifier.Qualify)
		fmt.Fprintf(w, "%[1]s {\n", sigBuf.String())

		// Method body
		if sig.Results().Len() == 0 {
			fmt.Fprintf(w, "\tm.ctl.M.%[1]s(", method.Name())
		} else {
			fmt.Fprintf(w, "\treturn m.ctl.M.%[1]s(", method.Name())
		}

		params := sig.Params()
		for i := 0; i < params.Len(); i++ {
			if i != 0 {
				fmt.Fprint(w, ", ")
			}
			param := params.At(i)
			fmt.Fprint(w, param.Name())
			if sig.Variadic() && i == params.Len()-1 {
				fmt.Fprint(w, "...")
			}
		}

		fmt.Fprintln(w, ")")

		// Method tail
		fmt.Fprintln(w, "}")
	}
}
