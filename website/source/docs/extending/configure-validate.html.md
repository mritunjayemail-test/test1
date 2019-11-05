### The "Configure" Method

`Configure` is Called for each plugin prior to any runs with the configuration
that was given in the template. Configuration is passed as an `interface{}`
type, that can be multiple types.

`Configure` is responsible for translating this configuration into an internal
structure, validating it, and returning any errors. Hopefully Packer defines
some helpers to help you get this sorted out depending on the underlying type
received.

 * a `[]map[string]interface{}` will be passed when using a single build file.
 * a `"github.com/zclconf/go-cty/cty".Value` when using HCL2 by building
   passing an HCL file, a folder or a `pkr.json` file.


For decoding a `map[string]interface` into a meaningful structure, the
[mapstructure](https://github.com/mitchellh/mapstructure) library is
recommended. Mapstructure will take an `interface{}` and decode it into an
arbitrarily complex struct. If there are any errors, it generates very human
friendly errors that can be returned directly from the prepare method.

While it is not actively enforced, **no side effects** should occur from
running the `Configure` method. Specifically, don't create files, don't launch
virtual machines, etc. `Configure`'s purpose is solely to configure the
builder.

The `Configure` method is called very early in the build process so that errors
may be displayed to the user before anything actually happens. `Configure` may
return basic validation errors during parse time.

### The "Validate" Method

`Validate` is Called after `Configure` it allows a plugin to tell wether
something is setup incorrectly. For example error if field `A` OR field `B`
MUST be set.
