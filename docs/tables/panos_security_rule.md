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
