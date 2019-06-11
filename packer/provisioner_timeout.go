package packer

import (
	"context"
	"fmt"
	"time"
)

// TimeoutProvisioner is a Provisioner implementation that can timeout after a
// duration
type TimeoutProvisioner struct {
	Provisioner
	Timeout time.Duration
}

func (p *TimeoutProvisioner) Provision(ctx context.Context, ui Ui, comm Communicator) error {
	errC := make(chan error)
	go func() {
		// init the timeout from the same goroutine
		ctx, cancel := context.WithTimeout(ctx, p.Timeout)
		defer cancel()

		ui.Say(fmt.Sprintf("Setting a %s timeout for the next provisioner...", p.Timeout))
		errC <- p.Provisioner.Provision(ctx, ui, comm)
	}()

	select {
	case err := <-errC:
		// all good
		return err
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded { 
			ui.Error("Cancelling provisioner after a timeout...")
		}
		select {
		case err := <-errC:
			// provisioner probably returned after a cleanup
			return err
		//  give the provisioner 1 minute to cleanup/return
		case <-time.After(1 * time.Minute):
			ui.Error("Provisioner did not return/cleanup quick enough and is probably locked, aborting.")
			return ctx.Err()
		}
	}
}
