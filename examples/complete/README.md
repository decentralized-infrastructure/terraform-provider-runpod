# Complete RunPod Terraform Example

This example demonstrates how to use all three RunPod resources:
- Network Volume (persistent storage)
- Pod (GPU compute instance)
- Serverless Endpoint

## Prerequisites

- RunPod API key
- Template ID for serverless endpoint (optional)

## Usage

1. Copy the example variables file:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Edit `terraform.tfvars` and add your RunPod API key:
   ```hcl
   runpod_api_key = "rpa_your_api_key_here"
   template_id    = "your_template_id_here"  # Optional
   ```

3. Initialize Terraform:
   ```bash
   terraform init
   ```

4. Review the planned changes:
   ```bash
   terraform plan
   ```

5. Apply the configuration:
   ```bash
   terraform apply
   ```

6. When done, destroy the resources:
   ```bash
   terraform destroy
   ```

## What This Creates

- **Network Volume**: 10GB persistent storage in US-CA-2
- **GPU Pod**: PyTorch container with RTX 4090 GPU, attached to the network volume
- **Serverless Endpoint**: Auto-scaling endpoint with 0-2 workers

## Outputs

After applying, you'll see:
- Network volume ID
- Pod ID and status
- Pod cost per hour
- Endpoint ID
- Counts of all resources in your account (from data sources)
