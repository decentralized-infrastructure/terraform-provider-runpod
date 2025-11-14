variable "runpod_api_key" {
  description = "RunPod API key"
  type        = string
  sensitive   = true
}

variable "template_id" {
  description = "Template ID for the serverless endpoint"
  type        = string
  default     = "qtis4aranz"  # Example template ID - replace with your own
}
