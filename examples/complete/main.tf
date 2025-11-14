terraform {
  required_providers {
    runpod = {
      source = "registry.terraform.io/decentralized-infrastructure/runpod"
    }
  }
}

provider "runpod" {
  api_key = var.runpod_api_key
}

# Create a network volume for persistent storage in a specific region
resource "runpod_network_volume" "regional_volume" {
  name           = "regional-test-volume"
  size           = 10
  data_center_id = "US-CA-2"
}

# Create a GPU pod with the network volume attached
resource "runpod_pod" "example" {
  name              = "example-pod"
  image_name        = "runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel"
  
  # Use valid GPU type names from RunPod API (from error message)
  gpu_type_ids      = [
    # High-end consumer GPUs
    "NVIDIA GeForce RTX 5090",
    "NVIDIA GeForce RTX 5080",
    "NVIDIA GeForce RTX 4090",
    "NVIDIA GeForce RTX 4080 SUPER",
    "NVIDIA GeForce RTX 4080",
    "NVIDIA GeForce RTX 4070 Ti",
    "NVIDIA GeForce RTX 3090 Ti",
    "NVIDIA GeForce RTX 3090",
    "NVIDIA GeForce RTX 3080 Ti",
    "NVIDIA GeForce RTX 3080",
    "NVIDIA GeForce RTX 3070",
    # Professional/Data Center GPUs
    "NVIDIA H200",
    "NVIDIA H100 NVL",
    "NVIDIA H100 PCIe",
    "NVIDIA H100 80GB HBM3",
    "NVIDIA B200",
    "NVIDIA L40S",
    "NVIDIA L40",
    "NVIDIA L4",
    "NVIDIA A100 80GB PCIe",
    "NVIDIA A100-SXM4-80GB",
    "NVIDIA A40",
    "NVIDIA A30",
    # RTX Professional
    "NVIDIA RTX 6000 Ada Generation",
    "NVIDIA RTX 5000 Ada Generation",
    "NVIDIA RTX 4000 Ada Generation",
    "NVIDIA RTX 4000 SFF Ada Generation",
    "NVIDIA RTX A6000",
    "NVIDIA RTX A5000",
    "NVIDIA RTX A4500",
    "NVIDIA RTX A4000",
    "NVIDIA RTX A2000",
    "NVIDIA RTX 2000 Ada Generation",
    # Tesla V100 variants
    "Tesla V100-PCIE-16GB",
    "Tesla V100-SXM2-16GB",
    "Tesla V100-SXM2-32GB",
    "Tesla V100-FHHL-16GB",
    # AMD
    "AMD Instinct MI300X OAM",
  ]
  
  # Restrict to US-CA-2 to match the network volume
  # If capacity is limited, we could expand this list while still preferring US-CA-2
  data_center_ids   = ["US-CA-2", "US-TX-3", "US-IL-1"]
  
  cloud_type        = "COMMUNITY"
  support_public_ip = true
  
  gpu_count  = 1
  volume_in_gb = 20
  container_disk_in_gb = 20
  
  # Try attaching the regional network volume
  network_volume_id = runpod_network_volume.regional_volume.id
  
  # Environment variables
  env = {
    "MY_ENV_VAR" = "example_value"
    "WORKSPACE"  = "/workspace"
  }
  
  # Expose ports
  ports = ["8888/http", "22/tcp"]
}

# Note: Endpoint creation commented out due to template already being bound
# Uncomment and update template_id when you have an available template
# resource "runpod_endpoint" "example" {
#   name        = "example-endpoint"
#   template_id = var.template_id
#   
#   gpu_type_ids = ["NVIDIA GeForce RTX 4090", "NVIDIA A40"]
#   
#   workers_min  = 0
#   workers_max  = 2
#   
#   idle_timeout = 5
#   
#   scaler_type  = "QUEUE_DELAY"
#   scaler_value = 4
# }

# Data sources to list all resources
data "runpod_pods" "all" {}

data "runpod_network_volumes" "all" {}

data "runpod_endpoints" "all" {}
