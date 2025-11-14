# RunPod Terraform Provider Examples

This directory contains examples demonstrating how to use the RunPod Terraform Provider.

## Prerequisites

- Terraform >= 1.0
- RunPod API key (set via `RUNPOD_API_KEY` environment variable or in provider config)

## Available Examples

### [basic-pod](./basic-pod)
Simple GPU pod creation example with basic configuration including:
- GPU type and data center selection
- Environment variables
- Port exposure
- Storage configuration

### [complete](./complete)
Comprehensive example demonstrating:
- Pod creation with network volume attachment
- Regional constraints for network volumes
- Data sources for listing resources
- Multiple outputs for monitoring

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
   terraform plan
   terraform apply
   ```

4. Clean up when done:
   ```bash
   terraform destroy
   ```

## Notes

- These examples use community cloud instances for cost-effectiveness
- GPU availability varies by region and type
- Network volumes must be in the same region as pods when attached
- State files (`.tfstate`) and variable files (`.tfvars`) are gitignored

