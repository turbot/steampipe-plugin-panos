package panos

import (
	"context"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tablePanosSecurityRule(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "panos_security_rule",
		Description: "Security rules for the PAN-OS endpoint.",
		List: &plugin.ListConfig{
			Hydrate: listPanosSecurityRule,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vsys", Require: plugin.Optional},
				{Name: "device_group", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
				{Name: "rule_base", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the rule."},
			{Name: "uuid", Type: proto.ColumnType_STRING, Transform: transform.FromField("Uuid"), Description: "The PAN-OS UUID."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of security rule. Default sets to universal. Other possibles values are: interzone, or intrazone.", Default: "universal"},
			{Name: "disabled", Type: proto.ColumnType_BOOL, Description: "Indicates whether this rule is disabled."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The security rule's description."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "A list of administrative tags assigned to the rule."},

			// Other columns
			{Name: "source_zones", Type: proto.ColumnType_JSON, Description: "A list of source zones."},
			{Name: "source_addresses", Type: proto.ColumnType_JSON, Description: "A list of source addresses."},
			{Name: "negate_source", Type: proto.ColumnType_BOOL, Description: "Indicates if the source is negated."},
			{Name: "source_users", Type: proto.ColumnType_JSON, Description: "A list of source users."},
			{Name: "hip_profiles", Type: proto.ColumnType_JSON, Description: "A list of HIP profiles."},
			{Name: "destination_zones", Type: proto.ColumnType_JSON, Description: "A list of destination zones."},
			{Name: "destination_addresses", Type: proto.ColumnType_JSON, Description: "A list of destination addresses."},
			{Name: "negate_destination", Type: proto.ColumnType_BOOL, Description: "Indicates if the destination is negated."},
			{Name: "applications", Type: proto.ColumnType_JSON, Description: "A list of applications."},
			{Name: "services", Type: proto.ColumnType_JSON, Description: "A list of services."},
			{Name: "categories", Type: proto.ColumnType_JSON, Description: "A list of categories."},
			{Name: "action", Type: proto.ColumnType_STRING, Description: "Specifies the action for the matched traffic. Possible values are: allow (default), deny, drop, reset-client, reset-server, or reset-both."},
			{Name: "log_setting", Type: proto.ColumnType_STRING, Description: "Log forwarding profile."},
			{Name: "log_start", Type: proto.ColumnType_BOOL, Description: "Log the start of the traffic flow."},
			{Name: "log_end", Type: proto.ColumnType_BOOL, Description: "Log the end of the traffic flow."},
			{Name: "schedule", Type: proto.ColumnType_STRING, Description: "The security rule schedule."},
			{Name: "icmp_unreachable", Type: proto.ColumnType_BOOL, Description: "Indicates whether ICMP is unreachable."},
			{Name: "disable_server_response_inspection", Type: proto.ColumnType_BOOL, Description: "Indicates if server response inspection is disabled."},
			{Name: "group_tag", Type: proto.ColumnType_STRING, Description: "Specifies the group tag assigned to this rule."},
			{Name: "targets", Type: proto.ColumnType_JSON, Description: "A dictionary of target definitions."},
			{Name: "negate_target", Type: proto.ColumnType_BOOL, Description: "Instead of applying the rule for the given serial numbers, it is applied to everything except them."},

			{Name: "group", Type: proto.ColumnType_STRING, Description: "The group profile name."},
			{Name: "virus", Type: proto.ColumnType_STRING, Description: "The antivirus setting."},
			{Name: "spyware", Type: proto.ColumnType_STRING, Description: "The anti-spyware setting."},
			{Name: "vulnerability", Type: proto.ColumnType_STRING, Description: "The vulnerability protection setting."},
			{Name: "url_filtering", Type: proto.ColumnType_STRING, Transform: transform.FromField("UrlFiltering").NullIfZero(), Description: "The URL filtering setting."},
			{Name: "file_blocking", Type: proto.ColumnType_STRING, Description: "The file blocking setting."},
			{Name: "wild_fire_analysis", Type: proto.ColumnType_STRING, Description: "The wildfire analysis setting."},
			{Name: "data_filtering", Type: proto.ColumnType_STRING, Description: "The data filtering setting."},

			{Name: "vsys", Type: proto.ColumnType_STRING, Transform: transform.FromField("VSys").NullIfZero(), Description: "[NGFW] The vsys the security rule belongs to (default: vsys1)."},
			{Name: "device_group", Type: proto.ColumnType_STRING, Description: "[Panorama] The device group location (default: shared)."},
			{Name: "rule_base", Type: proto.ColumnType_STRING, Description: "[Panorama] The rulebase. This can be either pre-rulebase (default for panorama), rulebase, or post-rulebase."},
		},
	}
}

type securityRuleStruct struct {
	VSys        string
	DeviceGroup string
	RuleBase    string
	security.Entry
}

//// LIST FUNCTION

func listPanosSecurityRule(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("panos_security_rule.listPanosSecurityRule", "connection_error", err)
		return nil, err
	}

	// URL parameters for all queries
	keyQuals := d.EqualsQuals

	var vsys, deviceGroup, name string
	var listing []security.Entry
	var entry security.Entry

	// Default to rulebase
	// Override if passed in quals
	ruleBase := util.PreRulebase
	if keyQuals["rule_base"] != nil {
		ruleBase = keyQuals["rule_base"].GetStringValue()
		plugin.Logger(ctx).Trace("panos_address_object.listAddressObject", "using rule_base qual", name)
	}

	// Additional filters
	if keyQuals["name"] != nil {
		name = keyQuals["name"].GetStringValue()
		plugin.Logger(ctx).Trace("panos_address_object.listAddressObject", "using name qual", name)
	}

	switch client := conn.(type) {
	case *pango.Firewall:
		{
			vsys = "vsys1"
			if keyQuals["vsys"] != nil {
				plugin.Logger(ctx).Trace("panos_security_rule.listPanosSecurityRule", "Firewall", "using vsys qual")
				vsys = keyQuals["vsys"].GetStringValue()
			}

			// Filter using name, if passed in qual
			if name != "" {
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
				plugin.Logger(ctx).Trace("panos_security_rule.listPanosSecurityRule", "Panorama", "using device_group qual")
				deviceGroup = keyQuals["device_group"].GetStringValue()
			}
			plugin.Logger(ctx).Trace("panos_security_rule.listPanosSecurityRule", "Panorama.device_group", deviceGroup)

			// Filter using name, if passed in qual
			if name != "" {
				entry, err = client.Policies.Security.Get(deviceGroup, ruleBase, name)
				listing = append(listing, entry)
			} else {
				listing, err = client.Policies.Security.GetAll(deviceGroup, ruleBase)
			}
		}
	}

	// Error handling
	if err != nil {
		plugin.Logger(ctx).Error("panos_security_rule.listPanosSecurityRule", "query_error", err)
		return nil, err
	}

	for _, i := range listing {
		d.StreamListItem(ctx, securityRuleStruct{vsys, deviceGroup, ruleBase, i})
	}

	return nil, nil
}
