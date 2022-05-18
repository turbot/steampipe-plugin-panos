# Table: panos_address_object

Address objects in the PAN-OS endpoint. An address object allows reusing the same address or group of addresses as a source or destination address in policy rules, filters, and other firewall functions without adding each address manually for each instance.

An address object can include either IPv4 or IPv6 addresses (a single IP address, a range of addresses, or a subnet), an FQDN, or a wildcard address (IPv4 address followed by a slash and wildcard mask).

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
  panos_address_object;
order by
  name;
```
