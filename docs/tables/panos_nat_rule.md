# Table: panos_nat_rule

NAT policies allow you to specify whether source or destination IP addresses and ports are converted between public and private addresses and ports. For example, private source addresses can be translated to public addresses on traffic sent from an internal (trusted) zone to a public (untrusted) zone. NAT is also supported on virtual wire interfaces.

NAT rules are based on source and destination zones, source and destination addresses, and application service (such as HTTP). Like security policies, NAT policy rules are compared against incoming traffic in sequence, and the first rule that matches the traffic is applied.

## Examples

### List disabled NAT rules

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

### Get count of NAT rules by distinct `group tag`

```sql
select
  case
    when group_tag is null then 'none'
    else group_tag
  end as group_tag,
  count(*) as count
from
  panos_nat_rule
group by
  group_tag;
```

### List NAT rules which contain any administrative tag with color yellow

```sql
with yellow_tags as (
  select
    name
  from
    panos_administrative_tag
  where
    color='color4' -- color4 :: Yellow
)
select
  panos_nat_rule.name,
  panos_nat_rule.type,
  panos_nat_rule.description
from
  panos_nat_rule
  join yellow_tags on panos_nat_rule.tags ? yellow_tags.name;
```

### List NAT rules which move packets between different zones

```sql
select
  *
from
  panos_nat_rule
where
  not (source_zones ? destination_zone);
```

### List NAT rules which translate to unknown addresses

```sql
select
  name,
  dat_address
from
  panos_nat_rule
where
  dat_address not in (
    select name from panos_address_object
  );
```
