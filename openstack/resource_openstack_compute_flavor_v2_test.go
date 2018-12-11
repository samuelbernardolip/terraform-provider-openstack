package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/samuelbernardolip/gophercloud/openstack/compute/v2/flavors"
)

func TestAccComputeV2Flavor_basic(t *testing.T) {
	var flavor flavors.Flavor
	var flavorName = acctest.RandomWithPrefix("tf-acc-flavor")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FlavorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2Flavor_basic(flavorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorExists("openstack_compute_flavor_v2.flavor_1", &flavor),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "ram", "2048"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "vcpus", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "disk", "5"),
				),
			},
		},
	})
}

func TestAccComputeV2Flavor_extraSpecs(t *testing.T) {
	var flavor flavors.Flavor
	var flavorName = acctest.RandomWithPrefix("tf-acc-flavor")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FlavorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2Flavor_extraSpecs_1(flavorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorExists("openstack_compute_flavor_v2.flavor_1", &flavor),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "extra_specs.%", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "extra_specs.hw:cpu_policy", "CPU-POLICY"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "extra_specs.hw:cpu_thread_policy", "CPU-THREAD-POLICY"),
				),
			},
			resource.TestStep{
				Config: testAccComputeV2Flavor_extraSpecs_2(flavorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorExists("openstack_compute_flavor_v2.flavor_1", &flavor),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "extra_specs.%", "1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "extra_specs.hw:cpu_policy", "CPU-POLICY-2"),
				),
			},
		},
	})
}

func testAccCheckComputeV2FlavorDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_compute_flavor_v2" {
			continue
		}

		_, err := flavors.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Flavor still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2FlavorExists(n string, flavor *flavors.Flavor) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		found, err := flavors.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Flavor not found")
		}

		*flavor = *found

		return nil
	}
}

func testAccComputeV2Flavor_basic(flavorName string) string {
	return fmt.Sprintf(`
    resource "openstack_compute_flavor_v2" "flavor_1" {
      name = "%s"
      ram = 2048
      vcpus = 2
      disk = 5

      is_public = true
    }
    `, flavorName)
}

func testAccComputeV2Flavor_extraSpecs_1(flavorName string) string {
	return fmt.Sprintf(`
    resource "openstack_compute_flavor_v2" "flavor_1" {
      name = "%s"
      ram = 2048
      vcpus = 2
      disk = 5

      is_public = true

      extra_specs {
        "hw:cpu_policy" = "CPU-POLICY",
        "hw:cpu_thread_policy" = "CPU-THREAD-POLICY"
      }
    }
    `, flavorName)
}

func testAccComputeV2Flavor_extraSpecs_2(flavorName string) string {
	return fmt.Sprintf(`
    resource "openstack_compute_flavor_v2" "flavor_1" {
      name = "%s"
      ram = 2048
      vcpus = 2
      disk = 5

      is_public = true

      extra_specs {
        "hw:cpu_policy" = "CPU-POLICY-2"
      }
    }
    `, flavorName)
}
