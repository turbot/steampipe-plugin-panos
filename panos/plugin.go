package panos

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-panos",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		DefaultGetConfig: &plugin.GetConfig{
			ShouldIgnoreError: isNotFoundError,
		},
		TableMap: map[string]*plugin.Table{
			"panos_address_object":     tablePanosAddressObject(ctx),
			"panos_administrative_tag": tablePanosAdministrativeTag(ctx),
			"panos_nat_rule":           tablePanosNATRule(ctx),
			"panos_security_rule":      tablePanosSecurityRule(ctx),
		},
	}
	return p
}
