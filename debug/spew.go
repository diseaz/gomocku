package debug

import (
	"github.com/diseaz/go-spew/spew"
)

var CompactDump = func() (r spew.ConfigState) {
	r = spew.Compact
	r.DisableMethods = true
	return r
}()
