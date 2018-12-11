package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/samuelbernardolip/terraform-provider-openstack/openstack"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: openstack.Provider})
}
