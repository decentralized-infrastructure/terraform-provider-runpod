# Terraform Provider Documentation

This directory will contain generated documentation for the RunPod Terraform Provider.

Documentation is automatically generated using `tfplugindocs`:

```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
tfplugindocs generate
```

The generated documentation will include:
- Provider configuration
- Resource schemas
- Data source schemas
- Examples from the `examples/` directory
