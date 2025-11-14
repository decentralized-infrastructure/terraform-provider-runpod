output "network_volume_id" {
  description = "ID of the created network volume"
  value       = runpod_network_volume.regional_volume.id
}

output "pod_id" {
  description = "ID of the created pod"
  value       = runpod_pod.example.id
}

output "pod_status" {
  description = "Status of the created pod"
  value       = runpod_pod.example.desired_status
}

output "pod_cost_per_hour" {
  description = "Cost per hour for the pod"
  value       = runpod_pod.example.cost_per_hr
}

output "endpoint_id" {
  description = "ID of the created endpoint"
  value       = null # runpod_endpoint.example.id
}

output "all_pods_count" {
  description = "Total number of pods in the account"
  value       = length(try(data.runpod_pods.all.pods, []))
}

output "all_volumes_count" {
  description = "Total number of network volumes in the account"
  value       = length(try(data.runpod_network_volumes.all.network_volumes, []))
}

output "all_endpoints_count" {
  description = "Total number of endpoints in the account"
  value       = length(try(data.runpod_endpoints.all.endpoints, []))
}
