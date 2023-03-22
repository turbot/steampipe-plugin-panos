package panos

import (
	"context"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/nat"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tablePanosNATRule(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "panos_nat_rule",
		Description: "NAT rules for the PAN-OS endpoint.",
		List: &plugin.ListConfig{
			Hydrate: listPanosNATRule,
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
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of NAT rule. Possible values are: ipv4 (default), nat64, or nptv6.", Default: "ipv4"},
			{Name: "disabled", Type: proto.ColumnType_BOOL, Description: "Indicates if a rule is disabled, or not."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The NAT rule's description."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "A list of administrative tags assigned to the rule."},

			// Other columns
			{Name: "targets", Type: proto.ColumnType_JSON, Description: "A dictionary of target definitions."},
			{Name: "negate_target", Type: proto.ColumnType_BOOL, Description: "Indicates if instead of applying the rule for the given serial numbers, it is applied to everything except them."},
			{Name: "group_tag", Type: proto.ColumnType_STRING, Description: "The NAT rule's group tag."},
			{Name: "source_zones", Type: proto.ColumnType_JSON, Description: "A list of source zone(s)."},
			{Name: "destination_zone", Type: proto.ColumnType_STRING, Description: "The NAT rule's destination zone."},
			{Name: "to_interface", Type: proto.ColumnType_STRING, Description: "Egress interface from route lookup (default: any)."},
			{Name: "service", Type: proto.ColumnType_STRING, Description: "Specifies the service (default: any)."},
			{Name: "source_addresses", Type: proto.ColumnType_JSON, Description: "A list of source addresses."},
			{Name: "destination_addresses", Type: proto.ColumnType_JSON, Description: "A list of destination address."},

			// Source Address Translation (SAT) config
			{Name: "sat_type", Type: proto.ColumnType_STRING, Description: "Specifies the type of source address translation. Possible values are: none (default), dynamic-ip-and-port, dynamic-ip, or static-ip."},
			{Name: "sat_address_type", Type: proto.ColumnType_STRING, Description: "Source address translation address type. Possible values are: interface-address or translated-address."},
			{Name: "sat_translated_addresses", Type: proto.ColumnType_JSON, Description: "A list of translated address."},
			{Name: "sat_interface", Type: proto.ColumnType_STRING, Description: "Describes the source address translation interface."},
			{Name: "sat_ip_address", Type: proto.ColumnType_STRING, Description: "Describes source address translation IP address."},
			{Name: "sat_fallback_type", Type: proto.ColumnType_STRING, Description: "Specifies source address translation fallback type. Possible values are: none, interface-address, or translated-address."},
			{Name: "sat_fallback_translated_addresses", Type: proto.ColumnType_JSON, Description: "A list of source address translation fallback translated address."},
			{Name: "sat_fallback_interface", Type: proto.ColumnType_STRING, Description: "Specifies source address translation interface."},
			{Name: "sat_fallback_ip_type", Type: proto.ColumnType_STRING, Description: "Specifies source address translation IP type. Possible values are: ip or floating."},
			{Name: "sat_fallback_ip_address", Type: proto.ColumnType_STRING, Description: "Specifies the source address translation fallback IP address."},
			{Name: "sat_static_translated_address", Type: proto.ColumnType_STRING, Description: "Specifies the statically translated source address."},
			{Name: "sat_static_bi_directional", Type: proto.ColumnType_BOOL, Description: "Indicates whether bi-directional source address translation is enabled, or not."},

			// Destination Address Translation (DAT) config
			{Name: "dat_type", Type: proto.ColumnType_STRING, Description: "Specifies the destination address translation type. Values can be either static or dynamic. The dynamic option is only available on PAN-OS 8.1+."},
			{Name: "dat_address", Type: proto.ColumnType_STRING, Description: "Specifies destination address translation's address. Possible values are: static or dynamic."},
			{Name: "dat_port", Type: proto.ColumnType_INT, Description: "Specifies the destination address port."},
			{Name: "dat_dynamic_distribution", Type: proto.ColumnType_STRING, Description: "Specifies the distribution algorithm for destination address pool."},

			{Name: "vsys", Type: proto.ColumnType_STRING, Transform: transform.FromField("VSys").NullIfZero(), Description: "[NGFW] The vsys to put the NAT rule into (default: vsys1)."},
			{Name: "device_group", Type: proto.ColumnType_STRING, Description: "[Panorama] The device group location (default: shared)."},
			{Name: "rule_base", Type: proto.ColumnType_STRING, Description: "[Panorama] The rulebase. For firewalls, there is only the rulebase value (default), but on Panorama, there is also pre-rulebase and post-rulebase."},
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

func listPanosNATRule(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("panos_nat_rule.listPanosNATRule", "connection_error", err)
		return nil, err
	}

	// URL parameters for all queries
	keyQuals := d.EqualsQuals

	var vsys, deviceGroup, name string
	var natRules []nat.Entry
	var entry nat.Entry

	// Default set to rulebase
	ruleBase := util.Rulebase

	// Additional filters
	if d.EqualsQuals["name"] != nil {
		name = d.EqualsQuals["name"].GetStringValue()
		plugin.Logger(ctx).Trace("panos_address_object.listAddressObject", "using name qual", name)
	}

	switch client := conn.(type) {
	case *pango.Firewall:
		{
			vsys = "vsys1"
			if keyQuals["vsys"] != nil {
				plugin.Logger(ctx).Trace("panos_nat_rule.listPanosNATRule", "Firewall", "using vsys qual")
				vsys = keyQuals["vsys"].GetStringValue()
			}
			plugin.Logger(ctx).Trace("panos_nat_rule.listPanosNATRule", "Firewall.vsys", vsys)

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
			deviceGroup = "shared"
			if keyQuals["device_group"] != nil {
				plugin.Logger(ctx).Trace("panos_nat_rule.listPanosNATRule", "Panorama", "using device_group qual")
				deviceGroup = keyQuals["device_group"].GetStringValue()
			}
			plugin.Logger(ctx).Trace("panos_nat_rule.listPanosNATRule", "Panorama.device_group", deviceGroup)

			// For Panorama, default set to pre_rulebase.
			// Override if passed in quals
			ruleBase = util.PreRulebase
			if keyQuals["rule_base"] != nil {
				plugin.Logger(ctx).Trace("panos_nat_rule.listPanosNATRule", "Panorama", "using rule_base qual")
				ruleBase = keyQuals["rule_base"].GetStringValue()
			}
			plugin.Logger(ctx).Trace("panos_nat_rule.listPanosNATRule", "Panorama.rule_base", ruleBase)

			// Filter using name, if passed in qual
			if name != "" {
				entry, err = client.Policies.Nat.Get(deviceGroup, ruleBase, name)
				natRules = append(natRules, entry)
			} else {
				natRules, err = client.Policies.Nat.GetAll(deviceGroup, ruleBase)
			}
		}
	}

	if err != nil {
		plugin.Logger(ctx).Error("panos_nat_rule.listPanosNATRule", "query_error", err)
		return nil, err
	}

	for _, i := range natRules {
		d.StreamListItem(ctx, natRuleStruct{vsys, deviceGroup, ruleBase, i})
	}

	return nil, nil
}
