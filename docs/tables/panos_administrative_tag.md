# Table: panos_administrative_tag

Tags allows grouping of objects using keywords or phrases. Tags can be applied to address objects, address groups (static and dynamic), applications, zones, services, service groups, and to policy rules.

Each object can have up to 64 tags.

Tags can also be used to sort or filter objects and to visually distinguish objects by color. When you apply a color to a tag, the Policy tab displays the object with a background color.

For a full reference on mapping `color` codes to `color` names, please refer to: https://registry.terraform.io/providers/PaloAltoNetworks/panos/latest/docs/resources/administrative_tag.

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
