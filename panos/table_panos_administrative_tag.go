package panos

import (
	"context"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/tags"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tablePanosTagObject(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "panos_administrative_tag",
		Description: "Tag objects for the PAN-OS endpoint.",
		List: &plugin.ListConfig{
			Hydrate: listTag,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vsys", Require: plugin.Optional},
				{Name: "device_group", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The tag object's name."},
			{Name: "vsys", Type: proto.ColumnType_STRING, Transform: transform.FromField("VSys").NullIfZero(), Description: "[NGFW] The vsys the address object belongs to (default: vsys1)."},
			{Name: "device_group", Type: proto.ColumnType_STRING, Transform: transform.FromField("DeviceGroup").NullIfZero(), Description: "[Panorama] The device group location (default: shared)"},
			{Name: "color", Type: proto.ColumnType_STRING, Description: "The color of the Tag."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "Tag Comment."},
			{Name: "raw", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Description: "Raw view of data for the tag object."},
		},
	}
}

type tagStruct struct {
	VSys        string
	DeviceGroup string
	tags.Entry
}

func listTag(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	plugin.Logger(ctx).Debug("panos_tag.listTag", "step", "about to connect")

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("panos_tag.listTag", "connection_error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("panos_tag.listTag", "conn", conn)

	// URL parameters for all queries
	keyQuals := d.KeyColumnQuals
	var vsys, deviceGroup, name string
	var listing []tags.Entry
	var entry tags.Entry

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
			plugin.Logger(ctx).Debug("panos_tag.listTag", "Firewall.id", vsys)
			plugin.Logger(ctx).Debug("panos_tag.listTag", "Firewall.name", vsys)

			if name != "" {
				entry, err = client.Objects.Tags.Get(vsys, name)
				listing = []tags.Entry{entry}
			} else {
				listing, err = client.Objects.Tags.GetAll(vsys)
			}
		}
	case *pango.Panorama:
		{
			deviceGroup = "shared"
			if keyQuals["device_group"] != nil {
				deviceGroup = keyQuals["device_group"].GetStringValue()
			}
			plugin.Logger(ctx).Debug("panos_tag.listTag", "Panorama.id", deviceGroup)
			plugin.Logger(ctx).Debug("panos_tag.listTag", "Panorama.name", name)
			if name != "" {
				entry, err = client.Objects.Tags.Get(deviceGroup, name)
				listing = []tags.Entry{entry}
			} else {
				listing, err = client.Objects.Tags.GetAll(deviceGroup)
			}
		}
	}

	if err != nil {
		plugin.Logger(ctx).Error("panos_tag.listTag", "query_error", err)
		return nil, err
	}

	plugin.Logger(ctx).Debug("panos_tag.listTag", "len(listing)", len(listing))

	for _, i := range listing {
		plugin.Logger(ctx).Debug("panos_tag.listTag", "listing.i", i)
		d.StreamListItem(ctx, tagStruct{vsys, deviceGroup, i})
	}

	return nil, nil
}
