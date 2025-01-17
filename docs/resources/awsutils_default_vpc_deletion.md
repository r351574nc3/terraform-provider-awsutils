---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "awsutils_default_vpc_deletion Resource - terraform-provider-aws-utils"
subcategory: ""
description: |-
  Deletes the default VPC along with the child resources of the VPC including Subnets, Route Tables, NACLs and Internet
  Gateways in the configured region.
  Best-practices call for not using the default VPC, but rather, creating a new set of VPCs as necessary. AWS Security
  Hub will flag the default VPCs as non-compliant if they aren't configured with best-practices. Rather than jumping
  through hoops, it's easier to delete to default VPCs. This task cannot be accomplished with the official AWS
  Terraform Provider, so this resource is necessary.
  Please note that applying this resource is destructive and nonreversible. This resource is unusual as it will
  DELETE infrastructure when terraform apply is run rather than creating it. This is a permanent
  deletion and nothing will be restored whenterraform destroy is run.
---

# awsutils_default_vpc_deletion

Deletes the default VPC along with the child resources of the VPC including Subnets, Route Tables, NACLs and Internet 
Gateways in the configured region.
		
Best-practices call for not using the default VPC, but rather, creating a new set of VPCs as necessary. AWS Security 
Hub will flag the default VPCs as non-compliant if they aren't configured with best-practices. Rather than jumping 
through hoops, it's easier to delete to default VPCs. This task cannot be accomplished with the official AWS 
Terraform Provider, so this resource is necessary. 
		
Please note that applying this resource is destructive and nonreversible. This resource is unusual as it will 
**DELETE** infrastructure when `terraform apply` is run rather than creating it. This is a permanent 
deletion and nothing will be restored when`terraform destroy` is run.

## Example Usage

```terraform
terraform {
  required_providers {
    awsutils = {
      source = "cloudposse/awsutils"
      # For local development,
      # install the provider on local computer by running `make install` from the root of the repo, and uncomment the 
      # version below
      # version = "9999.99.99"
    }
  }
}

provider "awsutils" {
  region = "us-east-1"
}

# Delete the default VPC in our account/region
resource "awsutils_default_vpc_deletion" "default" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- **id** (String) The ID of the VPC that was deleted.

