package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/samuelbernardolip/gophercloud"
	"github.com/samuelbernardolip/gophercloud/openstack/compute/v2/flavors"
	"github.com/samuelbernardolip/gophercloud/openstack/identity/v3/projects"
	"github.com/samuelbernardolip/gophercloud/pagination"
)

func TestAccComputeV2FlavorAccess_basic(t *testing.T) {
	var flavor flavors.Flavor
	var flavorName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	var project projects.Project
	var projectName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	var flavorAccess flavors.FlavorAccess

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FlavorAccessDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2FlavorAccess_basic(flavorName, projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckComputeV2FlavorExists("openstack_compute_flavor_v2.flavor_1", &flavor),
					testAccCheckComputeV2FlavorAccessExists("openstack_compute_flavor_access_v2.access_1", &flavorAccess),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_flavor_access_v2.access_1", "flavor_id", &flavor.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_flavor_access_v2.access_1", "tenant_id", &project.ID),
				),
			},
		},
	})
}

func testAccCheckComputeV2FlavorAccessDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_compute_flavor_access_v2" {
			continue
		}

		fid, tid, err := parseComputeFlavorAccessId(rs.Primary.ID)
		if err != nil {
			return err
		}

		pager := flavors.ListAccesses(computeClient, fid)
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			accessList, err := flavors.ExtractAccesses(page)
			if err != nil {
				return false, err
			}

			for _, a := range accessList {
				if a.TenantID == tid {
					return false, fmt.Errorf("Flavor Access still exists")
				}
			}

			return true, nil
		})

		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccCheckComputeV2FlavorAccessExists(n string, access *flavors.FlavorAccess) resource.TestCheckFunc {
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

		fid, tid, err := parseComputeFlavorAccessId(rs.Primary.ID)
		if err != nil {
			return err
		}

		pager := flavors.ListAccesses(computeClient, fid)
		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			accessList, err := flavors.ExtractAccesses(page)
			if err != nil {
				return false, err
			}

			for _, a := range accessList {
				if a.TenantID == tid {
					access = &a
					return false, nil
				}
			}

			return true, nil
		})

		return err
	}
}

func testAccComputeV2FlavorAccess_basic(flavorName, tenantName string) string {
	return fmt.Sprintf(`
    resource "openstack_compute_flavor_v2" "flavor_1" {
      name = "%s"
      ram = 512
      vcpus = 1
      disk = 5

      is_public = false
    }

    resource "openstack_identity_project_v3" "project_1" {
      name = "%s"
    }

    resource "openstack_compute_flavor_access_v2" "access_1" {
      flavor_id = "${openstack_compute_flavor_v2.flavor_1.id}"
      tenant_id = "${openstack_identity_project_v3.project_1.id}"
    }
    `, flavorName, tenantName)
}
