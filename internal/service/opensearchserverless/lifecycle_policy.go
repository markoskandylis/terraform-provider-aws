// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package opensearchserverless

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearchserverless"
	awstypes "github.com/aws/aws-sdk-go-v2/service/opensearchserverless/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/fwdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// Function annotations are used for resource registration to the Provider. DO NOT EDIT.
// @FrameworkResource(name="Lifecycle Policy")
func newResourceLifecyclePolicy(_ context.Context) (resource.ResourceWithConfigure, error) {
	return &resourceLifecyclePolicy{}, nil
}

const (
	ResNameLifecyclePolicy = "Lifecycle Policy"
)

type resourceLifecyclePolicy struct {
	framework.ResourceWithConfigure
	framework.WithTimeouts
}

func (r *resourceLifecyclePolicy) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "aws_opensearchserverless_lifecycle_policy"
}

func (r *resourceLifecyclePolicy) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 1000),
				},
			},
			"id": framework.IDAttribute(),
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 32),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 20480),
				},
			},
			"policy_version": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					enum.FrameworkValidate[awstypes.LifecyclePolicyType](),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *resourceLifecyclePolicy) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	conn := r.Meta().OpenSearchServerlessClient(ctx)

	var plan resourceLifecyclePolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := &opensearchserverless.CreateLifecyclePolicyInput{
		ClientToken: aws.String(id.UniqueId()),
		Name:        flex.StringFromFramework(ctx, plan.Name),
		Policy:      flex.StringFromFramework(ctx, plan.Policy),
		Type:        awstypes.LifecyclePolicyType(plan.Type.ValueString()),
	}

	if !plan.Description.IsNull() {
		in.Description = flex.StringFromFramework(ctx, plan.Description)
	}

	out, err := conn.CreateLifecyclePolicy(ctx, in)
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.OpenSearchServerless, create.ErrActionCreating, ResNameLifecyclePolicy, plan.Name.ValueString(), err),
			err.Error(),
		)
		return
	}
	if out == nil || out.LifecyclePolicyDetail == nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.OpenSearchServerless, create.ErrActionCreating, ResNameLifecyclePolicy, plan.Name.ValueString(), nil),
			errors.New("empty output").Error(),
		)
		return
	}

	state := plan
	resp.Diagnostics.Append(state.refreshFromOutput(ctx, out.LifecyclePolicyDetail)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceLifecyclePolicy) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	conn := r.Meta().OpenSearchServerlessClient(ctx)

	var state resourceLifecyclePolicyData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := findLifecyclePolicyByNameAndType(ctx, conn, state.ID.ValueString(), state.Type.ValueString())

	if tfresource.NotFound(err) {
		resp.Diagnostics.Append(fwdiag.NewResourceNotFoundWarningDiagnostic(err))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.OpenSearchServerless, create.ErrActionReading, ResNameLifecyclePolicy, state.ID.String(), err),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(state.refreshFromOutput(ctx, out)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceLifecyclePolicy) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	conn := r.Meta().OpenSearchServerlessClient(ctx)

	var plan, state resourceLifecyclePolicyData
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Description.Equal(state.Description) || !plan.Policy.Equal(state.Policy) {
		in := &opensearchserverless.UpdateLifecyclePolicyInput{
			ClientToken:   aws.String(id.UniqueId()),
			Name:          flex.StringFromFramework(ctx, plan.Name),
			PolicyVersion: flex.StringFromFramework(ctx, state.PolicyVersion),
			Type:          awstypes.LifecyclePolicyType(plan.Type.ValueString()),
		}

		if !plan.Policy.Equal(state.Policy) {
			in.Policy = flex.StringFromFramework(ctx, plan.Policy)
		}

		if !plan.Description.Equal(state.Description) {
			in.Description = flex.StringFromFramework(ctx, plan.Description)
		}

		out, err := conn.UpdateLifecyclePolicy(ctx, in)
		if err != nil {
			resp.Diagnostics.AddError(
				create.ProblemStandardMessage(names.OpenSearchServerless, create.ErrActionUpdating, ResNameLifecyclePolicy, plan.ID.ValueString(), err),
				err.Error(),
			)
			return
		}
		if out == nil || out.LifecyclePolicyDetail == nil {
			resp.Diagnostics.AddError(
				create.ProblemStandardMessage(names.OpenSearchServerless, create.ErrActionUpdating, ResNameLifecyclePolicy, plan.ID.ValueString(), nil),
				errors.New("empty output").Error(),
			)
			return
		}

		resp.Diagnostics.Append(state.refreshFromOutput(ctx, out.LifecyclePolicyDetail)...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *resourceLifecyclePolicy) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	conn := r.Meta().OpenSearchServerlessClient(ctx)

	var state resourceLifecyclePolicyData
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := &opensearchserverless.DeleteLifecyclePolicyInput{
		ClientToken: aws.String(id.UniqueId()),
		Name:        flex.StringFromFramework(ctx, state.Name),
		Type:        awstypes.LifecyclePolicyType(state.Type.ValueString()),
	}

	_, err := conn.DeleteLifecyclePolicy(ctx, in)

	if err != nil {
		var nfe *awstypes.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return
		}
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.OpenSearchServerless, create.ErrActionDeleting, ResNameLifecyclePolicy, state.ID.String(), err),
			err.Error(),
		)
		return
	}
}

func (r *resourceLifecyclePolicy) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, idSeparator)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		err := fmt.Errorf("unexpected format for ID (%[1]s), expected lifecycle-policy-name%[2]slifecycle-policy-type", req.ID, idSeparator)
		resp.Diagnostics.AddError(fmt.Sprintf("importing %s (%s)", ResNameLifecyclePolicy, req.ID), err.Error())
		return
	}

	state := resourceLifecyclePolicyData{
		ID:   types.StringValue(parts[0]),
		Name: types.StringValue(parts[0]),
		Type: types.StringValue(parts[1]),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// refreshFromOutput writes state data from an AWS response object
func (rd *resourceLifecyclePolicyData) refreshFromOutput(ctx context.Context, out *awstypes.LifecyclePolicyDetail) diag.Diagnostics {
	var diags diag.Diagnostics

	if out == nil {
		return diags
	}

	rd.ID = flex.StringToFramework(ctx, out.Name)
	rd.Description = flex.StringToFramework(ctx, out.Description)
	rd.Name = flex.StringToFramework(ctx, out.Name)
	rd.Type = flex.StringValueToFramework(ctx, out.Type)
	rd.PolicyVersion = flex.StringToFramework(ctx, out.PolicyVersion)

	policyBytes, err := out.Policy.MarshalSmithyDocument()
	if err != nil {
		diags.AddError(fmt.Sprintf("refreshing state for %s (%s)", ResNameLifecyclePolicy, rd.Name), err.Error())
		return diags
	}

	p := string(policyBytes)

	rd.Policy = flex.StringToFramework(ctx, &p)

	return diags
}

type resourceLifecyclePolicyData struct {
	Description   types.String `tfsdk:"description"`
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Policy        types.String `tfsdk:"policy"`
	PolicyVersion types.String `tfsdk:"policy_version"`
	Type          types.String `tfsdk:"type"`
}
