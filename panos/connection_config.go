package panos

import (
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/schema"
)

type panosConfig struct {
	Hostname *string `cty:"hostname"`
	APIKey   *string `cty:"api_key"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"hostname": {
		Type: schema.TypeString,
	},
	"api_key": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &panosConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) panosConfig {
	if connection == nil || connection.Config == nil {
		return panosConfig{}
	}
	config, _ := connection.Config.(panosConfig)
	return config
}
