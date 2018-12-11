package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/samuelbernardolip/gophercloud"
	"github.com/samuelbernardolip/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func TestAccNetworkingV2FloatingIPAssociate_basic(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2FloatingIPAssociate_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2FloatingIPExists(
						"openstack_networking_floatingip_associate_v2.fip_1", &fip),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_floatingip_associate_v2.fip_1", "floating_ip", &fip.FloatingIP),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_floatingip_associate_v2.fip_1", "port_id", &fip.PortID),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2FloatingIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack floating IP: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_floatingip_v2" {
			continue
		}

		fip, err := floatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nil
			}

			return fmt.Errorf("Error retrieving floating IP: %s", err)
		}

		if fip.PortID != "" {
			return fmt.Errorf("floating IP is still associated")
		}
	}

	return nil
}

func testAccCheckNetworkingV2FloatingIPAssociateExists(n string, fip *floatingips.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := floatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("FloatingIP not found")
		}

		*fip = *found

		return nil
	}
}

var testAccNetworkingV2FloatingIPAssociate_basic = fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_router_interface_v2" "router_interface_1" {
  router_id = "${openstack_networking_router_v2.router_1.id}"
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  external_gateway = "%s"
}

resource "openstack_networking_port_v2" "port_1" {
  admin_state_up = "true"
  network_id = "${openstack_networking_subnet_v2.subnet_1.network_id}"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.20"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
  pool = "%s"
}

resource "openstack_networking_floatingip_associate_v2" "fip_1" {
  floating_ip = "${openstack_networking_floatingip_v2.fip_1.address}"
  port_id = "${openstack_networking_port_v2.port_1.id}"
}
`, OS_EXTGW_ID, OS_POOL_NAME)
