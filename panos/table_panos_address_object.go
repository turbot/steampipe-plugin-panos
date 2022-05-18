package panos

import (
	"context"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/addr"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tablePanosAddressObject(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "panos_address_object",
		Description: "Address objects in the PAN-OS endpoint.",
		List: &plugin.ListConfig{
			Hydrate: listAddressObject,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vsys", Require: plugin.Optional},
				{Name: "device_group", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The address object's name."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of address - ip-netmask (default) | ip-range | ip-wildcard | fqdn."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "IP address or other value of the object."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of this object."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: " Administrative tags."},

			{Name: "vsys", Type: proto.ColumnType_STRING, Transform: transform.FromField("VSys").NullIfZero(), Description: "The vsys the address object belongs to (default: vsys1)."},
			{Name: "device_group", Type: proto.ColumnType_STRING, Transform: transform.FromField("DeviceGroup").NullIfZero(), Description: "The device group location (default: shared)."},
			{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Description: "Raw view of data for the address object."},
		},
	}
}

type addressStruct struct {
	VSys        string
	DeviceGroup string
	addr.Entry
}

//// LIST FUNCTION

func listAddressObject(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("panos_address_object.listAddressObject", "connection_error", err)
		return nil, err
	}

	// URL parameters for all queries
	keyQuals := d.KeyColumnQuals

	var vsys, deviceGroup, name string
	var listing []addr.Entry
	var entry addr.Entry

	// Additional filters
	if d.KeyColumnQuals["name"] != nil {
		name = d.KeyColumnQuals["name"].GetStringValue()
		plugin.Logger(ctx).Trace("panos_address_object.listAddressObject", "using name qual", name)
	}

	switch client := conn.(type) {
	case *pango.Firewall:
		{
			vsys = "vsys1"
			if keyQuals["vsys"] != nil {
				plugin.Logger(ctx).Trace("panos_address_object.listAddressObject", "Firewall", "using vsys qual")
				vsys = keyQuals["vsys"].GetStringValue()
			}
			plugin.Logger(ctx).Trace("panos_address_object.listAddressObject", "Firewall.vsys", vsys)

			// Filter using name, if passed in qual
			if name != "" {
				entry, err = client.Objects.Address.Get(vsys, name)
				listing = []addr.Entry{entry}
			} else {
				listing, err = client.Objects.Address.GetAll(vsys)
			}
		}
	case *pango.Panorama:
		{
			deviceGroup = "shared"
			if keyQuals["device_group"] != nil {
				plugin.Logger(ctx).Trace("panos_address_object.listAddressObject", "Panorama", "using device_group qual")
				deviceGroup = keyQuals["device_group"].GetStringValue()
			}
			plugin.Logger(ctx).Trace("panos_address_object.listAddressObject", "Panorama.device_group", deviceGroup)

			// Filter using name, if passed in qual
			if name != "" {
				entry, err = client.Objects.Address.Get(deviceGroup, name)
				listing = []addr.Entry{entry}
			} else {
				listing, err = client.Objects.Address.GetAll(deviceGroup)
			}
		}
	}

	if err != nil {
		plugin.Logger(ctx).Error("panos_address_object.listAddressObject", "query_error", err)
		return nil, err
	}

	for _, i := range listing {
		d.StreamListItem(ctx, addressStruct{vsys, deviceGroup, i})
	}

	return nil, nil
}
