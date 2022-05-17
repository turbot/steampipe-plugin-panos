# Table: panos_nat_rule

Network Address Translation (NAT) allows the source or destination IP address to be changed for traffic to transition through a router or gateway.

## Examples

### List disable NAT rules

```sql
select
  name,
  uuid,
  type
from
  panos_nat_rule
where
  disabled;
```

### List NAT rules for a specific `vsys`

```sql
select
  name,
  uuid,
  type,
  disabled,
  tags
from
  panos_nat_rule
where
  vsys = 'vsys1';
```

### List NAT rules for a **Panorama** device group

```sql
select
  name,
  uuid,
  type,
  disabled,
  tags
from
  panos_nat_rule
where
  device_group = 'group1';
```

### List of NAT rules without `application` tag

```sql
select
  name,
  uuid,
  type,
  disabled,
  tags
from
  panos_nat_rule
where
  tags is null
  or not tags ?| array['application'];
```

### Get NAT rules count by group

```sql
select
  case
    when group_tag is null then 'none'
    else group_tag
  end as group_tag,
  count(*)
from
  panos_nat_rule
group by group_tag;
```

### Lis NAT rules which contain any administrative tag with color yellow

```sql
select
  r.name,
  r.uuid,
  r.type,
  t.name,
  t.color
from
  panos_nat_rule as r,
  panos_administrative_tag as t
where
  t.color = 'color4'
  and r.tags ?| array[t.name]
```
