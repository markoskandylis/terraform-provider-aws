---
subcategory: "QuickSight"
layout: "aws"
page_title: "AWS: aws_quicksight_folder_membership"
description: |-
  Terraform resource for managing an AWS QuickSight Folder Membership.
---


<!-- Please do not edit this file, it is generated. -->
# Resource: aws_quicksight_folder_membership

Terraform resource for managing an AWS QuickSight Folder Membership.

## Example Usage

### Basic Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { Token, TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { QuicksightFolderMembership } from "./.gen/providers/aws/quicksight-folder-membership";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new QuicksightFolderMembership(this, "example", {
      folderId: Token.asString(awsQuicksightFolderExample.folderId),
      memberId: Token.asString(awsQuicksightDataSetExample.dataSetId),
      memberType: "DATASET",
    });
  }
}

```

## Argument Reference

The following arguments are required:

* `folderId` - (Required, Forces new resource) Identifier for the folder.
* `memberId` - (Required, Forces new resource) ID of the asset (the dashboard, analysis, or dataset).
* `memberType` - (Required, Forces new resource) Type of the member. Valid values are `ANALYSIS`, `DASHBOARD`, and `DATASET`.

The following arguments are optional:

* `awsAccountId` - (Optional, Forces new resource) AWS account ID.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - A comma-delimited string joining AWS account ID, folder ID, member type, and member ID.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import QuickSight Folder Membership using the AWS account ID, folder ID, member type, and member ID separated by commas (`,`). For example:

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { QuicksightFolderMembership } from "./.gen/providers/aws/quicksight-folder-membership";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    QuicksightFolderMembership.generateConfigForImport(
      this,
      "example",
      "123456789012,example-folder,DATASET,example-dataset"
    );
  }
}

```

Using `terraform import`, import QuickSight Folder Membership using the AWS account ID, folder ID, member type, and member ID separated by commas (`,`). For example:

```console
% terraform import aws_quicksight_folder_membership.example 123456789012,example-folder,DATASET,example-dataset
```

<!-- cache-key: cdktf-0.20.8 input-d86c3016c2fb02ce6bbff475bfb57b5c2ef1344603d2a5e1e68e88605800f66c -->