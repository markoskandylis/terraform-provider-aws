---
subcategory: "VPC Lattice"
layout: "aws"
page_title: "AWS: aws_vpclattice_service_network"
description: |-
  Terraform resource for managing an AWS VPC Lattice Service Network.
---

# Resource: aws_vpclattice_service_network

Terraform resource for managing an AWS VPC Lattice Service Network.

## Example Usage

### Basic Usage

```terraform
resource "aws_vpclattice_service_network" "example" {
<<<<<<< HEAD
=======
  name      = "example"
  auth_type = "AWS_IAM"
>>>>>>> origin/main
}
```

## Argument Reference

The following arguments are required:

<<<<<<< HEAD
* `example_arg` - (Required) Concise argument description. Do not begin the description with "An", "The", "Defines", "Indicates", or "Specifies," as these are verbose. In other words, "Indicates the amount of storage," can be rewritten as "Amount of storage," without losing any information.

The following arguments are optional:

* `optional_arg` - (Optional) Concise argument description. Do not begin the description with "An", "The", "Defines", "Indicates", or "Specifies," as these are verbose. In other words, "Indicates the amount of storage," can be rewritten as "Amount of storage," without losing any information.
=======
* `name` - (Required) Name of the service network

The following arguments are optional:

* `auth_type` - (Optional) Type of IAM policy. Either `NONE` or `AWS_IAM`.
* `tags` - (Optional) Key-value mapping of resource tags. If configured with a provider [`default_tags` configuration block](/docs/providers/aws/index.html#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.
>>>>>>> origin/main

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

<<<<<<< HEAD
* `arn` - ARN of the Service Network. Do not begin the description with "An", "The", "Defines", "Indicates", or "Specifies," as these are verbose. In other words, "Indicates the amount of storage," can be rewritten as "Amount of storage," without losing any information.
* `example_attribute` - Concise description. Do not begin the description with "An", "The", "Defines", "Indicates", or "Specifies," as these are verbose. In other words, "Indicates the amount of storage," can be rewritten as "Amount of storage," without losing any information.

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

* `create` - (Default `60m`)
* `update` - (Default `180m`)
* `delete` - (Default `90m`)

## Import

VPC Lattice Service Network can be imported using the `example_id_arg`, e.g.,

```
$ terraform import aws_vpclattice_service_network.example rft-8012925589
=======
* `arn` - ARN of the Service Network.
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](/docs/providers/aws/index.html#default_tags-configuration-block).

## Import

VPC Lattice Service Network can be imported using the `id`, e.g.,

```
$ terraform import aws_vpclattice_service_network.example sn-0158f91c1e3358dba
>>>>>>> origin/main
```
