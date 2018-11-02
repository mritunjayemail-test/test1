package packer

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/packer/packer/artifact"

	"github.com/hashicorp/packer/packer/config"
)

// Core is the main executor of Packer. If Packer is being used as a
// library, this is the struct you'll want to instantiate to get anything done.
type Core struct {
	config *config.Root

	getComponent ComponentFinder
	secrets      []string
	operations   []config.Operation

	debug bool
}

// ComponentFinder is a function
// pointer necessary to look up components of Packer such as builders,
// commands, etc.
type ComponentFinder func(name string) (artifact.Handler, error)

// An Option allows to configure a Core
// packer configuration.
type Option func(*Core) error

// AllOperations is the list of operation called
// on each packer handler referenced by config.
var AllOperations = []config.Operation{
	config.Validation,
	config.Build,
}

// NewCore creates and configures a new Core.
//
// NewCore will flatten cfg.Artifacts into a
// 1 dimensional array for traversal simplicity
//
func NewCore(cfg *config.Root, cpts ComponentFinder, options ...Option) (*Core, error) {
	result := &Core{
		config:       cfg,
		getComponent: cpts,
		operations:   AllOperations,
	}

	for _, opt := range options {
		if err := opt(result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Run executes set of instructions given configurations
func (c *Core) Run(ctx context.Context) {

	type ContextCancel struct {
		Context context.Context
		Cancel  func()
	}
	contexts := map[string]ContextCancel{}
	diags := diagnosticReceiver{}

	for _, operation := range c.operations {
		ctx := context.WithValue(ctx, "operation", operation)

		idx := artifact.Index{}

		// initialize each artifact's context
		for _, artifact := range c.config.Artifacts {
			childCtx, cancel := context.WithCancel(ctx)
			contexts[artifact.FullName()] = ContextCancel{childCtx, cancel}
		}

		// start artifacts
		wg := errgroup.Group{}
		for _, artifact := range c.config.Artifacts {
			artifact := artifact // avoid race contitions
			childCtx := contexts[artifact.FullName()].Context
			cancel := contexts[artifact.FullName()].Cancel

			// run artifact command
			wg.Go(func() error {
				defer cancel()

				if artifact.Source != nil {
					sourceCtx := contexts[*artifact.Source].Context

					<-sourceCtx.Done() // wait for source to be done

					if diags.Diagnostics().HasErrors() {
						// there is an error somewhere
						// let's leave quietly
						return nil
					}
				}

				handler, err := c.getComponent(artifact.Type)
				if err != nil {
					diags.Append(&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  fmt.Sprintf("Error getting component %s", artifact.Type),
						Detail:   fmt.Sprintf("%v", err),
						// TODO(azr): find way to add detail on where
						// type is defined in file for a more precise output
					})
				}

				newDiags := handler.Handle(childCtx, artifact, idx)
				diags.Extend(newDiags)
				return nil
			})

			if c.debug {
				//TODO(azr): ui say why we wait
				wg.Wait()
			}
		}

		wg.Wait()
		if diags.Diagnostics().HasErrors() {
			break
		}
	}

	wr := hcl.NewDiagnosticTextWriter(
		os.Stdout,      // writer to send messages to
		c.config.Files, // the parser's file cache, for source snippets
		78,             // wrapping width
		true,           // generate colored/highlighted output
	)
	err := wr.WriteDiagnostics(diags.Diagnostics())
	if err != nil {
		panic("error: " + err.Error()) // TODO(azr): remove me
	}
}
