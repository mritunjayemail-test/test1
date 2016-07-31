---
description: |
    The ProfitBricks builder is able to create images for ProfitBricks cloud.
layout: docs
page_title: ProfitBricks Builder
...

# ProfitBricks Builder

Type: `profitbricks`

The ProfitBricks Builder is able to create virtual machines for [ProfitBricks](https://www.profitbricks.com).

-> **Note:** This builder is not supported by HashiCorp. Please visit [ProfitBricks DevOps Central](https://devops.profitbricks.com/) for support. You may file issues on [GitHub](https://github.com/profitbricks/docker-machine-driver-profitbricks/issues).

## Configuration Reference

There are many configuration options available for the builder. They are
segmented below into two categories: required and optional parameters. Within
each category, the available configuration keys are alphabetized.

In addition to the options listed here, a
[communicator](/docs/templates/communicator.html) can be configured for this
builder.

### Required

-   `username` (string) - ProfitBricks username

-   `password` (string) - ProfitBrucks password

-   `image` (string) - ????

### Optional

-   `cores` (int) - Amount of CPU cores to use for this build. Defaults to `4`.

-   `disk_size` (string) - Amount of disk space for this image. Defaults to `50gb`

-   `disktype` (string) - Type of disk to use for this image. Defaults to `HDD`.

-   `ram` (int) - Amount of RAM to use for this image. Defalts to `2048`.

-   `region` (string) - Defaults to `us/las`.

-   `snapshot_name` (string) - If snapshot name is not provided Packer will generate it

-   `snapshot_password` (string) - ????

-   `ssh_key_path` (string) - Path to private SSHkey. If no path to the key is provided Packer will create one under the name [snapshot_name]

-   `url` (string) - Endpoint for the ProfitBricks REST API

## Example

Here is a basic example:

```json
{
  "builders": [
    {
      "image": "Ubuntu-16.04",
      "type": "profitbricks",
      "disk_size": "5",
      "snapshot_name": "double",
      "ssh_key_path": "/path/to/private/key",
      "snapshot_password": "test1234",
      "timeout": 100
    }
  ],
  "provisioners": [
    {
      "inline": [
        "echo foo"
      ],
      "type": "shell"
    }
  ]
}
```
