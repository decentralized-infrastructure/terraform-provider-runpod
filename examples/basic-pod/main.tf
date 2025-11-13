terraform {
  required_providers {
    runpod = {
      source = "registry.terraform.io/decentralized-infrastructure/runpod"
      version = "~> 0.1"
    }
  }
}

provider "runpod" {
  # API key from environment variable RUNPOD_API_KEY
}

# Example of creating a basic GPU pod
# resource "runpod_pod" "example" {
#   name            = "example-pod"
#   gpu_type        = "NVIDIA RTX A4000"
#   gpu_count       = 1
#   container_image = "runpod/pytorch:latest"
#   
#   # Optional configuration
#   env_vars = {
#     CUSTOM_VAR = "value"
#   }
# }
