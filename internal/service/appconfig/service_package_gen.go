// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package appconfig

import (
	"context"

	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newResourceEnvironment,
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceConfigurationProfile,
			TypeName: "aws_appconfig_configuration_profile",
		},
		{
			Factory:  DataSourceConfigurationProfiles,
			TypeName: "aws_appconfig_configuration_profiles",
		},
		{
			Factory:  DataSourceEnvironment,
			TypeName: "aws_appconfig_environment",
		},
		{
			Factory:  DataSourceEnvironments,
			TypeName: "aws_appconfig_environments",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceApplication,
			TypeName: "aws_appconfig_application",
			Name:     "Application",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceConfigurationProfile,
			TypeName: "aws_appconfig_configuration_profile",
			Name:     "Connection Profile",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceDeployment,
			TypeName: "aws_appconfig_deployment",
			Name:     "Deployment",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceDeploymentStrategy,
			TypeName: "aws_appconfig_deployment_strategy",
			Name:     "Deployment Strategy",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceExtension,
			TypeName: "aws_appconfig_extension",
			Name:     "Extension",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceExtensionAssociation,
			TypeName: "aws_appconfig_extension_association",
		},
		{
			Factory:  ResourceHostedConfigurationVersion,
			TypeName: "aws_appconfig_hosted_configuration_version",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.AppConfig
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
