package openstack

import (
	"fmt"
	"testing"
	"time"

	"github.com/samuelbernardolip/gophercloud"
	"github.com/samuelbernardolip/gophercloud/openstack/networking/v2/extensions/fwaas/rules"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccFWRuleV1_basic(t *testing.T) {
	var rule rules.Rule

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFW(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWRuleV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWRuleV1_basic_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV1Exists("openstack_fw_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "action", "deny"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "enabled", "true"),
				),
			},

			resource.TestStep{
				Config: testAccFWRuleV1_basic_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV1Exists("openstack_fw_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "action", "deny"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "description", "Terraform accept test"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "destination_ip_address", "4.3.2.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_port", "444"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "destination_port", "555"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "enabled", "true"),
				),
			},

			resource.TestStep{
				Config: testAccFWRuleV1_basic_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV1Exists("openstack_fw_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "action", "allow"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "description", "Terraform accept test updated"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_ip_address", "1.2.3.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "destination_ip_address", "4.3.2.8"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_port", "666"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "destination_port", "777"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "enabled", "false"),
				),
			},

			resource.TestStep{
				Config: testAccFWRuleV1_basic_4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV1Exists("openstack_fw_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "action", "allow"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "description", "Terraform accept test updated"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_ip_address", "1.2.3.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "destination_ip_address", "4.3.2.8"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_port", "666"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "destination_port", "777"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "enabled", "false"),
				),
			},

			resource.TestStep{
				Config: testAccFWRuleV1_basic_5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV1Exists("openstack_fw_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "action", "allow"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "description", "Terraform accept test updated"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_ip_address", "1.2.3.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "destination_ip_address", "4.3.2.8"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_port", "666"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccFWRuleV1_anyProtocol(t *testing.T) {
	var rule rules.Rule

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckFW(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWRuleV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWRuleV1_anyProtocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV1Exists("openstack_fw_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "action", "allow"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "protocol", "any"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "description", "Allow any protocol"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_ip_address", "192.168.199.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccFWRuleV1_updateName(t *testing.T) {
	var rule rules.Rule

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWRuleV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccFWRuleV1_updateName_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV1Exists("openstack_fw_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "action", "deny"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "description", "Terraform accept test"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "destination_ip_address", "4.3.2.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_port", "444"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "destination_port", "555"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "enabled", "true"),
				),
			},

			resource.TestStep{
				Config: testAccFWRuleV1_updateName_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV1Exists("openstack_fw_rule_v1.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "name", "updated_rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "action", "deny"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "description", "Terraform accept test"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "destination_ip_address", "4.3.2.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "source_port", "444"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "destination_port", "555"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v1.rule_1", "enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckFWRuleV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_firewall_rule" {
			continue
		}
		_, err = rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Firewall rule (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWRuleV1Exists(n string, rule *rules.Rule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		var found *rules.Rule
		for i := 0; i < 5; i++ {
			// Firewall rule creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = rules.Get(networkingClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(gophercloud.ErrDefault404); ok {
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		*rule = *found

		return nil
	}
}

const testAccFWRuleV1_basic_1 = `
resource "openstack_fw_rule_v1" "rule_1" {
	name = "rule_1"
	protocol = "udp"
	action = "deny"
}
`

const testAccFWRuleV1_basic_2 = `
resource "openstack_fw_rule_v1" "rule_1" {
	name = "rule_1"
	description = "Terraform accept test"
	protocol = "udp"
	action = "deny"
	ip_version = 4
	source_ip_address = "1.2.3.4"
	destination_ip_address = "4.3.2.0/24"
	source_port = "444"
	destination_port = "555"
	enabled = true
}
`

const testAccFWRuleV1_basic_3 = `
resource "openstack_fw_rule_v1" "rule_1" {
	name = "rule_1"
	description = "Terraform accept test updated"
	protocol = "tcp"
	action = "allow"
	ip_version = 4
	source_ip_address = "1.2.3.0/24"
	destination_ip_address = "4.3.2.8"
	source_port = "666"
	destination_port = "777"
	enabled = false
}
`

const testAccFWRuleV1_basic_4 = `
resource "openstack_fw_rule_v1" "rule_1" {
	name = "rule_1"
	description = "Terraform accept test updated"
	protocol = "udp"
	action = "allow"
	ip_version = 4
	source_ip_address = "1.2.3.0/24"
	destination_ip_address = "4.3.2.8"
	source_port = "666"
	destination_port = "777"
	enabled = false
}
`

const testAccFWRuleV1_basic_5 = `
resource "openstack_fw_rule_v1" "rule_1" {
	name = "rule_1"
	description = "Terraform accept test updated"
	protocol = "udp"
	action = "allow"
	ip_version = 4
	source_ip_address = "1.2.3.0/24"
	destination_ip_address = "4.3.2.8"
	source_port = "666"
	enabled = false
}
`

const testAccFWRuleV1_anyProtocol = `
resource "openstack_fw_rule_v1" "rule_1" {
	name = "rule_1"
	description = "Allow any protocol"
	protocol = "any"
	action = "allow"
	ip_version = 4
	source_ip_address = "192.168.199.0/24"
	enabled = true
}
`

const testAccFWRuleV1_updateName_1 = `
resource "openstack_fw_rule_v1" "rule_1" {
	name = "rule_1"
	description = "Terraform accept test"
	protocol = "udp"
	action = "deny"
	ip_version = 4
	source_ip_address = "1.2.3.4"
	destination_ip_address = "4.3.2.0/24"
	source_port = "444"
	destination_port = "555"
	enabled = true
}
`

const testAccFWRuleV1_updateName_2 = `
resource "openstack_fw_rule_v1" "rule_1" {
	name = "updated_rule_1"
	description = "Terraform accept test"
	protocol = "udp"
	action = "deny"
	ip_version = 4
	source_ip_address = "1.2.3.4"
	destination_ip_address = "4.3.2.0/24"
	source_port = "444"
	destination_port = "555"
	enabled = true
}
`
