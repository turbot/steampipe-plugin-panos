package panos

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-panos",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromJSONTag().NullIfZero(),
		DefaultGetConfig: &plugin.GetConfig{
			ShouldIgnoreError: isNotFoundError,
		},
		TableMap: map[string]*plugin.Table{
			"panos_address_object": tablePanosAddressObject(ctx),
		},
	}
	return p
}
