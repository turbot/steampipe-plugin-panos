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

func tablePanosNATRuleGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "panos_nat_rule_group",
		Description: "NAT rules for the PAN-OS device.",
		List: &plugin.ListConfig{
			Hydrate: listNATRuleGroup,
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
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of NAT rule."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The NAT rule's description."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "List of administrative tags."},
			{Name: "source_zones", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "source_addresses", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "negate_target", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "destination_zone", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "destination_addresses", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "service", Type: proto.ColumnType_STRING, Description: ""},
			{Name: "disabled", Type: proto.ColumnType_BOOL, Description: ""},
			{Name: "targets", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "vsys", Type: proto.ColumnType_STRING, Transform: transform.FromField("VSys").NullIfZero(), Description: "[NGFW] The vsys the NAT rule belongs to (default: vsys1)."},
			{Name: "device_group", Type: proto.ColumnType_STRING, Transform: transform.FromField("DeviceGroup").NullIfZero(), Description: "[Panorama] The device group location (default: shared)"},
			{Name: "rule_base", Type: proto.ColumnType_STRING, Transform: transform.FromField("RuleBase").NullIfZero(), Description: "The rulebase. This can be either pre-rulebase (default), rulebase, or post-rulebase."},
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

func listNATRuleGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	plugin.Logger(ctx).Debug("panos_nat_rule_group.listNATRuleGroup", "step", "about to connect")

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("panos_nat_rule_group.listNATRuleGroup", "connection_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("panos_nat_rule_group.listNATRuleGroup", "conn", conn)

	// URL parameters for all queries
	keyQuals := d.KeyColumnQuals
	var vsys, deviceGroup, name string
	var listing []nat.Entry
	var entry nat.Entry

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
			plugin.Logger(ctx).Debug("panos_nat_rule_group.listNATRuleGroup", "Firewall.id", vsys)

			if name != "" {
				plugin.Logger(ctx).Debug("panos_nat_rule_group.listNATRuleGroup", "Firewall.name", name)
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
			plugin.Logger(ctx).Debug("panos_nat_rule_group.listNATRuleGroup", "Panorama.id", deviceGroup)

			if name != "" {
				plugin.Logger(ctx).Debug("panos_nat_rule_group.listNATRuleGroup", "Panorama.name", name)

				entry, err = client.Policies.Nat.Get(deviceGroup, ruleBase, name)
				listing = append(listing, entry)
			} else {
				listing, err = client.Policies.Nat.GetAll(deviceGroup, ruleBase)
			}
		}
	}

	if err != nil {
		plugin.Logger(ctx).Error("panos_nat_rule_group.listNATRuleGroup", "query_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Debug("panos_nat_rule_group.listNATRuleGroup", "len(listing)", len(listing))

	for _, i := range listing {
		plugin.Logger(ctx).Debug("panos_nat_rule_group.listNATRuleGroup", "listing.i", i)
		d.StreamListItem(ctx, natRuleStruct{vsys, deviceGroup, ruleBase, i})
	}

	return nil, nil
}
