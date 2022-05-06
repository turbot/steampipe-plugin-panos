package panos

import (
	"context"

	"github.com/PaloAltoNetworks/pango"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tablePanosAddressObject(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "panos_address_object",
		Description: "Address objects for the PAN-OS device.",
		List: &plugin.ListConfig{
			Hydrate: listAddressObject,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vsys", Require: plugin.Optional},
				{Name: "device_group", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The address object's name."},
			{Name: "vsys", Type: proto.ColumnType_STRING, Description: "[NGFW] The vsys the address object belongs to (default: vsys1)."},
			{Name: "device_group", Type: proto.ColumnType_STRING, Description: "[Panorama] The device group location (default: shared)"},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of address object."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "The address object's value."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The address object's description."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "List of administrative tags."},
			{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Description: "Raw view of data for the address object."},
		},
	}
}

func listAddressObject(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	plugin.Logger(ctx).Debug("panos_address_object.listAddressObject", "step", "about to connect")

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("panos_address_object.listAddressObject", "connection_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Debug("panos_address_object.listAddressObject", "conn", conn)

	// URL parameters for all queries
	keyQuals := d.KeyColumnQuals
	var id string
	var listing []string

	switch client := conn.(type) {
	case *pango.Firewall:
		{
			id = "vsys1"
			if keyQuals["vsys"] != nil {
				id = keyQuals["vsys"].GetStringValue()
			}
			plugin.Logger(ctx).Debug("panos_address_object.listAddressObject", "Firewall.id", id)
			listing, err = client.Objects.Address.GetList(id)
		}
	case *pango.Panorama:
		{
			id = "shared"
			if keyQuals["device_group"] != nil {
				id = keyQuals["shared"].GetStringValue()
			}
			plugin.Logger(ctx).Debug("panos_address_object.listAddressObject", "Panorama.id", id)
			listing, err = client.Objects.Address.GetList(id)
		}
	}

	if err != nil {
		plugin.Logger(ctx).Error("panos_address_object.listAddressObject", "query_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Debug("panos_address_object.listAddressObject", "len(listing)", len(listing))

	for _, i := range listing {
		plugin.Logger(ctx).Debug("panos_address_object.listAddressObject", "listing.i", i)
		d.StreamListItem(ctx, i)
	}

	return nil, nil
}
