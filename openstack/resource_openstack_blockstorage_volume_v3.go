package openstack

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/samuelbernardolip/gophercloud"
	"github.com/samuelbernardolip/gophercloud/openstack/blockstorage/extensions/volumeactions"
	"github.com/samuelbernardolip/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/samuelbernardolip/gophercloud/openstack/compute/v2/extensions/volumeattach"
)

func resourceBlockStorageVolumeV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlockStorageVolumeV3Create,
		Read:   resourceBlockStorageVolumeV3Read,
		Update: resourceBlockStorageVolumeV3Update,
		Delete: resourceBlockStorageVolumeV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"enable_online_resize": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_vol_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"volume_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"consistency_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_replica": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"multiattach": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"attachment": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"device": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourceVolumeV3AttachmentHash,
			},
		},
	}
}

func resourceBlockStorageVolumeV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.blockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	createOpts := &volumes.CreateOpts{
		AvailabilityZone:   d.Get("availability_zone").(string),
		ConsistencyGroupID: d.Get("consistency_group_id").(string),
		Description:        d.Get("description").(string),
		ImageID:            d.Get("image_id").(string),
		Metadata:           resourceContainerMetadataV2(d),
		Name:               d.Get("name").(string),
		Size:               d.Get("size").(int),
		SnapshotID:         d.Get("snapshot_id").(string),
		SourceReplica:      d.Get("source_replica").(string),
		SourceVolID:        d.Get("source_vol_id").(string),
		VolumeType:         d.Get("volume_type").(string),
		Multiattach:        d.Get("multiattach").(bool),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := volumes.Create(blockStorageClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenStack volume: %s", err)
	}
	log.Printf("[INFO] Volume ID: %s", v.ID)

	// Wait for the volume to become available.
	log.Printf(
		"[DEBUG] Waiting for volume (%s) to become available",
		v.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"downloading", "creating"},
		Target:     []string{"available"},
		Refresh:    VolumeV3StateRefreshFunc(blockStorageClient, v.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to become ready: %s",
			v.ID, err)
	}

	// Store the ID now
	d.SetId(v.ID)

	return resourceBlockStorageVolumeV3Read(d, meta)
}

func resourceBlockStorageVolumeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.blockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	v, err := volumes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "volume")
	}

	log.Printf("[DEBUG] Retrieved volume %s: %+v", d.Id(), v)

	d.Set("size", v.Size)
	d.Set("description", v.Description)
	d.Set("availability_zone", v.AvailabilityZone)
	d.Set("name", v.Name)
	d.Set("snapshot_id", v.SnapshotID)
	d.Set("source_vol_id", v.SourceVolID)
	d.Set("volume_type", v.VolumeType)
	d.Set("metadata", v.Metadata)
	d.Set("region", GetRegion(d, config))

	attachments := make([]map[string]interface{}, len(v.Attachments))
	for i, attachment := range v.Attachments {
		attachments[i] = make(map[string]interface{})
		attachments[i]["id"] = attachment.ID
		attachments[i]["instance_id"] = attachment.ServerID
		attachments[i]["device"] = attachment.Device
		log.Printf("[DEBUG] attachment: %v", attachment)
	}
	d.Set("attachment", attachments)

	return nil
}

func resourceBlockStorageVolumeV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.blockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	updateOpts := volumes.UpdateOpts{
		Name:        &name,
		Description: &description,
	}

	if d.HasChange("metadata") {
		updateOpts.Metadata = resourceVolumeMetadataV3(d)
	}

	var v *volumes.Volume
	if d.HasChange("size") {
		v, err = volumes.Get(blockStorageClient, d.Id()).Extract()
		if err != nil {
			return fmt.Errorf(
				"Error extending volume (%s) (%s)",
				d.Id(), err)
		}

		if v.Status == "in-use" {
			if !d.Get("enable_online_resize").(bool) {
				return fmt.Errorf(
					`Error extending volume (%s),
					volume is attached to the instance and
					resizing online is disabled,
					see enable_online_resize option`, d.Id())
			}

			blockStorageClient.Microversion = "3.42"
		}

		extendOpts := volumeactions.ExtendSizeOpts{
			NewSize: d.Get("size").(int),
		}
		err = volumeactions.ExtendSize(blockStorageClient, d.Id(), extendOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf(
				"Error extending volume (%s) size (%s)",
				d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"extending"},
			Target:     []string{"available", "in-use"},
			Refresh:    VolumeV3StateRefreshFunc(blockStorageClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err := stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for volume (%s) to become ready (%s)",
				d.Id(), err)
		}
	}

	_, err = volumes.Update(blockStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenStack volume: %s", err)
	}

	return resourceBlockStorageVolumeV3Read(d, meta)
}

func resourceBlockStorageVolumeV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.blockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	v, err := volumes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "volume")
	}

	// make sure this volume is detached from all instances before deleting
	if len(v.Attachments) > 0 {
		log.Printf("[DEBUG] detaching volumes")
		if computeClient, err := config.computeV2Client(GetRegion(d, config)); err != nil {
			return err
		} else {
			for _, volumeAttachment := range v.Attachments {
				log.Printf("[DEBUG] Attachment: %v", volumeAttachment)
				if err := volumeattach.Delete(computeClient, volumeAttachment.ServerID, volumeAttachment.ID).ExtractErr(); err != nil {
					return err
				}
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"in-use", "attaching", "detaching"},
				Target:     []string{"available"},
				Refresh:    VolumeV3StateRefreshFunc(blockStorageClient, d.Id()),
				Timeout:    10 * time.Minute,
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForState()
			if err != nil {
				return fmt.Errorf(
					"Error waiting for volume (%s) to become available: %s",
					d.Id(), err)
			}
		}
	}

	// It's possible that this volume was used as a boot device and is currently
	// in a "deleting" state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if v.Status != "deleting" {
		if err := volumes.Delete(blockStorageClient, d.Id(), nil).ExtractErr(); err != nil {
			return CheckDeleted(d, err, "volume")
		}
	}

	// Wait for the volume to delete before moving on.
	log.Printf("[DEBUG] Waiting for volume (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "downloading", "available"},
		Target:     []string{"deleted"},
		Refresh:    VolumeV3StateRefreshFunc(blockStorageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceVolumeMetadataV3(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

// VolumeV3StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an OpenStack volume.
func VolumeV3StateRefreshFunc(client *gophercloud.ServiceClient, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := volumes.Get(client, volumeID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return v, "deleted", nil
			}
			return nil, "", err
		}

		if v.Status == "error" {
			return v, v.Status, fmt.Errorf("There was an error creating the volume. " +
				"Please check with your cloud admin or check the Block Storage " +
				"API logs to see why this error occurred.")
		}

		return v, v.Status, nil
	}
}

func resourceVolumeV3AttachmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["instance_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["instance_id"].(string)))
	}
	return hashcode.String(buf.String())
}
