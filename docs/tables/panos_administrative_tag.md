# Table: panos_administrative_tag

Lists the administrative tags in PAN-OS.

## Tag Colors

Tag colors are as follows:

* `color1`: Red
* `color2`: Green
* `color3`: Blue
* `color4`: Yellow
* `color5`: Copper
* `color6`: Orange
* `color7`: Purple
* `color8`: Gray
* `color9`: Light Green
* `color10`: Cyan
* `color11`: Light Gray
* `color12`: Blue Gray
* `color13`: Lime
* `color14`: Black
* `color15`: Gold
* `color16`: Brown
* `color17`: Olive
* `color18`: (Reserved for internal use)
* `color19`: Maroon
* `color20`: Red Orange
* `color21`: Yellow Orange
* `color22`: Forest Green
* `color23`: Turquoise Blue
* `color24`: Azure Blue
* `color25`: Cerulean Blue
* `color26`: Midnight Blue
* `color27`: Medium Blue
* `color28`: Cobalt Blue
* `color29`: Violet Blue
* `color30`: Blue Violet
* `color31`: Medium Violet
* `color32`: Medium Rose
* `color33`: Lavender
* `color34`: Orchid
* `color35`: Thistle
* `color36`: Peach
* `color37`: Salmon
* `color38`: Magenta
* `color39`: Red Violet
* `color40`: Mahogany
* `color41`: Burnt Sienna
* `color42`: Chestnut

## Examples

### Basic info

```sql
select
  name,
  color,
  comment
from
  panos_administrative_tag;
```

### List administrative tags for a specific `vsys`

```sql
select
  name,
  color,
  comment
from
  panos_administrative_tag;
where
  vsys='vsys1'
```

### List administrative tags for a **Panorama** device group

```sql
select
  name,
  color,
  comment
from
  panos_administrative_tag;
where
  device_group='group1'
```

