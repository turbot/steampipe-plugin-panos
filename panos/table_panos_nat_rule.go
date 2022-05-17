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
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The NAT rule's name."},
			{Name: "uuid", Type: proto.ColumnType_STRING, Description: "The PAN-OS UUID."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of NAT rule. This can be ipv4 (default), nat64, or nptv6.", Default: "ipv4"},
			{Name: "disabled", Type: proto.ColumnType_BOOL, Description: "Indicates if a rule is disabled, or not."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The NAT rule's description."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "List of administrative tags."},
			{Name: "target", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "negate_target", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "group_tag", Type: proto.ColumnType_STRING, Description: "The NAT rule's group tag."},
			{Name: "original_packet", Type: proto.ColumnType_JSON, Description: "The original packet specification."},
			{Name: "translated_packet", Type: proto.ColumnType_JSON, Description: "The translated packet specification."},
			{Name: "vsys", Type: proto.ColumnType_STRING, Transform: transform.FromField("VSys").NullIfZero(), Description: "The vsys to put the NAT rule into (default: vsys1)."},
			{Name: "device_group", Type: proto.ColumnType_STRING, Description: "The device group location (default: shared)"},
			{Name: "rule_base", Type: proto.ColumnType_STRING, Description: "The rulebase. For firewalls, there is only the rulebase value (default), but on Panorama, there is also pre-rulebase and post-rulebase."},
			{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Description: "Raw view of data for the NAT rule."},
		},
	}
}

//// LIST FUNCTION

func listNATRule(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	plugin.Logger(ctx).Trace("listNATRule")

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("panos_nat_rule.listNATRule", "connection_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("panos_nat_rule.listNATRule", "conn", conn)

	// URL parameters for all queries
	keyQuals := d.KeyColumnQuals

	var vsys, deviceGroup, name string
	var listing []nat.Entry
	var entry nat.Entry

	// Default set to rulebase
	ruleBase := util.Rulebase

	if d.KeyColumnQuals["name"] != nil {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	switch client := conn.(type) {
	case *pango.Firewall:
		{
			vsys = "vsys1"
			if keyQuals["vsys"] != nil {
				vsys = keyQuals["vsys"].GetStringValue()
			}
			plugin.Logger(ctx).Debug("panos_nat_rule.listNATRule", "Firewall.id", vsys)

			if name != "" {
				plugin.Logger(ctx).Debug("panos_nat_rule.listNATRule", "Firewall.name", name)
				entry, err = client.Policies.Nat.Get(vsys, name)
				listing = []nat.Entry{entry}
			} else {
				listing, err = client.Policies.Nat.GetAll(vsys)
			}
		}
	case *pango.Panorama:
		{
			deviceGroup = "shared"
			if keyQuals["device_group"] != nil {
				deviceGroup = keyQuals["device_group"].GetStringValue()
			}
			plugin.Logger(ctx).Debug("panos_nat_rule.listNATRule", "Panorama.id", deviceGroup)

			// For Panorama, default set to pre_rulebase.
			// Override if passed in quals
			ruleBase = util.PreRulebase
			if keyQuals["rule_base"] != nil {
				ruleBase = keyQuals["rule_base"].GetStringValue()
			}
			plugin.Logger(ctx).Debug("panos_nat_rule.listNATRule", "Panorama.rule_base", ruleBase)

			if name != "" {

				plugin.Logger(ctx).Debug("panos_nat_rule.listNATRule", "Panorama.name", name)

				entry, err = client.Policies.Nat.Get(deviceGroup, ruleBase, name)
				listing = append(listing, entry)
			} else {
				listing, err = client.Policies.Nat.GetAll(deviceGroup, ruleBase)
			}
		}
	}

	if err != nil {
		plugin.Logger(ctx).Error("panos_nat_rule.listNATRule", "query_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Debug("panos_nat_rule.listNATRule", "len(listing)", len(listing))

	for _, i := range listing {
		plugin.Logger(ctx).Debug("panos_nat_rule.listNATRule", "listing.i", i)

		natRule := buildNATRuleResultObj(vsys, deviceGroup, ruleBase, i)
		d.StreamListItem(ctx, natRule)
	}

	return nil, nil
}

func buildNATRuleResultObj(vsys string, dg string, ruleBase string, o nat.Entry) map[string]interface{} {
	m := map[string]interface{}{
		"Name":         o.Name,
		"UUID":         o.Uuid,
		"Description":  o.Description,
		"Type":         o.Type,
		"Disabled":     o.Disabled,
		"Tags":         o.Tags,
		"Target":       o.Targets,
		"NegateTarget": o.NegateTarget,
		"GroupTag":     o.GroupTag,
		"VSys":         vsys,
		"DeviceGroup":  dg,
		"RuleBase":     ruleBase,
	}

	op := map[string]interface{}{
		"source_zones":          o.SourceZones,
		"destination_zone":      o.DestinationZone,
		"destination_interface": o.ToInterface,
		"service":               o.Service,
		"source_addresses":      o.SourceAddresses,
		"destination_addresses": o.DestinationAddresses,
	}
	m["OriginalPacket"] = []interface{}{op}

	tp := make(map[string]interface{})
	src := make(map[string]interface{})
	dst := make(map[string]interface{})
	switch o.SatType {
	case nat.DynamicIpAndPort:
		diap := make(map[string]interface{})
		switch o.SatAddressType {
		case nat.TranslatedAddress:
			diap["translated_address"] = []interface{}{
				map[string]interface{}{
					"translated_addresses": o.SatTranslatedAddresses,
				},
			}
		case nat.InterfaceAddress:
			diap["interface_address"] = []interface{}{
				map[string]interface{}{
					"interface":  o.SatInterface,
					"ip_address": o.SatIpAddress,
				},
			}
		}
		src["dynamic_ip_and_port"] = []interface{}{diap}
	case nat.DynamicIp:
		di := map[string]interface{}{
			"translated_addresses": o.SatTranslatedAddresses,
		}
		switch o.SatFallbackType {
		case nat.TranslatedAddress:
			di["fallback"] = []interface{}{
				map[string]interface{}{
					"translated_address": []interface{}{
						map[string]interface{}{
							"translated_addresses": o.SatFallbackTranslatedAddresses,
						},
					},
				},
			}
		case nat.InterfaceAddress:
			di["fallback"] = []interface{}{
				map[string]interface{}{
					"interface_address": []interface{}{
						map[string]interface{}{
							"interface":  o.SatFallbackInterface,
							"type":       o.SatFallbackIpType,
							"ip_address": o.SatFallbackIpAddress,
						},
					},
				},
			}
		case nat.None:
			di["fallback"] = []interface{}{}
		}
		src["dynamic_ip"] = []interface{}{di}
	case nat.StaticIp:
		src["static_ip"] = []interface{}{
			map[string]interface{}{
				"translated_address": o.SatStaticTranslatedAddress,
				"bi_directional":     o.SatStaticBiDirectional,
			},
		}
	}
	switch o.DatType {
	case nat.DatTypeStatic:
		dst["static_translation"] = []interface{}{
			map[string]interface{}{
				"address": o.DatAddress,
				"port":    o.DatPort,
			},
		}
	case nat.DatTypeDynamic:
		dst["dynamic_translation"] = []interface{}{
			map[string]interface{}{
				"address":      o.DatAddress,
				"port":         o.DatPort,
				"distribution": o.DatDynamicDistribution,
			},
		}
	}
	tp["source"] = []interface{}{src}
	tp["destination"] = []interface{}{dst}
	m["TranslatedPacket"] = []interface{}{tp}

	return m
}
