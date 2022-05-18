# Table: panos_nat_rule

NAT policies allow you to specify whether source or destination IP addresses and ports are converted between public and private addresses and ports. For example, private source addresses can be translated to public addresses on traffic sent from an internal (trusted) zone to a public (untrusted) zone. NAT is also supported on virtual wire interfaces.

NAT rules are based on source and destination zones, source and destination addresses, and application service (such as HTTP). Like security policies, NAT policy rules are compared against incoming traffic in sequence, and the first rule that matches the traffic is applied.

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
  not tags ? 'application';
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
  and r.tags ? t.name
```
