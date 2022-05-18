# Table: panos_security_rule

Find list of security rules present in the PAN-OS endpoint. Security policies allow you to enforce rules and take action, and can be as general or specific as needed. The policy rules are compared against the incoming traffic in sequence, and because the first rule that matches the traffic is applied, the more specific rules must precede the more general ones.

## Examples

### Basic ingress rule info

```sql
select
  name,
  type,
  action,
  source_zones,
  source_addresses,
  destination_zones,
  destination_addresses
  source_users
from
  panos_security_rule;
```

### List disabled security rules

```sql
select
  name,
  type,
  description
from
  panos_security_rule
where
  disabled;
```

### Get security rules count by group

```sql
select
  case
    when group_tag is null then 'none'
    else group_tag
  end as group_tag,
  count(*)
from
  panos_security_rule
group by group_tag;
```

### List security rules having public access to specific tagged addresses

```sql
with high_impact_tags as (
  select
    name
  from
    panos_administrative_tag
  where
    color = 'color1' -- red
),
address_with_high_impact_tags as (
  select
    a.name
  from
    panos_address_object as a
    join high_impact_tags as t on a.tags ? t.name
)
select
  r.name,
  r.source_addresses,
  r.destination_addresses
from
  panos_security_rule as r
  join address_with_high_impact_tags as ht on r.source_addresses ? 'any' and r.destination_addresses ? ht.name;
```

### List of security rules without `application` tag

```sql
select
  name,
  type,
  description,
  source_zones,
  destination_zones
from
  panos_security_rule
where
  not tags ? 'application';
```

### Lis security rules which contain any administrative tag with color yellow

```sql
with yellow_tags as (
  select
    name
  from
    panos_administrative_tag
  where
    color = 'color4' -- yellow
)
select
  panos_security_rule.name,
  panos_security_rule.type,
  panos_security_rule.description
from
  panos_security_rule
  join yellow_tags on panos_security_rule.tags ? yellow_tags.name;
```
