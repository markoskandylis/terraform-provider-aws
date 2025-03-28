---
subcategory: "SSM Contacts"
layout: "aws"
page_title: "AWS: aws_ssmcontacts_plan"
description: |-
  Terraform data source for managing an AWS SSM Contact Plan.
---


<!-- Please do not edit this file, it is generated. -->
# Data Source: aws_ssmcontacts_plan

Terraform data source for managing a Plan of an AWS SSM Contact.

## Example Usage

### Basic Usage

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.aws.data_aws_ssmcontacts_plan import DataAwsSsmcontactsPlan
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        DataAwsSsmcontactsPlan(self, "test",
            contact_id="arn:aws:ssm-contacts:us-west-2:123456789012:contact/contactalias"
        )
```

## Argument Reference

The following arguments are required:

* `contact_id` - (Required) The Amazon Resource Name (ARN) of the contact or escalation plan.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `stage` - List of stages. A contact has an engagement plan with stages that contact specified contact channels. An escalation plan uses stages that contact specified contacts.

<!-- cache-key: cdktf-0.20.8 input-92a72d4072ec7f614f162d990ca656e1ab1bfa4324ea98598545cebc4173cb62 -->