package bootfromvolume

import "github.com/samuelbernardolip/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("os-volumes_boot")
}
