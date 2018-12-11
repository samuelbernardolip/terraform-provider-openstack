package openstack

import (
	"testing"

	"github.com/samuelbernardolip/gophercloud"
	"github.com/samuelbernardolip/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/samuelbernardolip/gophercloud/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestExpandNetworkingPortDHCPOptsV2Create(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"ip_version": 4,
		"name":       "A",
		"value":      "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"ip_version": 6,
		"name":       "B",
		"value":      "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	expectedDHCPOptions := []extradhcpopts.CreateExtraDHCPOpt{
		{
			OptName:   "B",
			OptValue:  "false",
			IPVersion: gophercloud.IPVersion(6),
		},
		{
			OptName:   "A",
			OptValue:  "true",
			IPVersion: gophercloud.IPVersion(4),
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Create(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsEmptyV2Create(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	expectedDHCPOptions := []extradhcpopts.CreateExtraDHCPOpt{}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Create(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsV2Update(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"ip_version": 4,
		"name":       "A",
		"value":      "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"ip_version": 6,
		"name":       "B",
		"value":      "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	optsValueTrue := "true"
	optsValueFalse := "false"
	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{
		{
			OptName:   "B",
			OptValue:  &optsValueFalse,
			IPVersion: gophercloud.IPVersion(6),
		},
		{
			OptName:   "A",
			OptValue:  &optsValueTrue,
			IPVersion: gophercloud.IPVersion(4),
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Update(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsEmptyV2Update(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Update(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortDHCPOptsV2Delete(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"ip_version": 4,
		"name":       "A",
		"value":      "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"ip_version": 6,
		"name":       "B",
		"value":      "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{
		{
			OptName: "B",
		},
		{
			OptName: "A",
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Delete(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestFlattenNetworkingPort2DHCPOptionsV2(t *testing.T) {
	dhcpOptions := extradhcpopts.ExtraDHCPOptsExt{
		ExtraDHCPOpts: []extradhcpopts.ExtraDHCPOpt{
			{
				OptName:   "A",
				OptValue:  "true",
				IPVersion: 4,
			},
			{
				OptName:   "B",
				OptValue:  "false",
				IPVersion: 6,
			},
		},
	}

	expectedDHCPOptions := []map[string]interface{}{
		map[string]interface{}{
			"ip_version": 4,
			"name":       "A",
			"value":      "true",
		},
		map[string]interface{}{
			"ip_version": 6,
			"name":       "B",
			"value":      "false",
		},
	}

	actualDHCPOptions := flattenNetworkingPortDHCPOptsV2(dhcpOptions)

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestExpandNetworkingPortAllowedAddressPairsV2(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	addressPairs1 := map[string]interface{}{
		"ip_address":  "192.0.2.1",
		"mac_address": "mac1",
	}
	addressPairs2 := map[string]interface{}{
		"ip_address":  "198.51.100.1",
		"mac_address": "mac2",
	}
	allowedAddressPairs := []map[string]interface{}{addressPairs1, addressPairs2}
	d.Set("allowed_address_pairs", allowedAddressPairs)

	expectedAllowedAddressPairs := []ports.AddressPair{
		{
			IPAddress:  "192.0.2.1",
			MACAddress: "mac1",
		},
		{
			IPAddress:  "198.51.100.1",
			MACAddress: "mac2",
		},
	}

	actualAllowedAddressPairs := expandNetworkingPortAllowedAddressPairsV2(d.Get("allowed_address_pairs").(*schema.Set))

	assert.ElementsMatch(t, expectedAllowedAddressPairs, actualAllowedAddressPairs)
}

func TestFlattenNetworkingPortAllowedAddressPairsV2(t *testing.T) {
	allowedAddressPairs := []ports.AddressPair{
		{
			IPAddress:  "192.0.2.1",
			MACAddress: "mac1",
		},
		{
			IPAddress:  "198.51.100.1",
			MACAddress: "mac2",
		},
	}
	mac := "mac3"

	expectedAllowedAddressPairs := []map[string]interface{}{
		map[string]interface{}{
			"ip_address":  "192.0.2.1",
			"mac_address": "mac1",
		},
		map[string]interface{}{
			"ip_address":  "198.51.100.1",
			"mac_address": "mac2",
		},
	}

	actualAllowedAddressPairs := flattenNetworkingPortAllowedAddressPairsV2(mac, allowedAddressPairs)

	assert.ElementsMatch(t, expectedAllowedAddressPairs, actualAllowedAddressPairs)
}
