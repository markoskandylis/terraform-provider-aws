package kafka_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kafka"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccKafkaVpcConnectionDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)

	var vpcconnection kafka.DescribeVpcConnectionOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_msk_vpc_connection.test"
	resourceName := "aws_msk_vpc_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Kafka)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Kafka),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVPCConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCConnectionDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCConnectionExists(ctx, dataSourceName, &vpcconnection),
					resource.TestCheckResourceAttrPair(resourceName, "arn", dataSourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "authentication", dataSourceName, "authentication"),
					resource.TestCheckResourceAttrPair(resourceName, "target_cluster_arn", dataSourceName, "target_cluster_arn"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id ", dataSourceName, "vpc_id "),
					resource.TestCheckResourceAttrPair(resourceName, "client_subnets", dataSourceName, "client_subnets"),
					resource.TestCheckResourceAttrPair(resourceName, "security_groups", dataSourceName, "security_groups"),
				),
			},
		},
	})
}

func testAccVPCConnectionDataSourceConfig_base(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  count = 3

  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(aws_vpc.test.cidr_block, 8, count.index)

  tags = {
    Name = %[1]q
  }
}


`, rName))
}

func testAccVPCConnectionDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccVPCConnectionDataSourceConfig_base(rName), fmt.Sprintf(`
resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}
resource "aws_msk_vpc_connection" "test" {
  authentication     = "SASL_IAM"
  target_cluster_arn = "arn:aws:kafka:eu-west-2:926562225508:cluster/demo-cluster-1/a7640874-7bdf-4a38-be10-24465449a333-2"
  vpc_id             = aws_vpc.test.id
  client_subnets     = aws_subnet.test[*].id
  security_groups    = [aws_security_group.test.id]
}

data "aws_msk_vpc_connection" "test" {
	arn = aws_msk_vpc_connection.test.arn
}

`, rName))
}
