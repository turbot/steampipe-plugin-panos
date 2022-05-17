# Table: panos_administrative_tag

Lists the administrative tags in PAN-OS.

For a full reference on mapping `color` codes to `color` names, please refer to: https://registry.terraform.io/providers/PaloAltoNetworks/panos/latest/docs/resources/administrative_tag

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
  vsys = 'vsys1'
```

### List administrative tags for a **Panorama** device group

```sql
select
  name,
  color,
  comment
from
  panos_administrative_tag
where
  device_group = 'group1'
```

