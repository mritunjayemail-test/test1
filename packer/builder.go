package packer

import "context"

// Packer Builders are the components of Packer responsible for creating a
// machine, bringing it to a point where it can be provisioned, and then
// turning that provisioned machine into some sort of machine image.
//
// Implementers of Builder are responsible for actually building images on some
// platform given some configuration.
//
type Builder interface {
	Plugin
	// Build is where the actual build should take place. It takes a Build and
	// a Ui.
	//
	// Build is where all the interesting stuff happens. Build is executed,
	// often in parallel for multiple builders, to actually build the machine,
	// provision it, and create the resulting machine image, which is returned
	// as an implementation of the `packer.Artifact` interface.
	//
	// Build takes three parameters. The `packer.Ui` object is used to send
	// output to the console. `packer.Hook` is used to execute hooks, which are
	// covered in more detail in the hook section below.
	//
	// Because builder runs are typically a complex set of many steps, the
	// [multistep](https://github.com/hashicorp/packer/blob/master/helper/multistep)
	// helper is recommended to bring order to the complexity. Multistep is a
	// library which allows you to separate your logic into multiple distinct
	// "steps" and string them together. It fully supports cancellation
	// mid-step and so on. Please check it out, it is how the built-in builders
	// are all implemented.
	//
	// Build returns an implementation of `packer.Artifact`. If something goes
	// wrong during the build, an error can be returned, as well. Note that it
	// is perfectly fine to produce no artifact and no error, although this is
	// rare.
	//
	Build(context.Context, Ui, Hook) (Artifact, error)
}
