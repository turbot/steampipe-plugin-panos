package main

import (
	"github.com/turbot/steampipe-plugin-panos/panos"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: panos.Plugin})
}
