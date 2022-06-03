package panos

import (
	"context"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/tags"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tablePanosAdministrativeTag(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "panos_administrative_tag",
		Description: "Tag objects in the PAN-OS endpoint.",
		List: &plugin.ListConfig{
			Hydrate: listPanosAdministrativeTag,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vsys", Require: plugin.Optional},
				{Name: "device_group", Require: plugin.Optional},
				{Name: "name", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "An unique identifier of the tag."},
			{Name: "color", Type: proto.ColumnType_STRING, Description: "Specifies the color ID for the tag."},
			{Name: "comment", Type: proto.ColumnType_STRING, Description: "Specifies a label or description to describe for what the tag is used."},

			{Name: "vsys", Type: proto.ColumnType_STRING, Transform: transform.FromField("VSys").NullIfZero(), Description: "[NGFW] The vsys the address object belongs to (default: vsys1)."},
			{Name: "device_group", Type: proto.ColumnType_STRING, Description: "[Panorama] The device group location (default: shared)."},
		},
	}
}

type tagStruct struct {
	VSys        string
	DeviceGroup string
	tags.Entry
}

//// LIST FUNCTION

func listPanosAdministrativeTag(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("panos_administrative_tag.listPanosAdministrativeTag", "connection_error", err)
		return nil, err
	}

	// URL parameters for all queries
	keyQuals := d.KeyColumnQuals

	var vsys, deviceGroup, name string
	var listing []tags.Entry
	var entry tags.Entry

	// Additional filter
	if d.KeyColumnQuals["name"] != nil {
		name = d.KeyColumnQuals["name"].GetStringValue()
		plugin.Logger(ctx).Trace("panos_administrative_tag.listPanosAdministrativeTag", "using name qual", name)
	}

	switch client := conn.(type) {
	case *pango.Firewall:
		{
			vsys = "vsys1"
			if keyQuals["vsys"] != nil {
				plugin.Logger(ctx).Trace("panos_administrative_tag.listPanosAdministrativeTag", "Firewall", "using vsys qual")
				vsys = keyQuals["vsys"].GetStringValue()
			}
			plugin.Logger(ctx).Trace("panos_administrative_tag.listPanosAdministrativeTag", "Firewall.vsys", vsys)

			// Filter using name, if passed in qual
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
				plugin.Logger(ctx).Trace("panos_administrative_tag.listPanosAdministrativeTag", "Panorama", "using device_group qual")
				deviceGroup = keyQuals["device_group"].GetStringValue()
			}
			plugin.Logger(ctx).Trace("panos_administrative_tag.listPanosAdministrativeTag", "Panorama.device_group", deviceGroup)

			// Filter using name, if passed in qual
			if name != "" {
				entry, err = client.Objects.Tags.Get(deviceGroup, name)
				listing = []tags.Entry{entry}
			} else {
				listing, err = client.Objects.Tags.GetAll(deviceGroup)
			}
		}
	}

	if err != nil {
		plugin.Logger(ctx).Error("panos_administrative_tag.listPanosAdministrativeTag", "query_error", err)
		return nil, err
	}

	for _, i := range listing {
		d.StreamListItem(ctx, tagStruct{vsys, deviceGroup, i})
	}

	return nil, nil
}
