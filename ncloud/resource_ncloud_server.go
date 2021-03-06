package ncloud

import (
	"fmt"
	"log"
	"time"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNcloudServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceNcloudServerCreate,
		Read:   resourceNcloudServerRead,
		Delete: resourceNcloudServerDelete,
		Update: resourceNcloudServerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(DefaultCreateTimeout),
			Delete: schema.DefaultTimeout(DefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"server_image_product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server image product code to determine which server image to create. It can be obtained through getServerImageProductList. You are required to select one among two parameters: server image product code (server_image_product_code) and member server image number(member_server_image_no).",
			},
			"server_product_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server product code to determine the server specification to create. It can be obtained through the getServerProductList action. Default : Selected as minimum specification. The minimum standards are 1. memory 2. CPU 3. basic block storage size 4. disk type (NET,LOCAL)",
			},
			"member_server_image_no": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Required value when creating a server from a manually created server image. It can be obtained through the getMemberServerImageList action.",
			},
			"server_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateServerName,
				Description:  "Server name to create. default: Assigned by ncloud",
			},
			"server_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server description to create",
			},
			"login_key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The login key name to encrypt with the public key. Default : Uses the most recently created login key name",
			},
			"is_protect_server_termination": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "You can set whether or not to protect return when creating. default : false",
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInternetLineTypeCode,
				Description:  "Internet line identification code. PUBLC(Public), GLBL(Global). default : PUBLC(Public)",
			},
			"fee_system_type_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A rate system identification code. There are time plan(MTRAT) and flat rate (FXSUM). Default : Time plan(MTRAT)",
			},
			"zone_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone code. You can determine the ZONE where the server will be created. It can be obtained through the getZoneList action. Default : Assigned by NAVER Cloud Platform.",
				ConflictsWith: []string{"zone_no"},
			},
			"zone_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone number. You can determine the ZONE where the server will be created. It can be obtained through the getZoneList action. Default : Assigned by NAVER Cloud Platform.",
				ConflictsWith: []string{"zone_code"},
			},

			"access_control_group_configuration_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
				Description: "You can set the ACG created when creating the server. ACG setting number can be obtained through the getAccessControlGroupList action. Default : Default ACG number",
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The server will execute the user data script set by the user at first boot. To view the column, it is returned only when viewing the server instance. You must need base64 Encoding, URL Encoding before put in value of userData. If you don't URL Encoding again it occurs signature invalid error.",
			},
			"raid_type_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Raid Type Name",
			},
			"tag_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        tagListSchemaResource,
				Description: "Instance tag list",
			},

			"server_instance_no": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"base_block_storage_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"platform_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"is_fee_charging_monitoring": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_instance_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"server_instance_operation": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"server_instance_status_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uptime": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_forwarding_public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_forwarding_external_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_forwarding_internal_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     zoneSchemaResource,
			},
			"region": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     regionSchemaResource,
			},
			"base_block_storage_disk_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"base_block_storage_disk_detail_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
			"internet_line_type": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     commonCodeSchemaResource,
			},
		},
	}
}

func resourceNcloudServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	reqParams, err := buildCreateServerInstanceReqParams(client, d)
	if err != nil {
		return err
	}

	var resp *server.CreateServerInstancesResponse
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		var err error
		logCommonRequest("CreateServerInstances", reqParams)
		resp, err = client.server.V2Api.CreateServerInstances(reqParams)

		log.Printf("[DEBUG] resourceNcloudServerCreate resp: %v", resp)
		if resp != nil && isRetryableErr(GetCommonResponse(resp), []string{ApiErrorUnknown, ApiErrorAuthorityParameter, ApiErrorServerObjectInOperation, ApiErrorPreviousServersHaveNotBeenEntirelyTerminated}) {
			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		logErrorResponse("CreateServerInstances", err, reqParams)
		return err
	}
	logCommonResponse("CreateServerInstances", GetCommonResponse(resp))

	serverInstance := resp.ServerInstanceList[0]
	d.SetId(ncloud.StringValue(serverInstance.ServerInstanceNo))

	if err := waitForServerInstance(client, ncloud.StringValue(serverInstance.ServerInstanceNo), "RUN"); err != nil {
		return err
	}
	return resourceNcloudServerRead(d, meta)
}

func resourceNcloudServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	instance, err := getServerInstance(client, d.Id())
	if err != nil {
		return err
	}

	if instance != nil {
		d.Set("server_instance_no", instance.ServerInstanceNo)
		d.Set("server_name", instance.ServerName)
		d.Set("server_image_product_code", instance.ServerImageProductCode)
		d.Set("server_instance_status_name", instance.ServerInstanceStatusName)
		d.Set("uptime", instance.Uptime)
		d.Set("server_image_name", instance.ServerImageName)
		d.Set("private_ip", instance.PrivateIp)
		d.Set("cpu_count", instance.CpuCount)
		d.Set("memory_size", instance.MemorySize)
		d.Set("base_block_storage_size", instance.BaseBlockStorageSize)
		d.Set("is_fee_charging_monitoring", instance.IsFeeChargingMonitoring)
		d.Set("public_ip", instance.PublicIp)
		d.Set("private_ip", instance.PrivateIp)
		d.Set("create_date", instance.CreateDate)
		d.Set("uptime", instance.Uptime)
		d.Set("port_forwarding_public_ip", instance.PortForwardingPublicIp)
		d.Set("port_forwarding_external_port", instance.PortForwardingExternalPort)
		d.Set("port_forwarding_internal_port", instance.PortForwardingInternalPort)
		d.Set("user_data", d.Get("user_data").(string))

		if err := d.Set("server_instance_status", flattenCommonCode(instance.ServerInstanceStatus)); err != nil {
			return err
		}
		if err := d.Set("platform_type", flattenCommonCode(instance.PlatformType)); err != nil {
			return err
		}
		if err := d.Set("server_instance_operation", flattenCommonCode(instance.ServerInstanceOperation)); err != nil {
			return err
		}
		if err := d.Set("zone", flattenZone(instance.Zone)); err != nil {
			return err
		}
		if err := d.Set("region", flattenRegion(instance.Region)); err != nil {
			return err
		}
		if err := d.Set("base_block_storage_disk_type", flattenCommonCode(instance.BaseBlockStorageDiskType)); err != nil {
			return err
		}
		if err := d.Set("base_block_storage_disk_detail_type", flattenCommonCode(instance.BaseBlockStroageDiskDetailType)); err != nil {
			return err
		}
		if err := d.Set("internet_line_type", flattenCommonCode(instance.InternetLineType)); err != nil {
			return err
		}
		if len(instance.InstanceTagList) != 0 {
			d.Set("load_balancer_rule_list", flattenInstanceTagList(instance.InstanceTagList))
		}
	}

	return nil
}

func resourceNcloudServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)
	serverInstance, err := getServerInstance(client, d.Id())
	if err != nil {
		return err
	}

	if serverInstance == nil || ncloud.StringValue(serverInstance.ServerInstanceStatus.Code) != "NSTOP" {
		if err := stopServerInstance(client, d.Id()); err != nil {
			return err
		}
		if err := waitForServerInstance(client, ncloud.StringValue(serverInstance.ServerInstanceNo), "NSTOP"); err != nil {
			return err
		}
	}

	err = detachBlockStorageByServerInstanceNo(d, client, d.Id())
	if err != nil {
		log.Printf("[ERROR] detachBlockStorageByServerInstanceNo err: %s", err)
		return err
	}

	if err := terminateServerInstance(client, d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceNcloudServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	if d.HasChange("server_product_code") {
		reqParams := &server.ChangeServerInstanceSpecRequest{
			ServerInstanceNo:  ncloud.String(d.Get("server_instance_no").(string)),
			ServerProductCode: ncloud.String(d.Get("server_product_code").(string)),
		}

		var resp *server.ChangeServerInstanceSpecResponse
		err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			var err error
			logCommonRequest("ChangeServerInstanceSpec", reqParams)
			resp, err = client.server.V2Api.ChangeServerInstanceSpec(reqParams)

			if resp != nil && isRetryableErr(GetCommonResponse(resp), []string{ApiErrorUnknown, ApiErrorObjectInOperation, ApiErrorObjectInOperation}) {
				logErrorResponse("retry ChangeServerInstanceSpec", err, reqParams)
				time.Sleep(time.Second * 5)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err != nil {
			logErrorResponse("ChangeServerInstanceSpec", err, reqParams)
			return err
		}
		logCommonResponse("ChangeServerInstanceSpec", GetCommonResponse(resp))
	}

	return resourceNcloudServerRead(d, meta)
}

func buildCreateServerInstanceReqParams(client *NcloudAPIClient, d *schema.ResourceData) (*server.CreateServerInstancesRequest, error) {

	var paramAccessControlGroupConfigurationNoList []*string
	if param, ok := d.GetOk("access_control_group_configuration_no_list"); ok {
		paramAccessControlGroupConfigurationNoList = expandStringInterfaceList(param.([]interface{}))
	}
	zoneNo, err := parseZoneNoParameter(client, d)
	if err != nil {
		return nil, err
	}
	reqParams := &server.CreateServerInstancesRequest{
		ServerImageProductCode:                ncloud.String(d.Get("server_image_product_code").(string)),
		ServerProductCode:                     ncloud.String(d.Get("server_product_code").(string)),
		MemberServerImageNo:                   ncloud.String(d.Get("member_server_image_no").(string)),
		ServerName:                            ncloud.String(d.Get("server_name").(string)),
		ServerDescription:                     ncloud.String(d.Get("server_description").(string)),
		LoginKeyName:                          ncloud.String(d.Get("login_key_name").(string)),
		InternetLineTypeCode:                  StringPtrOrNil(d.GetOk("internet_line_type_code")),
		FeeSystemTypeCode:                     ncloud.String(d.Get("fee_system_type_code").(string)),
		ZoneNo:                                zoneNo,
		AccessControlGroupConfigurationNoList: paramAccessControlGroupConfigurationNoList,
		UserData:                              ncloud.String(d.Get("user_data").(string)),
		RaidTypeName:                          ncloud.String(d.Get("raid_type_name").(string)),
	}

	if instanceTagList, err := expandTagListParams(d.Get("tag_list").([]interface{})); err == nil {
		reqParams.InstanceTagList = instanceTagList
	}

	if IsProtectServerTermination, ok := d.GetOk("is_protect_server_termination"); ok {
		reqParams.IsProtectServerTermination = ncloud.Bool(IsProtectServerTermination.(bool))
	}

	return reqParams, nil
}

func getServerInstance(client *NcloudAPIClient, serverInstanceNo string) (*server.ServerInstance, error) {
	reqParams := new(server.GetServerInstanceListRequest)
	reqParams.ServerInstanceNoList = []*string{ncloud.String(serverInstanceNo)}
	logCommonRequest("GetServerInstanceList", reqParams)

	resp, err := client.server.V2Api.GetServerInstanceList(reqParams)

	if err != nil {
		logErrorResponse("GetServerInstanceList", err, reqParams)
		return nil, err
	}
	logCommonResponse("GetServerInstanceList", GetCommonResponse(resp))
	if len(resp.ServerInstanceList) > 0 {
		inst := resp.ServerInstanceList[0]
		return inst, nil
	}
	return nil, nil
}

func getServerZoneNo(client *NcloudAPIClient, serverInstanceNo string) (string, error) {
	serverInstance, err := getServerInstance(client, serverInstanceNo)
	if err != nil || serverInstance == nil || serverInstance.Zone == nil {
		return "", err
	}
	return *serverInstance.Zone.ZoneNo, nil
}

func stopServerInstance(client *NcloudAPIClient, serverInstanceNo string) error {
	reqParams := &server.StopServerInstancesRequest{
		ServerInstanceNoList: []*string{ncloud.String(serverInstanceNo)},
	}
	logCommonRequest("StopServerInstances", reqParams)
	resp, err := client.server.V2Api.StopServerInstances(reqParams)
	if err != nil {
		logErrorResponse("StopServerInstances", err, reqParams)
		return err
	}
	logCommonResponse("StopServerInstances", GetCommonResponse(resp))

	return nil
}

func terminateServerInstance(client *NcloudAPIClient, serverInstanceNo string) error {
	reqParams := &server.TerminateServerInstancesRequest{
		ServerInstanceNoList: []*string{ncloud.String(serverInstanceNo)},
	}

	var resp *server.TerminateServerInstancesResponse
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		var err error
		logCommonRequest("TerminateServerInstances", reqParams)
		resp, err = client.server.V2Api.TerminateServerInstances(reqParams)
		if err == nil && resp == nil {
			return resource.NonRetryableError(err)
		}
		if resp != nil && isRetryableErr(GetCommonResponse(resp), []string{ApiErrorUnknown, ApiErrorServerObjectInOperation2}) {
			logErrorResponse("retry TerminateServerInstances", err, reqParams)
			return resource.RetryableError(err)
		}
		logCommonResponse("TerminateServerInstances", GetCommonResponse(resp))
		return resource.NonRetryableError(err)
	})

	if err != nil {
		logErrorResponse("TerminateServerInstances", err, reqParams)
		return err
	}
	return nil
}

func waitForServerInstance(client *NcloudAPIClient, instanceId string, status string) error {

	c1 := make(chan error, 1)

	go func() {
		for {
			instance, err := getServerInstance(client, instanceId)

			if err != nil {
				c1 <- err
				return
			}
			if instance == nil || ncloud.StringValue(instance.ServerInstanceStatus.Code) == status {
				c1 <- nil
				return
			}
			log.Printf("[DEBUG] Wait server instance [%s] status [%s] to be [%s]", instanceId, ncloud.StringValue(instance.ServerInstanceStatus.Code), status)
			time.Sleep(time.Second * 1)
		}
	}()

	select {
	case res := <-c1:
		return res
	case <-time.After(DefaultCreateTimeout):
		return fmt.Errorf("TIMEOUT : Wait to server instance  (%s)", instanceId)
	}
}

var tagListSchemaResource = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"tag_key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Instance Tag Key",
		},
		"tag_value": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Instance Tag Value",
		},
	},
}
