package packer

import (
	"sync"

	"github.com/hashicorp/hcl2/hcl"
)

type diagnosticReceiver struct {
	sync.Locker
	diagnostics hcl.Diagnostics
}

func (dr *diagnosticReceiver) Append(diag *hcl.Diagnostic) {
	dr.Lock()
	defer dr.Unlock()
	dr.diagnostics.Append(diag)
}

func (dr *diagnosticReceiver) Extend(diags hcl.Diagnostics) {
	dr.Lock()
	defer dr.Unlock()
	dr.diagnostics.Extend(diags)
}

func (dr *diagnosticReceiver) Diagnostics() hcl.Diagnostics {
	dr.Lock()
	defer dr.Unlock()
	return dr.diagnostics
}
