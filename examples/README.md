# RunPod Terraform Provider Examples

This directory contains examples demonstrating how to use the RunPod Terraform Provider.

## Prerequisites

- Terraform >= 1.0
- RunPod API key (set via `RUNPOD_API_KEY` environment variable)

## Available Examples

- `basic-pod/` - Simple GPU pod creation example
- More examples coming soon...

## Quick Start

1. Set your API key:
   ```bash
   export RUNPOD_API_KEY="your-api-key-here"
   ```

2. Choose an example and navigate to its directory:
   ```bash
   cd basic-pod
   ```

3. Run Terraform:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Notes

These examples are for demonstration purposes. Adjust configurations according to your specific requirements.
