package panos

import (
	"context"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/nat"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tablePanosNATRule(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "panos_nat_rule",
		Description: "NAT rules for the PAN-OS device.",
		List: &plugin.ListConfig{
			Hydrate: listNATRule,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vsys", Require: plugin.Optional},
				{Name: "device_group", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
				{Name: "rule_base", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The NAT rule's name."},
			{Name: "uuid", Type: proto.ColumnType_STRING, Transform: transform.FromField("Uuid"), Description: "The PAN-OS UUID."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of NAT rule. This can be ipv4 (default), nat64, or nptv6.", Default: "ipv4"},
			{Name: "disabled", Type: proto.ColumnType_BOOL, Description: "Indicates if a rule is disabled, or not."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The NAT rule's description."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "List of administrative tags."},

			// Other columns
			{Name: "targets", Type: proto.ColumnType_JSON, Description: "Specifies a map of target definitions."},
			{Name: "negate_target", Type: proto.ColumnType_BOOL, Description: "Indicates if instead of applying the rule for the given serial numbers, it is applied to everything except them."},
			{Name: "group_tag", Type: proto.ColumnType_STRING, Description: "The NAT rule's group tag."},
			{Name: "source_zones", Type: proto.ColumnType_JSON, Description: "The list of source zone(s)."},
			{Name: "destination_zone", Type: proto.ColumnType_STRING, Description: "The NAT rule's destination zone."},
			{Name: "to_interface", Type: proto.ColumnType_STRING, Description: "Egress interface from route lookup (default: any)."},
			{Name: "service", Type: proto.ColumnType_STRING, Description: "Specifies the service (default: any)."},
			{Name: "source_addresses", Type: proto.ColumnType_JSON, Description: "A list of source addresses."},
			{Name: "destination_addresses", Type: proto.ColumnType_JSON, Description: "A list of destination address."},

			// Source Address Translation (SAT) config
			{Name: "sat_type", Type: proto.ColumnType_STRING, Description: "Type of source address translation. This can be none (default), dynamic-ip-and-port, dynamic-ip, or static-ip."},
			{Name: "sat_address_type", Type: proto.ColumnType_STRING, Description: "Source address translation address type. This can be interface-address or translated-address."},
			{Name: "sat_translated_addresses", Type: proto.ColumnType_JSON, Description: "The NAT rule's name."},
			{Name: "sat_interface", Type: proto.ColumnType_STRING, Description: "Describes the source address translation interface."},
			{Name: "sat_ip_address", Type: proto.ColumnType_INET, Description: "Describes source address translation IP address."},
			{Name: "sat_fallback_type", Type: proto.ColumnType_STRING, Description: "Specifies source address translation fallback type. This can be none, interface-address, or translated-address."},
			{Name: "sat_fallback_translated_addresses", Type: proto.ColumnType_JSON, Description: "A list of source address translation fallback translated address."},
			{Name: "sat_fallback_interface", Type: proto.ColumnType_STRING, Description: "Specifies source address translation interface."},
			{Name: "sat_fallback_ip_type", Type: proto.ColumnType_STRING, Description: "Specifies source address translation IP type. This can be ip or floating."},
			{Name: "sat_fallback_ip_address", Type: proto.ColumnType_INET, Description: "Specifies the source address translation fallback IP address."},
			{Name: "sat_static_translated_address", Type: proto.ColumnType_INET, Description: "Specifies the statically translated source address."},
			{Name: "sat_static_bi_directional", Type: proto.ColumnType_BOOL, Description: "Indicates whether bi-directional source address translation is enabled, or not."},

			// Destination Address Translation (DAT) config
			{Name: "dat_type", Type: proto.ColumnType_STRING, Description: "Specifies the destination address translation type. This should be either static or dynamic. The dynamic option is only available on PAN-OS 8.1+."},
			{Name: "dat_address", Type: proto.ColumnType_INET, Description: "Specifies destination address translation's address. Requires dat_type be set to \"static\" or \"dynamic\"."},
			{Name: "dat_port", Type: proto.ColumnType_INT, Description: "Specifies the destination address port. Requires dat_type be set to \"static\" or \"dynamic\"."},
			{Name: "dat_dynamic_distribution", Type: proto.ColumnType_STRING, Description: "Distribution algorithm for destination address pool. The PAN-OS 8.1 GUI doesn't seem to set this anywhere, but this is added here for completeness' sake. Requires dat_type of \"dynamic\"."},

			{Name: "vsys", Type: proto.ColumnType_STRING, Transform: transform.FromField("VSys").NullIfZero(), Description: "The vsys to put the NAT rule into (default: vsys1)."},
			{Name: "device_group", Type: proto.ColumnType_STRING, Description: "The device group location (default: shared)"},
			{Name: "rule_base", Type: proto.ColumnType_STRING, Description: "The rulebase. For firewalls, there is only the rulebase value (default), but on Panorama, there is also pre-rulebase and post-rulebase."},
			{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Description: "Raw view of data for the NAT rule."},
		},
	}
}

type natRuleStruct struct {
	VSys        string
	DeviceGroup string
	RuleBase    string
	nat.Entry
}

//// LIST FUNCTION

func listNATRule(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	plugin.Logger(ctx).Trace("listNATRule")

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("panos_nat_rule.listNATRule", "connection_error", err)
		return nil, err
	}

	// URL parameters for all queries
	keyQuals := d.KeyColumnQuals

	var vsys, deviceGroup, name string
	var natRules []nat.Entry
	var entry nat.Entry

	// Default set to rulebase
	ruleBase := util.Rulebase

	// Additional filters
	if d.KeyColumnQuals["name"] != nil {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	switch client := conn.(type) {
	case *pango.Firewall:
		{
			plugin.Logger(ctx).Debug("panos_nat_rule.listNATRule", "Firewall.id")
			vsys = "vsys1"
			if keyQuals["vsys"] != nil {
				vsys = keyQuals["vsys"].GetStringValue()
			}

			// Filter using name, if passed in qual
			if name != "" {
				entry, err = client.Policies.Nat.Get(vsys, name)
				natRules = append(natRules, entry)
			} else {
				natRules, err = client.Policies.Nat.GetAll(vsys)
			}
		}
	case *pango.Panorama:
		{
			plugin.Logger(ctx).Debug("panos_nat_rule.listNATRule", "Panorama.id")
			deviceGroup = "shared"
			if keyQuals["device_group"] != nil {
				deviceGroup = keyQuals["device_group"].GetStringValue()
			}

			// For Panorama, default set to pre_rulebase.
			// Override if passed in quals
			ruleBase = util.PreRulebase
			if keyQuals["rule_base"] != nil {
				ruleBase = keyQuals["rule_base"].GetStringValue()
			}

			// Filter using name, if passed in qual
			if name != "" {
				entry, err = client.Policies.Nat.Get(deviceGroup, ruleBase, name)
				natRules = append(natRules, entry)
			} else {
				natRules, err = client.Policies.Nat.GetAll(deviceGroup, ruleBase)
			}
		}
	}

	// Error handling
	if err != nil {
		plugin.Logger(ctx).Error("panos_nat_rule.listNATRule", "query_error", err)
		return nil, err
	}

	for _, i := range natRules {
		plugin.Logger(ctx).Debug("panos_nat_rule.listNATRule", "natRules.i", i)
		d.StreamListItem(ctx, natRuleStruct{vsys, deviceGroup, ruleBase, i})
	}

	return nil, nil
}
