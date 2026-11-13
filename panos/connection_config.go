package panos

import (
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type panosConfig struct {
	Hostname *string `hcl:"hostname"`
	APIKey   *string `hcl:"api_key"`
	Username *string `hcl:"username"`
	Password *string `hcl:"password"`
	Timeout  *int    `hcl:"timeout"`
}

func ConfigInstance() interface{} {
	return &panosConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) panosConfig {
	if connection == nil || connection.GetConfig() == nil {
		return panosConfig{}
	}
	config, _ := connection.GetConfig().(panosConfig)
	return config
}
