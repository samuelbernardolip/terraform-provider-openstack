package openstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/samuelbernardolip/gophercloud/openstack/compute/v2/extensions/keypairs"
)

func dataSourceComputeKeypairV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceComputeKeypairV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			// computed-only
			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeKeypairV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	name := d.Get("name").(string)
	kp, err := keypairs.Get(computeClient, name).Extract()
	if err != nil {
		return fmt.Errorf("Error retrieving openstack_compute_keypair_v2 %s: %s", name, err)
	}

	d.SetId(name)

	log.Printf("[DEBUG] Retrieved openstack_compute_keypair_v2 %s: %#v", d.Id(), kp)

	d.Set("fingerprint", kp.Fingerprint)
	d.Set("public_key", kp.PublicKey)
	d.Set("region", GetRegion(d, config))

	return nil
}
