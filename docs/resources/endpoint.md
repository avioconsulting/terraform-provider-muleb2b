# Endpoint Resource

Provides [Mule B2B Endpoint][1] resource. 

## Example Usage

```hcl
data "muleb2b_environment" "sbx" {
  name = "Sandbox"
}

data "muleb2b_identifier_type" "as2" {
  environment_id = data.muleb2b_environment.sbx.id
  name = "AS2"
}

resource "muleb2b_partner" "test" {
  name           = "Test"
  environment_id = data.muleb2b_environment.sbx.id
  identifier {
    identifier_type_id = data.muleb2b_identifier_type.as2.id
    value = "ab43c"
  }
}

resource "muleb2b_endpoint" "test" {
  name = "test-receive"
  role = "receive_ack"
  type = "http"
  partner_id = muleb2b_partner.test.id
  environment_id = data.muleb2b_environment.sbx.id
  http_config {
    server_address = "test.mytest.com"
    server_port = 80
    path = "/"
    protocol = "http"
    auth_mode  {
      type = "none"
    }
  }
}
```

## Argument Reference

* `name` - (Required) Name for the endpoint
* `role` - (Required) The role the endpoint will play. Can be `"send"`, `"receive"`, `"receive_ack"`, `"storage_api"`.
* `type` - (Required) The type of endpoint. Can be `"http"` or `"sftp"`
* `partner_id` - (Required) The id of the partner for which the endpoint will be created
* `environment_id` - (Required) The id of the environment in which the endpoint will be created
* `description` - (Optional) Description of the endpoint's use
* `partner_certificate_id` - (Optional) The id of the certificate to use when one is needed
* `http_config` - (Optional) Either this argument, or `sftp_config` must be present
* `sftp_config` - (Optional) Either this argument, or `http_config` must be present

#### HTTP Config
The `http_config` block allows one to configure the endpoint's HTTP settings

* `config_name` - (Optional) name of the endpoint configuration. Defaults to `"http"`
* `server_address` - (Required) Address of the HTTP server
* `server_port` - (Required) Port of the HTTP server
* `path` - (Required) Path on the HTTP server
* `protocol` - (Required) Protocol for the HTTP server. Can be `"http"` or `"https"`
* `response_timeout` - (Optional) Timeout in milliseconds. Must be greater than `1000`. Default is `1000`.
* `connection_idle_timeout` - (Optional) Timeout in milliseconds. Must be greater than `1000`. Default is `3000`.
* `tls_context` - (Optional) TLS settings
* `auth_mode` - (Required) Auth mode for the HTTP service

#### SFTP Config
The `sftp_config` block allows one to configure the endpoint's SFTP settings

* `config_name` - (Optional) name of the endpoint configuration. Defaults to `"sftp"`
* `server_address` - (Required) Address of the SFTP server
* `server_port` - (Required) Port of the SFTP server
* `path` - (Required) Path for files on the SFTP server
* `archive_path` - (Optional) Path files will be archived to after being processed
* `size_check_wait_time` - (Optional) The wait time in milliseconds between size checks to determine if a file is ready to be processed. Defaults to `1000`
* `polling_frequency` - (Optional) Frequency in milliseconds to check the source path for new files. Defaults to `1000`
* `auth_mode` - (Required) Auth mode for the SFTP service

##### TLS Context
The `tls_context` block, as part of the `http_config` block, allows one to configure the TLS settings 
* `insecure` - (Optional)  `true` if the connection can be insecure. Defaults to `false`
* `need_certificate` - (Optional) `true` if providing custom certificate. Defaults to `false` 

##### Auth Mode
The `auth_mode` block, as part of the `http_config` and `sftp_config` blocks, allows one to configure the authentication on an endpoint
* `type` - (Required) Authentication Type. Can be `"none"`, `"basic"`, `"api_key"`, `"client_credentials"`, `"oauth_token"`
* `username` - (Optional) Required when `type` is `"basic"`
* `password` - (Optional) Required when `type` is `"basic"`
* `http_header_name` - (Optional) Header parameter associated to the API Key. Required when `type` is `"api_key"`
* `api_key` - (Optional) The value for the access key. Required when `type` is `"api_key"`
* `client_id` - (Optional) Client ID provided when registering your application. Required when `type` is `"client_credentials"` or `"oauth_token"` 
* `client_secret` - (Optional) The client secret provided when registering your application. Required when `type` is `"client_credentials"` or `"oauth_token"`
* `client_id_header` - (Optional) The header used for client id. Required when `type` is `"client_credentials"` 
* `client_secret_header` - (Optional) The header used for client secret. Required when `type` is `"client_credentials"`
* `token_url` - (Optional) The authorization URL used when `type` is `"oauth_token"`

## Attribute Reference

* `id` - ID of the endpoint

[1]: https://docs.mulesoft.com/partner-manager/2.0/endpoints