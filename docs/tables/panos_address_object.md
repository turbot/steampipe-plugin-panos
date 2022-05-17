# Table: panos_group

Address objects in the PAN-OS endpoint.

## Examples

### List all address objects

```sql
select
  *
from
  panos_address_object;
```

### List all address objects ordered by object name

```sql
select
  name,
  value,
  description
from
  panos_address_object
order by name;
```

```sh
+----------+-----------------+-------------------------+
| name     | value           | description             |
+----------+-----------------+-------------------------+
| localnet | 192.168.80.0/24 | The 192.168.80 network. |
+----------+-----------------+-------------------------+
```
