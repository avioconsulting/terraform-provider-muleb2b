package b2b

import (
	"fmt"
	"github.com/avioconsulting/muleb2b-api-go/muleb2b"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func readSftpConfig(data interface{}) (*muleb2b.EndpointConfig, error) {
	config := data.(*schema.Set).List()
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Failed to parse: %#v", raw)
		}

		address := cfg["server_address"].(string)
		port := cfg["server_port"].(int)
		movedPath := cfg["moved_path"].(string)
		scwt := cfg["size_check_wait_time"].(int)
		pollingFreq := cfg["polling_frequency"].(int)
		path := cfg["path"].(string)

		configName, ok := cfg["config_name"].(string)
		if !ok || configName == "" {
			configName = "sftp"
		}

		amCfg, ok := cfg["auth_mode"]
		if ok {
			authMode, err := readAuthModeConfig(amCfg)

			if err == nil {
				endpointConfig := muleb2b.EndpointConfig{
					MovedPath:         muleb2b.String(movedPath),
					SizeCheckWaitTime: muleb2b.Integer(scwt),
					PollingFrequency:  muleb2b.Integer(pollingFreq),
					Path:              muleb2b.String(path),
					ServerAddress:     muleb2b.String(address),
					ServerPort:        muleb2b.Integer(port),
					ConfigName:        muleb2b.String(configName),
					AuthMode:          authMode,
				}
				return &endpointConfig, nil
			} else {
				return nil, err
			}

		} else {
			return nil, fmt.Errorf("auth_mode is required in sftp_config")
		}
	}
	return nil, fmt.Errorf("sftp_config is required when type is sftp")
}

func readHttpConfig(data interface{}) (*muleb2b.EndpointConfig, error) {
	config := data.(*schema.Set).List()
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Failed to parse: %#v", raw)
		}

		address := cfg["server_address"].(string)
		port := cfg["server_port"].(int)
		path := cfg["path"].(string)
		protocol := cfg["protocol"].(string)
		responseTimeout := cfg["response_timeout"].(int)
		idleTimeout := cfg["connection_idle_timeout"].(int)

		configName, ok := cfg["config_name"].(string)
		if !ok || configName == "" {
			configName = "http"
		}

		amCfg, ok := cfg["auth_mode"]
		var authMode *muleb2b.AuthMode = nil
		var err error = nil
		if ok {
			authMode, err = readAuthModeConfig(amCfg)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("auth_mode is required in http_config")
		}

		tlsCfg, ok := cfg["tls_context"]
		var tlsContext *muleb2b.TlsContext = nil
		if protocol == "https" {
			if ok {
				tlsContext, err = readTlsContextConfig(tlsCfg)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("tls_context is required when protocol is https")
			}
		}

		endpointConfig := muleb2b.EndpointConfig{
			ConfigName:            muleb2b.String(configName),
			ServerAddress:         muleb2b.String(address),
			ServerPort:            muleb2b.Integer(port),
			Path:                  muleb2b.String(path),
			Protocol:              muleb2b.String(strings.ToUpper(protocol)),
			ResponseTimeout:       muleb2b.Integer(responseTimeout),
			ConnectionIdleTimeout: muleb2b.Integer(idleTimeout),
			AuthMode:              authMode,
			TlsContext:            tlsContext,
		}

		return &endpointConfig, nil

	}
	return nil, fmt.Errorf("http_config is required when type is http")
}

func readAuthModeConfig(data interface{}) (*muleb2b.AuthMode, error) {
	config := data.(*schema.Set).List()
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Failed to parse: %#v", raw)
		}

		authType := cfg["type"].(string)

		switch authType {
		case "none":
			authMode := muleb2b.AuthMode{
				AuthType: muleb2b.String("NONE"),
			}
			return &authMode, nil

		case "basic":
			user, uok := cfg["username"].(string)
			password, pok := cfg["password"].(string)
			if uok && pok {
				authMode := muleb2b.AuthMode{
					AuthType: muleb2b.String(strings.ToUpper(authType)),
					Username: muleb2b.String(user),
					Password: muleb2b.String(password),
				}
				return &authMode, nil
			} else {
				return nil, fmt.Errorf("username and password are required when auth_mode.type is basic")
			}

		case "api_key":
			apiKey, kok := cfg["api_key"].(string)
			httpheaderName, hok := cfg["http_header_name"].(string)
			if kok && hok {
				authMode := muleb2b.AuthMode{
					AuthType:       muleb2b.String(strings.ToUpper(authType)),
					ApiKey:         muleb2b.String(apiKey),
					HttpHeaderName: muleb2b.String(httpheaderName),
				}
				return &authMode, nil
			} else {
				return nil, fmt.Errorf("api_key and http_header_name are required when auth_mode.type is api_key")
			}

		case "client_credentials":
			clientId, idok := cfg["client_id"].(string)
			clientSecret, sok := cfg["client_secret"].(string)
			clientIdHeader, idhok := cfg["client_id_header"].(string)
			clientSecretHeader, shok := cfg["client_secret_header"].(string)
			if idhok && idok && sok && shok {
				authMode := muleb2b.AuthMode{
					AuthType:           muleb2b.String(strings.ToUpper(authType)),
					ClientId:           muleb2b.String(clientId),
					ClientSecret:       muleb2b.String(clientSecret),
					ClientIdHeader:     muleb2b.String(clientIdHeader),
					ClientSecretHeader: muleb2b.String(clientSecretHeader),
				}
				return &authMode, nil
			} else {
				return nil, fmt.Errorf("client_id, client_secret, client_id_header, and client_secret_header are required when auth_mode.type is client_credentials")
			}

		case "oauth_token":
			tokenUrl, tok := cfg["token_url"].(string)
			clientId, idok := cfg["client_id"].(string)
			clientSecret, sok := cfg["client_secret"].(string)
			if tok && idok && sok {
				authMode := muleb2b.AuthMode{
					AuthType:     muleb2b.String(strings.ToUpper(authType)),
					TokenUrl:     muleb2b.String(tokenUrl),
					ClientId:     muleb2b.String(clientId),
					ClientSecret: muleb2b.String(clientSecret),
				}
				return &authMode, nil
			} else {
				return nil, fmt.Errorf("token_url, client_id, and client_secret are required when auth_mode.type is oauth_token")
			}

		default:
			return nil, fmt.Errorf("invalid auth_mode.type specified")
		}
	}
	return nil, nil
}

func readTlsContextConfig(data interface{}) (*muleb2b.TlsContext, error) {
	config := data.(*schema.Set).List()
	for _, raw := range config {
		cfg, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Failed to parse: %#v", raw)
		}

		// Both values should be present since it was validated by
		insecure := cfg["insecure"].(bool)
		needCert := cfg["need_certificate"].(bool)

		tlsContext := muleb2b.TlsContext{
			Insecure:        muleb2b.Boolean(insecure),
			NeedCertificate: muleb2b.Boolean(needCert),
		}

		return &tlsContext, nil
	}
	return nil, nil
}

func flattenHttpConfig(endpointConfig *muleb2b.EndpointConfig) []interface{} {
	m := make(map[string]interface{})

	if endpointConfig != nil {
		m["config_name"] = *endpointConfig.ConfigName
		m["server_address"] = *endpointConfig.ServerAddress
		m["server_port"] = *endpointConfig.ServerPort
		m["path"] = *endpointConfig.Path
		m["protocol"] = strings.ToLower(*endpointConfig.Protocol)
		m["response_timeout"] = *endpointConfig.ResponseTimeout
		m["connection_idle_timeout"] = *endpointConfig.ConnectionIdleTimeout
		m["auth_mode"] = flattenAuthMode(endpointConfig.AuthMode)
		if strings.ToLower(*endpointConfig.Protocol) == "https" {
			m["tls_context"] = flattenTlsContext(endpointConfig.TlsContext)
		}
	}

	return []interface{}{m}
}

func flattenSftpConfig(endpointConfig *muleb2b.EndpointConfig) []interface{} {
	m := make(map[string]interface{})

	if endpointConfig != nil {
		m["config_name"] = *endpointConfig.ConfigName
		m["server_address"] = *endpointConfig.ServerAddress
		m["server_port"] = *endpointConfig.ServerPort
		m["path"] = *endpointConfig.Path
		m["moved_path"] = *endpointConfig.MovedPath
		m["size_check_wait_time"] = *endpointConfig.SizeCheckWaitTime
		m["polling_frequency"] = *endpointConfig.PollingFrequency
		m["auth_mode"] = flattenAuthMode(endpointConfig.AuthMode)
	}

	return []interface{}{m}
}

func flattenAuthMode(authMode *muleb2b.AuthMode) []interface{} {
	m := make(map[string]interface{})

	if authMode != nil {
		m["type"] = strings.ToLower(*authMode.AuthType)
		switch *authMode.AuthType {
		case "BASIC":
			m["username"] = *authMode.Username
			m["password"] = *authMode.Password
		case "API_KEY":
			m["api_key"] = *authMode.ApiKey
			m["http_header_name"] = *authMode.HttpHeaderName
		case "CLIENT_CREDENTIALS":
			m["client_id"] = *authMode.ClientId
			m["client_secret"] = *authMode.ClientSecret
			m["client_id_header"] = *authMode.ClientIdHeader
			m["client_secret_header"] = *authMode.ClientSecretHeader
		case "OAUTH_TOKEN":
			m["token_url"] = *authMode.TokenUrl
			m["client_id"] = *authMode.ClientId
			m["client_secret"] = *authMode.ClientSecret
		}
	}

	return []interface{}{m}
}

func flattenTlsContext(context *muleb2b.TlsContext) []interface{} {
	m := make(map[string]interface{})
	if context != nil {
		m["insecure"] = *context.Insecure
		m["need_certificate"] = *context.NeedCertificate
	}
	return []interface{}{m}
}

func expandHttpConfig(d interface{}) *muleb2b.EndpointConfig {
	if d != nil {
		configList := d.(*schema.Set).List()
		if len(configList) > 0 {
			configData := configList[0].(map[string]interface{})
			endpointConfig := muleb2b.EndpointConfig{
				ConfigName:            muleb2b.String(configData["config_name"].(string)),
				ServerAddress:         muleb2b.String(configData["server_address"].(string)),
				ServerPort:            muleb2b.Integer(configData["server_port"].(int)),
				Path:                  muleb2b.String(configData["path"].(string)),
				Protocol:              muleb2b.String(strings.ToUpper(configData["protocol"].(string))),
				ResponseTimeout:       muleb2b.Integer(configData["response_timeout"].(int)),
				ConnectionIdleTimeout: muleb2b.Integer(configData["connection_idle_timeout"].(int)),
			}

			endpointConfig.AuthMode = expandAuthMode(configData["auth_mode"])

			if strings.ToLower(*endpointConfig.Protocol) == "https" {
				endpointConfig.TlsContext = expandTlsContext(configData["tls_context"])
			}

			return &endpointConfig
		}
	}
	return nil
}

func expandSftpConfig(d interface{}) *muleb2b.EndpointConfig {
	if d != nil {
		configList := d.(*schema.Set).List()
		if len(configList) > 0 {
			configData := configList[0].(map[string]interface{})
			endpointConfig := muleb2b.EndpointConfig{
				ConfigName:        muleb2b.String(configData["config_name"].(string)),
				ServerAddress:     muleb2b.String(configData["server_address"].(string)),
				ServerPort:        muleb2b.Integer(configData["server_port"].(int)),
				Path:              muleb2b.String(configData["path"].(string)),
				MovedPath:         muleb2b.String(configData["moved_path"].(string)),
				SizeCheckWaitTime: muleb2b.Integer(configData["size_check_wait_time"].(int)),
				PollingFrequency:  muleb2b.Integer(configData["polling_frequency"].(int)),
			}

			endpointConfig.AuthMode = expandAuthMode(configData["auth_mode"])

			return &endpointConfig
		}
	}
	return nil
}

func expandAuthMode(d interface{}) *muleb2b.AuthMode {
	if d != nil {
		authList := d.(*schema.Set).List()
		if len(authList) > 0 {
			authData := authList[0].(map[string]interface{})
			authType := strings.ToUpper(authData["type"].(string))
			switch authType {
			case "NONE":
				authMode := muleb2b.AuthMode{
					AuthType: muleb2b.String(authType),
				}
				return &authMode

			case "BASIC":
				authMode := muleb2b.AuthMode{
					AuthType: muleb2b.String(authType),
					Username: muleb2b.String(authData["username"].(string)),
					Password: muleb2b.String(authData["password"].(string)),
				}
				return &authMode

			case "API_KEY":
				authMode := muleb2b.AuthMode{
					AuthType:       muleb2b.String(authType),
					ApiKey:         muleb2b.String(authData["api_key"].(string)),
					HttpHeaderName: muleb2b.String(authData["http_header_name"].(string)),
				}
				return &authMode

			case "CLIENT_CREDENTIALS":
				authMode := muleb2b.AuthMode{
					AuthType:           muleb2b.String(authType),
					ClientId:           muleb2b.String(authData["client_id"].(string)),
					ClientSecret:       muleb2b.String(authData["client_secret"].(string)),
					ClientIdHeader:     muleb2b.String(authData["client_id_header"].(string)),
					ClientSecretHeader: muleb2b.String(authData["client_secret_header"].(string)),
				}
				return &authMode

			case "OAUTH_TOKEN":
				authMode := muleb2b.AuthMode{
					AuthType:     muleb2b.String(authType),
					ClientId:     muleb2b.String(authData["client_id"].(string)),
					ClientSecret: muleb2b.String(authData["client_secret"].(string)),
					TokenUrl:     muleb2b.String(authData["token_url"].(string)),
				}
				return &authMode
			}
		}
	}
	return nil
}

func expandTlsContext(d interface{}) *muleb2b.TlsContext {
	if d != nil {
		tlsList := d.(*schema.Set).List()
		if len(tlsList) > 0 {
			tlsData := tlsList[0].(map[string]interface{})
			tlsContext := muleb2b.TlsContext{
				Insecure:        muleb2b.Boolean(tlsData["insecure"].(bool)),
				NeedCertificate: muleb2b.Boolean(tlsData["need_certificate"].(bool)),
			}
			return &tlsContext
		}
	}
	return nil
}
