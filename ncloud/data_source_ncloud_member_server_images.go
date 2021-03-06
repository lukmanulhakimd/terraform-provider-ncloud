package ncloud

import (
	"fmt"
	"regexp"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudMemberServerImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudMemberServerImagesRead,

		Schema: map[string]*schema.Schema{
			"member_server_image_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew:     true,
				ValidateFunc: validateRegexp,
				Description:  "A regex string to apply to the member server image list returned by ncloud",
			},
			"member_server_image_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of member server images to view",
			},
			"platform_type_code_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of platform codes of server images to view",
			},
			"region_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region code. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_no"},
			},
			"region_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region number. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_code"},
			},

			"member_server_images": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Member server image list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_server_image_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member server image no",
						},
						"member_server_image_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member server image name",
						},
						"member_server_image_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member server image description",
						},
						"original_server_instance_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original server instance no",
						},
						"original_server_product_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original server product code",
						},
						"original_server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original server name",
						},
						"original_base_block_storage_disk_type": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        commonCodeSchemaResource,
							Description: "Original base block storage disk type",
						},
						"original_server_image_product_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original server image product code",
						},
						"original_os_information": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original os information",
						},
						"original_server_image_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original server image name",
						},
						"member_server_image_status_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member server image status name",
						},
						"member_server_image_status": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        commonCodeSchemaResource,
							Description: "Member server image status",
						},
						"member_server_image_operation": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        commonCodeSchemaResource,
							Description: "Member server image operation",
						},
						"member_server_image_platform_type": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        commonCodeSchemaResource,
							Description: "Member server image platform type",
						},
						"create_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation date of the member server image",
						},
						"region": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        regionSchemaResource,
							Description: "Region info",
						},
						"member_server_image_block_storage_total_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Member server image block storage total rows",
						},
						"member_server_image_block_storage_total_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Member server image block storage total size",
						},
					},
				},
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of file that can save data source after running `terraform plan`.",
			},
		},
	}
}

func dataSourceNcloudMemberServerImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return err
	}
	reqParams := server.GetMemberServerImageListRequest{
		MemberServerImageNoList: expandStringInterfaceList(d.Get("member_server_image_no_list").([]interface{})),
		PlatformTypeCodeList:    expandStringInterfaceList(d.Get("platform_type_code_list").([]interface{})),
		RegionNo:                regionNo,
	}

	logCommonRequest("GetMemberServerImageList", reqParams)

	resp, err := client.server.V2Api.GetMemberServerImageList(&reqParams)
	if err != nil {
		logErrorResponse("GetMemberServerImageList", err, reqParams)
		return err
	}
	logCommonResponse("GetMemberServerImageList", GetCommonResponse(resp))

	allMemberServerImages := resp.MemberServerImageList
	var filteredMemberServerImages []*server.MemberServerImage
	nameRegex, nameRegexOk := d.GetOk("member_server_image_name_regex")
	if nameRegexOk {
		r := regexp.MustCompile(nameRegex.(string))
		for _, memberServerImage := range allMemberServerImages {
			if r.MatchString(ncloud.StringValue(memberServerImage.MemberServerImageName)) {
				filteredMemberServerImages = append(filteredMemberServerImages, memberServerImage)
			}
		}
	} else {
		filteredMemberServerImages = allMemberServerImages[:]
	}

	if len(filteredMemberServerImages) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	return memberServerImagesAttributes(d, filteredMemberServerImages)
}

func memberServerImagesAttributes(d *schema.ResourceData, memberServerImages []*server.MemberServerImage) error {
	var ids []string

	for _, m := range memberServerImages {
		ids = append(ids, *m.MemberServerImageNo)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("member_server_images", flattenMemberServerImages(memberServerImages)); err != nil {
		return err
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), d.Get("member_server_images"))
	}

	return nil
}
