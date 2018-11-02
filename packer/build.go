package packer

const (
	// This is the key in configurations that is set to the name of the
	// build.
	BuildNameConfigKey = "packer_build_name"

	// This is the key in the configuration that is set to the type
	// of the builder that is run. This is useful for provisioners and
	// such who want to make use of this.
	BuilderTypeConfigKey = "packer_builder_type"

	// This is the key in configurations that is set to "true" when Packer
	// debugging is enabled.
	DebugConfigKey = "packer_debug"

	// This is the key in configurations that is set to "true" when Packer
	// force build is enabled.
	ForceConfigKey = "packer_force"

	// This key determines what to do when a normal multistep step fails
	// - "cleanup" - run cleanup steps
	// - "abort" - exit without cleanup
	// - "ask" - ask the user
	OnErrorConfigKey = "packer_on_error"

	// TemplatePathKey is the path to the template that configured this build
	TemplatePathKey = "packer_template_path"

	// This key contains a map[string]string of the user variables for
	// template processing.
	UserVariablesConfigKey = "packer_user_variables"
)
