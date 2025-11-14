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

# Example of creating a basic GPU pod
resource "runpod_pod" "example" {
  name              = "basic-example-pod"
  image_name        = "runpod/pytorch:2.1.0-py3.10-cuda11.8.0-devel"
  
  # Specify GPU types (will use first available)
  gpu_type_ids      = [
    "NVIDIA GeForce RTX 4090",
    "NVIDIA GeForce RTX 3090",
    "NVIDIA A40",
    "NVIDIA RTX A6000",
  ]
  
  # Specify data centers
  data_center_ids   = [
    "US-CA-2",
    "US-TX-3",
    "EU-RO-1",
  ]
  
  gpu_count            = 1
  cloud_type           = "COMMUNITY"
  support_public_ip    = true
  
  # Storage
  volume_in_gb         = 20
  container_disk_in_gb = 20
  
  # Environment variables
  env = {
    MY_VAR = "example_value"
  }
  
  # Expose ports
  ports = ["8888/http", "22/tcp"]
}

# Output the pod details
output "pod_id" {
  value = runpod_pod.example.id
}

output "pod_status" {
  value = runpod_pod.example.desired_status
}

output "pod_cost" {
  value = runpod_pod.example.cost_per_hr
}

output "public_ip" {
  value = runpod_pod.example.public_ip
}
