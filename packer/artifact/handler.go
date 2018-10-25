package artifact

import (
	"context"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/packer/packer/config"
)

// A Handler prepares or builds artifacts
// given a packer artifact configuration and
// and an artifact Index.
//
// Handle could build or validate said artifact,
// depending on artifactConfig.Operation
//
// Handle can only be called after a parent (or source), if any.
//
// For a preparation, Handle should add an expected
// built artifact to Index, so other artifacts can validate
// and/or prepare based on output artifact(s). The index
// will be different from the build step.
//
// In case of a build, Handle should add the built artifact
// to Index as a following build may use output.
//
// Context will serve a as a bag for request scoped objects
// such as Ui and logging facilities.
// cancelling ctx will
type Handler interface {
	Handle(ctx context.Context, artifactConfig config.Artifact, idx Index) hcl.Diagnostics
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as artifact handlers.
// If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.
type HandlerFunc func(ctx context.Context, artifactConfig config.Artifact, idx Index) hcl.Diagnostics

// Handle calls f(ctx, artifactConfig, idx).
func (f HandlerFunc) Handle(ctx context.Context, artifactConfig config.Artifact, idx Index) hcl.Diagnostics {
	return f(ctx, artifactConfig, idx)
}
