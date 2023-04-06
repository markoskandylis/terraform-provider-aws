package vpclattice_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfvpclattice "github.com/hashicorp/terraform-provider-aws/internal/service/vpclattice"

	// "github.com/hashicorp/terraform-provider-aws/internal/service/vpclattice"

	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccVPCLatticeTargetGroup_basic(t *testing.T) {
	ctx := acctest.Context(t)

	var targetGroup vpclattice.GetTargetGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	// rType := "IP"
	resourceName := "aws_vpclattice_target_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.VPCLatticeEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.VPCLatticeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCLatticeTargetGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "vpc-lattice", regexp.MustCompile("targetgroup/.+$")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVPCLatticeTargetGroup_full(t *testing.T) {
	ctx := acctest.Context(t)

	var targetGroup vpclattice.GetTargetGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_vpclattice_target_group.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.VPCLatticeEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.VPCLatticeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCLatticeTargetGroupConfig_full(rName, "IP"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "config.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "config.0.protocol", "HTTPS"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "vpc-lattice", regexp.MustCompile("targetgroup/.+$")),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			{
				Config: testAccVPCLatticeTargetGroupConfig_full(rName, "INSTANCE"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "config.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "config.0.protocol", "HTTPS"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "vpc-lattice", regexp.MustCompile("targetgroup/.+$")),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			{
				Config: testAccVPCLatticeTargetGroupConfig_full(rName, "ALB"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "config.0.port", "443"),
					resource.TestCheckResourceAttr(resourceName, "config.0.protocol", "HTTPS"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "vpc-lattice", regexp.MustCompile("targetgroup/.+$")),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
		},
	})
}

func testAccVPCLatticeTargetGroupConfig_basic(rName string) string {
	return fmt.Sprintf(`
	data "aws_region" "current" {}
resource "aws_vpclattice_target_group" "test" {
  name = %[1]q
  type = "LAMBDA"
}
`, rName)
}

func testAccVPCLatticeTargetGroupConfig_full(rName, rType string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
data "aws_region" "current" {}

resource "aws_vpclattice_target_group" "test" {
	name     = %[1]q
	type     = %[2]q
	config {
	  port             = 443
	  protocol         = "HTTPS"
	  vpc_identifier   = "vpc-08ae9103897a90142"
	  ip_address_type  = "IPV4"
	  protocol_version = "HTTP1"
	  health_check {
		enabled             = false
		interval            = 30
		timeout             = 5
		healthy_threshold   = 2
		unhealthy_threshold = 2
		matcher 		 		  = "200-299"
		path             		  = "/"
		port             		  = 80
		protocol         		  = "HTTP"
		protocol_version 		  = "HTTP1"
	  }
	}
  }
`, rName, rType))
}

func TestAccVPCLatticeTargetGroup_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var targetGroup1, targetGroup2, targetGroup3 vpclattice.GetTargetGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_vpclattice_target_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.VPCLatticeEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckServiceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccServiceConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccServiceConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupExists(ctx, resourceName, &targetGroup3),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckTargetGroupDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).VPCLatticeClient()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_vpclattice_target_group" {
				continue
			}

			_, err := conn.GetTargetGroup(ctx, &vpclattice.GetTargetGroupInput{
				TargetGroupIdentifier: aws.String(rs.Primary.ID),
			})
			if err != nil {
				var nfe *types.ResourceNotFoundException
				if errors.As(err, &nfe) {
					return nil
				}
				return err
			}

			return create.Error(names.VPCLattice, create.ErrActionCheckingDestroyed, tfvpclattice.ResNameService, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckTargetGroupExists(ctx context.Context, name string, targetGroup *vpclattice.GetTargetGroupOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.VPCLattice, create.ErrActionCheckingExistence, tfvpclattice.ResNameService, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.VPCLattice, create.ErrActionCheckingExistence, tfvpclattice.ResNameService, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).VPCLatticeClient()
		resp, err := conn.GetTargetGroup(ctx, &vpclattice.GetTargetGroupInput{
			TargetGroupIdentifier: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return create.Error(names.VPCLattice, create.ErrActionCheckingExistence, tfvpclattice.ResNameService, rs.Primary.ID, err)
		}

		*targetGroup = *resp

		return nil
	}
}
