package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func resourceEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointCreate,
		Read:   resourceEndpointRead,
		Update: resourceEndpointUpdate,
		Delete: resourceEndpointDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name to give the endpoint",
			},
			"role": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRole,
				Description:  "The role the endpoint will play: send, receive, receive_ack, storage_api",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateType,
				Description:  "The type of endpoint, http or sftp",
			},
			"partner_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the partner that owns the endpoint",
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the environment that the endpoint is placed into",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Description of the endpoint",
			},
			"partner_certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the certificate to use when a certificate is needed",
			},
			"http_config": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "HTTP configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "http",
							Description: "name of the endpoint configuration",
						},
						"server_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Address of the HTTP service",
						},
						"server_port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Port of the HTTP service",
						},
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path of the HTTP service",
						},
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateEndpointHttpProtocol,
							Description:  "Protocol of the service. http or https",
						},
						"response_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     15000,
							Description: "Time to wait for a service response (ms)",
						},
						"connection_idle_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     30000,
							Description: "Time to wait before the connection is considered idle (ms)",
						},
						"auth_mode": endpointAuthModeSchema(),
						"tls_context": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"insecure": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not the connection is insecure",
									},
									"need_certificate": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "Whether or not a certificate is needed",
									},
								},
							},
						},
					},
				},
			},
			"sftp_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "sftp",
							Description: "name of the endpoint configuration",
						},
						"server_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Address of the sftp server",
						},
						"server_port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Port of the sftp server",
						},
						"archive_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path to move read files into",
						},
						"size_check_wait_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1000,
							Description: "Time to wait to check the file size",
						},
						"polling_frequency": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1000,
							Description: "Time to wait between checking for new files",
						},
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to look for new files",
						},
						"auth_mode": endpointAuthModeSchema(),
					},
				},
			},
		},
	}
}

func endpointAuthModeSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validateAuthModeType,
					Description:  "Authentication Mode selected: none, basic, api_key, client_credentials, or oauth_token",
				},
				"username": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The username to use for authentication",
				},
				"password": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "The password to use for authentication",
				},
				"http_header_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The name of the HTTP header to use",
				},
				"api_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "API key to use for authentication",
				},
				"client_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "ID of the client to use for authentication",
				},
				"client_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "Secret to use in authentication",
				},
				"client_id_header": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Header to use for the client id",
				},
				"client_secret_header": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Header to use for the client secret",
				},
				"token_url": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "URL for the OAUTH token",
				},
			},
		},
	}
}

func validateAuthModeType(value interface{}, key string) (warnings []string, errors []error) {
	v, ok := value.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", key))
	}

	if v != "basic" && v != "none" && v != "api_key" && v != "client_credentials" && v != "oauth_token" {
		errors = append(errors, fmt.Errorf("value of %q must be none, basic, api_key, client_credentials, or oauth_token", key))
	}
	return warnings, errors
}

func validateType(value interface{}, key string) (warnings []string, errors []error) {
	v, ok := value.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", key))
	}

	if v != "sftp" && v != "http" {
		errors = append(errors, fmt.Errorf("value of %q must be sftp or http", key))
	}
	return warnings, errors
}

func validateRole(value interface{}, key string) (warnings []string, errors []error) {
	v, ok := value.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", key))
	}

	if v != "send" && v != "storage_api" && v != "receive_ack" && v != "receive" {
		errors = append(errors, fmt.Errorf("value of %q must be send, receive, receive_ack, storage_api", key))
	}
	return warnings, errors
}

func validateEndpointHttpProtocol(value interface{}, key string) (warnings []string, errors []error) {
	v, ok := value.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", key))
	}

	if v != "http" && v != "https" {
		errors = append(errors, fmt.Errorf("value of %q must be http or https", key))
	}
	return warnings, errors
}

func resourceEndpointCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*muleb2b.Client)

	name := d.Get("name").(string)
	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	role := d.Get("role").(string)
	endType := d.Get("type").(string)

	partnerId := d.Get("partner_id").(string)

	endpoint := muleb2b.Endpoint{
		Name:          muleb2b.String(name),
		EndpointRole:  muleb2b.String(strings.ToUpper(role)),
		EnvironmentID: muleb2b.String(envId),
		EndpointType:  muleb2b.String(endType),
		PartnerID:     muleb2b.String(partnerId),
	}

	desc, ok := d.Get("description").(string)
	if ok {
		endpoint.Description = muleb2b.String(desc)
	} else {
		endpoint.Description = muleb2b.String("")
	}

	if v, ok := d.GetOk("partner_certificate_id"); ok {
		endpoint.PartnerCertificateID = muleb2b.String(v.(string))
	}

	// Preserve initial settings
	if endType == "sftp" {
		cfg, ok := d.GetOk("sftp_config")
		if ok {
			endpointCfg, err := readSftpConfig(cfg)
			if err != nil {
				return err
			}
			endpoint.Config = endpointCfg
		} else {
			return fmt.Errorf("sftp_config is required when type is set to sftp")
		}

	} else if endType == "http" {
		cfg, ok := d.GetOk("http_config")
		if ok {
			endpointCfg, err := readHttpConfig(cfg)
			if err != nil {
				return err
			}
			endpoint.Config = endpointCfg
		} else {
			return fmt.Errorf("http_config is required when type is set to http")
		}
	}

	if *endpoint.EndpointType == "sftp" {
		if err := d.Set("sftp_config", flattenSftpConfig(endpoint.Config, nil)); err != nil {
			return err
		}
	} else if *endpoint.EndpointType == "http" {
		if err := d.Set("http_config", flattenHttpConfig(endpoint.Config, nil)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported endpoint type: %s", *endpoint.EndpointType)
	}

	id, err := client.CreateEndpoint(endpoint)

	if err != nil {
		return err
	}

	d.SetId(*id)

	return resourceEndpointRead(d, m)
}

func resourceEndpointRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)
	id := d.Id()

	envId := d.Get("environment_id").(string) // Should be set on the resource and in the state
	client.SetEnvironment(envId)

	endpoint, err := client.GetEndpoint(id)

	if err != nil {
		return err
	}

	d.Set("name", *endpoint.Name)
	d.Set("role", strings.ToLower(*endpoint.EndpointRole))
	d.Set("type", *endpoint.EndpointType)
	d.Set("partner_id", *endpoint.PartnerID)
	d.Set("environment_id", *endpoint.EnvironmentID)
	if endpoint.PartnerCertificateID != nil {
		d.Set("partner_certificate_id", *endpoint.PartnerCertificateID)
	}
	if endpoint.Description != nil {
		d.Set("description", *endpoint.Description)
	}

	// Retrieve sensitive data from state - this can be improved
	var sensitive *sensitiveData = nil
	if *endpoint.EndpointType == "sftp" {
		endpoint.Config = expandSftpConfig(d.Get("sftp_config"))
		if endpoint.Config.AuthMode.Password != nil {
			sensitive = &sensitiveData{
				password: endpoint.Config.AuthMode.Password,
			}
		} else if endpoint.Config.AuthMode.ClientSecret != nil {
			sensitive = &sensitiveData{
				clientSecret: endpoint.Config.AuthMode.ClientSecret,
			}
		} else if endpoint.Config.AuthMode.ApiKey != nil {
			sensitive = &sensitiveData{
				apiKey: endpoint.Config.AuthMode.ApiKey,
			}
		}
	} else if *endpoint.EndpointType == "http" {
		endpoint.Config = expandHttpConfig(d.Get("http_config"))
		if endpoint.Config.AuthMode.Password != nil {
			sensitive = &sensitiveData{
				password: endpoint.Config.AuthMode.Password,
			}
		} else if endpoint.Config.AuthMode.ClientSecret != nil {
			sensitive = &sensitiveData{
				clientSecret: endpoint.Config.AuthMode.ClientSecret,
			}
		} else if endpoint.Config.AuthMode.ApiKey != nil {
			sensitive = &sensitiveData{
				apiKey: endpoint.Config.AuthMode.ApiKey,
			}
		}
	}

	if *endpoint.EndpointType == "sftp" {
		if err = d.Set("sftp_config", flattenSftpConfig(endpoint.Config, sensitive)); err != nil {
			return err
		}
	} else if *endpoint.EndpointType == "http" {
		if err = d.Set("http_config", flattenHttpConfig(endpoint.Config, sensitive)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported endpoint type: %s", *endpoint.EndpointType)
	}

	return nil
}

func resourceEndpointUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*muleb2b.Client)

	envId := d.Get("environment_id").(string) // Should be set on the resource
	client.SetEnvironment(envId)

	endpoint := muleb2b.Endpoint{
		ID:            muleb2b.String(d.Id()),
		Name:          muleb2b.String(d.Get("name").(string)),
		EndpointRole:  muleb2b.String(strings.ToUpper(d.Get("role").(string))),
		EndpointType:  muleb2b.String(d.Get("type").(string)),
		EnvironmentID: muleb2b.String(d.Get("environment_id").(string)),
		PartnerID:     muleb2b.String(d.Get("partner_id").(string)),
	}

	desc, ok := d.Get("description").(string)
	if ok {
		endpoint.Description = muleb2b.String(desc)
	}

	if v, ok := d.GetOk("partner_certificate_id"); ok {
		endpoint.PartnerCertificateID = muleb2b.String(v.(string))
	}

	if *endpoint.EndpointType == "sftp" {
		endpoint.Config = expandSftpConfig(d.Get("sftp_config"))
	} else if *endpoint.EndpointType == "http" {
		endpoint.Config = expandHttpConfig(d.Get("http_config"))
	}

	err := client.UpdateEndpoint(endpoint)
	if err != nil {
		return err
	}

	return resourceEndpointRead(d, m)
}

func resourceEndpointDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*muleb2b.Client)
	id := d.Id()

	envId := d.Get("environment_id").(string) // Should be set on the resource and in the state
	client.SetEnvironment(envId)

	err := client.DeleteEndpoint(id)

	return err
}
