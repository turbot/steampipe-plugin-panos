# Table: panos_security_rule

Find list of security rules present in the PAN-OS endpoint.

## Examples

### Basic info

```sql
select
  name,
  type,
  description,
  source_zones,
  destination_zones
from
  panos_security_rule;
```

```
+------------+-----------+-------------+----------------+-------------------+
| name       | type      | description | source_zones   | destination_zones |
+------------+-----------+-------------+----------------+-------------------+
| turbot_ccu | universal | test policy | ["zone1"]      | ["zone2"]         |
+------------+-----------+-------------+----------------+-------------------+
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
