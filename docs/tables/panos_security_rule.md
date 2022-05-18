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
  tags is null
  or not tags ? 'application';
```

### Lis security rules which contain any administrative tag with color yellow

```sql
select
  r.name,
  r.type,
  t.name,
  t.color
from
  panos_security_rule as r,
  panos_administrative_tag as t
where
  t.color = 'color4'
  and r.tags ? t.name;
```
