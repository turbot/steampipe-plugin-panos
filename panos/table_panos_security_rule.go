package panos

import (
	"context"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tablePanosSecurityRule(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "panos_security_rule",
		Description: "Security rules for the PAN-OS device.",
		List: &plugin.ListConfig{
			Hydrate: listSecurityRule,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vsys", Require: plugin.Optional},
				{Name: "device_group", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
				{Name: "rule_base", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The security rule's name."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of security rule."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The security rule's description."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "List of administrative tags."},
			{Name: "source_zones", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "source_addresses", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "negate_source", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "source_users", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "hip_profiles", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "destination_zones", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "destination_addresses", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "negate_destination", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "applications", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "services", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "categories", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "action", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "log_setting", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "log_start", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "log_end", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "disabled", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "schedule", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "icmp_unreachable", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "disable_server_response_inspection", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "group", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "targets", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "negate_target", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "virus", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "spyware", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "vulnerability", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "url_filtering", Type: proto.ColumnType_STRING, Transform: transform.FromField("UrlFiltering"), Description: ""},
			{Name: "file_blocking", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "wild_fire_analysis", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "data_filtering", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "vsys", Type: proto.ColumnType_STRING, Transform: transform.FromField("VSys").NullIfZero(), Description: "[NGFW] The vsys the security rule belongs to (default: vsys1)."},
			{Name: "device_group", Type: proto.ColumnType_STRING, Transform: transform.FromField("DeviceGroup").NullIfZero(), Description: "[Panorama] The device group location (default: shared)"},
			{Name: "rule_base", Type: proto.ColumnType_STRING, Transform: transform.FromField("RuleBase").NullIfZero(), Description: "The rulebase. This can be either pre-rulebase (default), rulebase, or post-rulebase."},
			{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Description: "Raw view of data for the security rule."},
		},
	}
}

type securityRuleStruct struct {
	VSys        string
	DeviceGroup string
	RuleBase    string
	security.Entry
}

func listSecurityRule(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	plugin.Logger(ctx).Debug("panos_security_rule.listSecurityRule", "step", "about to connect")

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("panos_security_rule.listSecurityRule", "connection_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("panos_security_rule.listSecurityRule", "conn", conn)

	// URL parameters for all queries
	keyQuals := d.KeyColumnQuals
	var vsys, deviceGroup, name string
	var listing []security.Entry
	var entry security.Entry

	ruleBase := util.PreRulebase
	if keyQuals["rule_base"] != nil {
		ruleBase = keyQuals["rule_base"].GetStringValue()
	}

	switch client := conn.(type) {
	case *pango.Firewall:
		{
			vsys = "vsys1"
			if keyQuals["vsys"] != nil {
				vsys = keyQuals["vsys"].GetStringValue()
			}
			plugin.Logger(ctx).Debug("panos_security_rule.listSecurityRule", "Firewall.id", vsys)

			if name != "" {
				plugin.Logger(ctx).Debug("panos_security_rule.listSecurityRule", "Firewall.name", name)
				entry, err = client.Policies.Security.Get(vsys, name)
				listing = []security.Entry{entry}
			} else {
				listing, err = client.Policies.Security.GetAll(vsys)
			}
		}
	case *pango.Panorama:
		{
			deviceGroup = "shared"
			if keyQuals["device_group"] != nil {
				deviceGroup = keyQuals["device_group"].GetStringValue()
			}
			plugin.Logger(ctx).Debug("panos_security_rule.listSecurityRule", "Panorama.id", deviceGroup)

			if name != "" {
				plugin.Logger(ctx).Debug("panos_security_rule.listSecurityRule", "Panorama.name", name)

				entry, err = client.Policies.Security.Get(deviceGroup, ruleBase, name)
				listing = append(listing, entry)
			} else {
				listing, err = client.Policies.Security.GetAll(deviceGroup, ruleBase)
			}
		}
	}

	if err != nil {
		plugin.Logger(ctx).Error("panos_security_rule.listSecurityRule", "query_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Debug("panos_security_rule.listSecurityRule", "len(listing)", len(listing))

	for _, i := range listing {
		plugin.Logger(ctx).Debug("panos_security_rule.listSecurityRule", "listing.i", i)
		d.StreamListItem(ctx, securityRuleStruct{vsys, deviceGroup, ruleBase, i})
	}

	return nil, nil
}