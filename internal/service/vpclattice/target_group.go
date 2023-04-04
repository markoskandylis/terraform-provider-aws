package vpclattice

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"

	// "github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_vpclattice_target_group")
func ResourceTargetGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceTargetGroupCreate,
		ReadWithoutTimeout:   resourceTargetGroupRead,
		UpdateWithoutTimeout: resourceTargetGroupUpdate,
		DeleteWithoutTimeout: resourceTargetGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"interval": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"healthy_threshold": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"unhealthy_threshold": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"matcher": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 65535),
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"protocol_version": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
										StateFunc: func(v interface{}) string {
											return strings.ToUpper(v.(string))
										},
										ValidateFunc: validation.StringInSlice([]string{
											"GRPC",
											"HTTP1",
											"HTTP2",
										}, true),
									},
								},
							},
						},
						"ip_address_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
							StateFunc: func(v interface{}) string {
								return strings.ToUpper(v.(string))
							},
							ValidateFunc: validation.StringInSlice([]string{
								"GRPC",
								"HTTP1",
								"HTTP2",
							}, true),
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},
	}
}

const (
	ResNameTargetGroup = "Target Group"
)

func resourceTargetGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient()

	in := &vpclattice.CreateTargetGroupInput{
		Name: aws.String(d.Get("name").(string)),
		Type: types.TargetGroupType(d.Get("type").(string)),
	}

	if v, ok := d.GetOk("config"); ok && len(v.([]interface{})) > 0 {
		config := expandConfigAttributes(v.([]interface{})[0].(map[string]interface{}))
		in.Config = &types.TargetGroupConfig{
			Port:            config.Port,
			Protocol:        config.Protocol,
			VpcIdentifier:   config.VpcIdentifier,
			IpAddressType:   config.IpAddressType,
			ProtocolVersion: config.ProtocolVersion,
			HealthCheck:     config.HealthCheck,
		}
	}

	// in.Tags = d.Get("tags")

	out, err := conn.CreateTargetGroup(ctx, in)
	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionCreating, ResNameTargetGroup, d.Get("name").(string), err)
	}

	if out == nil || out.Config == nil {
		return create.DiagError(names.VPCLattice, create.ErrActionCreating, ResNameTargetGroup, d.Get("name").(string), errors.New("empty output"))
	}

	d.SetId(aws.ToString(out.Id))

	if _, err := waitTargetGroupCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionWaitingForCreation, ResNameTargetGroup, d.Id(), err)
	}

	return resourceTargetGroupRead(ctx, d, meta)
}

func resourceTargetGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient()

	out, err := findTargetGroupByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] VpcLattice TargetGroup (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionReading, ResNameTargetGroup, d.Id(), err)
	}

	d.Set("arn", out.Arn)
	d.Set("name", out.Name)

	if err := d.Set("config", flattenTargetGroupConfig(out.Config)); err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionSetting, ResNameTargetGroup, d.Id(), err)
	}

	// defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	// ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig
	// tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	// if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
	// 	return create.DiagError(names.VpcLattice, create.ErrActionSetting, ResNameTargetGroup, d.Id(), err)
	// }

	// if err := d.Set("tags_all", tags.Map()); err != nil {
	// 	return create.DiagError(names.VpcLattice, create.ErrActionSetting, ResNameTargetGroup, d.Id(), err)
	// }

	return nil
}

func resourceTargetGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient()

	update := false

	in := &vpclattice.UpdateTargetGroupInput{
		TargetGroupIdentifier: aws.String(d.Id()),
	}

	if d.HasChange("name") {
		update = true
		in. = aws.String(d.Get("name").(string))
	}

	if d.HasChange("type") {
		update = true
		in.Type = types.TargetGroupType(d.Get("type").(string))
	}

	if d.HasChange("config") {
		update = true
		if v, ok := d.GetOk("config"); ok && len(v.([]interface{})) > 0 {
			config := expandConfigAttributes(v.([]interface{})[0].(map[string]interface{}))
			in.Config = &types.TargetGroupConfig{
				Port:            config.Port,
				Protocol:        config.Protocol,
				VpcIdentifier:   config.VpcIdentifier,
				IpAddressType:   config.IpAddressType,
				ProtocolVersion: config.ProtocolVersion,
				HealthCheck:     config.HealthCheck,
			}
		}
	}

	if update {
		log.Printf("[DEBUG] Updating VpcLattice TargetGroup (%s): %#v", d.Id(), in)
		out, err := conn.UpdateTargetGroup(ctx, in)
		if err != nil {
			return create.DiagError(names.VPCLattice, create.ErrActionUpdating, ResNameTargetGroup, d.Id(), err)
		}

		if _, err := waitTargetGroupUpdated(ctx, conn, aws.ToString(out.OperationId), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return create.DiagError(names.VPCLattice, create.ErrActionWaitingForUpdate, ResNameTargetGroup, d.Id(), err)
		}
	}

	return resourceTargetGroupRead(ctx, d, meta)
}

func resourceTargetGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).VPCLatticeClient()

	log.Printf("[INFO] Deleting VpcLattice TargetGroup %s", d.Id())

	_, err := conn.DeleteTargetGroup(ctx, &vpclattice.DeleteTargetGroupInput{
		Id: aws.String(d.Id()),
	})
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil
		}

		return create.DiagError(names.VPCLattice, create.ErrActionDeleting, ResNameTargetGroup, d.Id(), err)
	}

	// TIP: -- 4. Use a waiter to wait for delete to complete
	if _, err := waitTargetGroupDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return create.DiagError(names.VPCLattice, create.ErrActionWaitingForDeletion, ResNameTargetGroup, d.Id(), err)
	}

	// TIP: -- 5. Return nil
	return nil
}

const (
	statusChangePending = "Pending"
	statusDeleting      = "Deleting"
	statusNormal        = "Normal"
	statusUpdated       = "Updated"
)

func waitTargetGroupCreated(ctx context.Context, conn *vpclattice.Client, id string, timeout time.Duration) (*vpclattice.CreateTargetGroupOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{},
		Target:                    []string{statusNormal},
		Refresh:                   statusTargetGroup(ctx, conn, id),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*vpclattice.CreateTargetGroupOutput); ok {
		return out, err
	}

	return nil, err
}

func waitTargetGroupUpdated(ctx context.Context, conn *vpclattice.Client, id string, timeout time.Duration) (*vpclattice.UpdateTargetGroupOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{statusChangePending},
		Target:                    []string{statusUpdated},
		Refresh:                   statusTargetGroup(ctx, conn, id),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*vpclattice.UpdateTargetGroupOutput); ok {
		return out, err
	}

	return nil, err
}

func waitTargetGroupDeleted(ctx context.Context, conn *vpclattice.Client, id string, timeout time.Duration) (*vpclattice.DeleteTargetGroupOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{statusDeleting, statusNormal},
		Target:  []string{},
		Refresh: statusTargetGroup(ctx, conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*vpclattice.DeleteTargetGroupOutput); ok {
		return out, err
	}

	return nil, err
}

func statusTargetGroup(ctx context.Context, conn *vpclattice.Client, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		out, err := findTargetGroupByID(ctx, conn, id)
		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return out, string(out.Status), nil
	}
}

func findTargetGroupByID(ctx context.Context, conn *vpclattice.Client, id string) (*vpclattice.GetTargetGroupOutput, error) {
	in := &vpclattice.GetTargetGroupInput{
		TargetGroupIdentifier: aws.String(id),
	}
	out, err := conn.GetTargetGroup(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &resource.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.Id == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

func findTargetGroupByARN(ctx context.Context, conn *vpclattice.Client, arn string) (*vpclattice.GetTargetGroupOutput, error) {
	in := &vpclattice.GetTargetGroupInput{
		TargetGroupIdentifier: aws.String(arn),
	}
	out, err := conn.GetTargetGroup(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &resource.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.Id == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

func findTargetGroupByName(ctx context.Context, conn *vpclattice.Client, name string) (*vpclattice.GetTargetGroupOutput, error) {
	in := &vpclattice.GetTargetGroupInput{
		TargetGroupIdentifier: aws.String(name),
	}
	out, err := conn.GetTargetGroup(ctx, in)
	if err != nil {
		var nfe *types.ResourceNotFoundException
		if errors.As(err, &nfe) {
			return nil, &resource.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil || out.Id == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}

func flattenTargetGroupConfig(apiObject *types.TargetGroupConfig) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{
		"port":           aws.Int32(*apiObject.Port),
		"protocol":       string(apiObject.Protocol),
		"vpc_identifier": aws.String(*apiObject.VpcIdentifier),
	}

	if apiObject.IpAddressType != "" {
		m["ip_address_type"] = string(apiObject.IpAddressType)
	}

	if apiObject.ProtocolVersion != "" {
		m["protocol_version"] = string(apiObject.ProtocolVersion)
	}

	if apiObject.HealthCheck != nil {
		m["health_check"] = flattenHealthCheckConfig(apiObject.HealthCheck)
	}

	return m
}

func flattenHealthCheckConfig(apiObject *types.HealthCheckConfig) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	m := map[string]interface{}{
		"enable":              aws.Bool(*apiObject.Enabled),
		"interval":            aws.Int32(*apiObject.HealthCheckIntervalSeconds),
		"timeout":             aws.Int32(*apiObject.HealthCheckTimeoutSeconds),
		"healthy_threshold":   aws.Int32(*apiObject.HealthyThresholdCount),
		"unhealthy_threshold": aws.Int32(*apiObject.UnhealthyThresholdCount),
		"path":                aws.String(*apiObject.Path),
		"port":                aws.Int32(*apiObject.Port),
		"protocol":            string(apiObject.Protocol),
		"protocol_version":    string(apiObject.ProtocolVersion),
	}

	if matcher, ok := apiObject.Matcher.(*types.MatcherMemberHttpCode); ok {
		m["matcher"] = aws.ToString(&matcher.Value)
	}

	return m
}

func expandTargetGroupAttributes(tfMap map[string]interface{}) *types.TargetGroup {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.TargetGroup{}

	if v, ok := tfMap["name"].(string); ok {
		apiObject.Name = aws.String(v)
	}

	if v, ok := tfMap["type"].(string); ok {
		targetGroupType := types.TargetGroupTypeEnum(v)
		apiObject.Type = targetGroupType
	}

	if v, ok := tfMap["tags"].(map[string]interface{}); ok {
		tags := tftags.Expand(v)
		apiObject.Tags = tags
	}

	if v, ok := tfMap["config"].([]interface{}); ok && len(v) > 0 {
		config := expandConfigAttributes(v[0].(map[string]interface{}))
		apiObject.TargetGroupConfig = config
	}

	return apiObject
}

func expandConfigAttributes(tfMap map[string]interface{}) *types.TargetGroupConfig {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.TargetGroupConfig{}

	if v, ok := tfMap["port"].(int); ok {
		apiObject.Port = aws.Int32(int32(v))
	}

	if v, ok := tfMap["protocol"].(string); ok {
		protocol := types.TargetGroupProtocol(v)
		apiObject.Protocol = protocol
	}

	if v, ok := tfMap["vpc_identifier"].(string); ok {
		apiObject.VpcIdentifier = aws.String(v)
	}

	if v, ok := tfMap["health_check"].(map[string]interface{}); ok {
		hc := expandHealthCheckConfigAttributes(v)
		apiObject.HealthCheck = hc
	}

	if v, ok := tfMap["ip_address_type"].(string); ok {
		ipAddressType := types.IpAddressType(v)
		apiObject.IpAddressType = ipAddressType
	}

	if v, ok := tfMap["protocol_version"].(string); ok {
		protocolVersion := types.TargetGroupProtocolVersion(v)
		apiObject.ProtocolVersion = protocolVersion
	}

	return apiObject
}

func expandHealthCheckConfigAttributes(tfMap map[string]interface{}) *types.HealthCheckConfig {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.HealthCheckConfig{}

	if v, ok := tfMap["enable"].(bool); ok {
		apiObject.Enabled = aws.Bool(v)
	}

	if v, ok := tfMap["interval"].(int); ok {
		apiObject.HealthCheckIntervalSeconds = aws.Int32(int32(v))
	}

	if v, ok := tfMap["timeout"].(int); ok {
		apiObject.HealthCheckTimeoutSeconds = aws.Int32(int32(v))
	}

	if v, ok := tfMap["healthy_threshold"].(int); ok {
		apiObject.HealthyThresholdCount = aws.Int32(int32(v))
	}

	if v, ok := tfMap["unhealthy_threshold"].(int); ok {
		apiObject.UnhealthyThresholdCount = aws.Int32(int32(v))
	}

	if v, ok := tfMap["path"].(string); ok {
		apiObject.Path = aws.String(v)
	}

	if v, ok := tfMap["port"].(int); ok {
		apiObject.Port = aws.Int32(int32(v))
	}

	if v, ok := tfMap["protocol"].(string); ok {
		apiObject.Protocol = types.TargetGroupProtocol(v)
	}

	if v, ok := tfMap["protocol_version"].(string); ok {
		apiObject.ProtocolVersion = types.HealthCheckProtocolVersion(v)
	}

	if v, ok := tfMap["matcher"].(map[string]interface{}); ok {
		matcher := &types.MatcherMemberHttpCode{}
		if httpCode, ok := v["httpCode"].(string); ok {
			matcher.Value = httpCode
		}
		apiObject.Matcher = matcher
	}

	return apiObject
}
