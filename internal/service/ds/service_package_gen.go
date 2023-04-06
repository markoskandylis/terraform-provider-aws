// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package ds

import (
	"context"

	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceDirectory,
			TypeName: "aws_directory_service_directory",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceConditionalForwarder,
			TypeName: "aws_directory_service_conditional_forwarder",
		},
		{
			Factory:  ResourceDirectory,
			TypeName: "aws_directory_service_directory",
			Name:     "Directory",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "id",
			},
		},
		{
			Factory:  ResourceLogSubscription,
			TypeName: "aws_directory_service_log_subscription",
		},
		{
			Factory:  ResourceRadiusSettings,
			TypeName: "aws_directory_service_radius_settings",
		},
		{
			Factory:  ResourceRegion,
			TypeName: "aws_directory_service_region",
		},
		{
			Factory:  ResourceSharedDirectory,
			TypeName: "aws_directory_service_shared_directory",
		},
		{
			Factory:  ResourceSharedDirectoryAccepter,
			TypeName: "aws_directory_service_shared_directory_accepter",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.DS
}

var ServicePackage = &servicePackage{}
