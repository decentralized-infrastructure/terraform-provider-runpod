# RunPod Terraform Provider

This repository contains the Community [RunPod](https://www.runpod.io) Terraform Provider. 

**Note that this provider is NOT built by RunPod - it is a community effort.** Please do not ask RunPod for support, as they will not be able to provide it. Please report issues in the GitHub repository.

## Features

- **Pod Management**: Create, read, update, and delete GPU pods
- **Serverless Endpoints**: Manage serverless endpoints for inference
- **Template Management**: Create and manage pod templates
- **Network Volumes**: Manage persistent network storage
- **Container Registry**: Configure container registry authentication

## Supported Resources

### Resources
- `runpod_pod` - Manage GPU pods
- `runpod_endpoint` - Manage serverless endpoints
- `runpod_template` - Manage pod templates
- `runpod_network_volume` - Manage network volumes
- `runpod_container_registry_auth` - Manage container registry authentication

### Data Sources
- `runpod_pods` - List all pods in your account
- `runpod_endpoints` - List all serverless endpoints
- `runpod_templates` - List available templates
- `runpod_network_volumes` - List network volumes

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23
- RunPod API key (get one at [runpod.io](https://www.runpod.io/console/user/settings))

## Quick Start

1. Set your API key:
   ```bash
   export RUNPOD_API_KEY="your-api-key-here"
   ```

2. Create a basic configuration:
   ```hcl
   terraform {
     required_providers {
       runpod = {
         source = "registry.terraform.io/decentralized-infrastructure/runpod"
       }
     }
   }

   provider "runpod" {
     # API key from environment variable RUNPOD_API_KEY
   }

   # Create a GPU pod
   resource "runpod_pod" "example" {
     name          = "my-gpu-pod"
     gpu_type      = "NVIDIA RTX A4000"
     gpu_count     = 1
     container_image = "runpod/pytorch:latest"
   }
   ```

3. Deploy:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Examples

See the [examples/](./examples/) directory for comprehensive usage examples.

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider:

```shell
go build -o terraform-provider-runpod
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).

To add a new dependency:

```shell
go get github.com/author/dependency
go mod tidy
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine.

To compile the provider, run `go build`. This will build the provider and put the provider binary in the current directory.

For development testing, you can use a `.terraformrc` file to override the provider location:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/decentralized-infrastructure/runpod" = "/path/to/your/go/bin"
  }
  direct {}
}
```

## Documenting the Provider

In order to generate documentation for the provider, the following command can be run:
```
go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
tfplugindocs generate
```

## License

This project is licensed under the Mozilla Public License 2.0 - see the LICENSE file for details.
