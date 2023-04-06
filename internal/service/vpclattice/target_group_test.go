package vpclattice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/networkmanager"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/service/vpclattice"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccVPCLatticeTargetGroup_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_vpclattice_target_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, networkmanager.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCLatticeTargetGroupConfig_base(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func testAccCheckTargetGroupDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).VPCLatticeClient()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_networkmanager_connect_peer" {
				continue
			}

			_, err := vpclattice.FindTargetGroupByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("VPC Lattice Target Group %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccVPCLatticeTargetGroupConfig_base(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
data "aws_region" "current" {}

resource "aws_vpclattice_target_group" "test" {
	name     = %[1]q
	type     = "IP"
	config {
	  port             = 80
	  protocol         = "HTTP"
	  vpc_identifier   = "vpc-00731d8d223dc0c6e"
	  ip_address_type  = "IPV4"
	  protocol_version = "HTTP1"
	  health_check {
		enabled             = true
		health_check_interval_seconds            = 30
		health_check_timeout_seconds             = 5
		healthy_threshold_count   = 2
		unhealthy_threshold_count = 2
		matcher = {
		  httpCode = "200-299"
		}
		path             = "/health"
		port             = 80
		protocol         = "HTTP"
		protocol_version = "HTTP1"
	  }
	}
  }
`, rName))
}
