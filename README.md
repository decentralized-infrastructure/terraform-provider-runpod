# RunPod Terraform Provider

Official Terraform provider for [RunPod](https://www.runpod.io) - manage GPU cloud infrastructure as code.

## Features

- **Pod Management**: Create, read, and delete GPU/CPU pods with flexible configuration
- **Network Volumes**: Manage persistent storage volumes across pods
- **Serverless Endpoints**: Deploy and manage serverless GPU endpoints
- **Data Sources**: Query existing resources including pods, volumes, endpoints, and templates

## Supported Resources

### Resources
- `runpod_pod` - Manage GPU/CPU pods
- `runpod_network_volume` - Manage persistent network storage
- `runpod_endpoint` - Manage serverless endpoints

### Data Sources
- `runpod_pods` - List all pods
- `runpod_network_volumes` - List all network volumes
- `runpod_endpoints` - List all serverless endpoints
- `runpod_templates` - List available templates

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
         source = "decentralized-infrastructure/runpod"
       }
     }
   }

   provider "runpod" {
     # API key from environment variable RUNPOD_API_KEY
   }

   # Create a GPU pod
   resource "runpod_pod" "example" {
     name              = "my-gpu-pod"
     image_name        = "runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel"
     gpu_type_ids      = ["NVIDIA GeForce RTX 4090", "NVIDIA A40"]
     data_center_ids   = ["US-CA-2", "US-TX-3"]
     
     gpu_count            = 1
     cloud_type           = "COMMUNITY"
     support_public_ip    = true
     volume_in_gb         = 20
     container_disk_in_gb = 20
     
     ports = ["8888/http", "22/tcp"]
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
