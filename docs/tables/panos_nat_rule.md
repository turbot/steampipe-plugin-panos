---
title: "Steampipe Table: panos_nat_rule - Query Palo Alto Networks NAT Rules using SQL"
description: "Allows users to query Network Address Translation (NAT) Rules in Palo Alto Networks, providing insights into rule configurations, source and destination interfaces, and translation types."
---

# Table: panos_nat_rule - Query Palo Alto Networks NAT Rules using SQL

A Network Address Translation (NAT) Rule in Palo Alto Networks is a policy that specifies how to translate the source and destination IP addresses of packets as they traverse through a device. These rules are essential for directing traffic correctly through firewalls and other network devices. They provide a means of modifying network address information in packet headers while in transit across a traffic routing device.

## Table Usage Guide

The `panos_nat_rule` table provides insights into NAT rules within Palo Alto Networks. As a network administrator, explore rule-specific details through this table, including source and destination interfaces, translation types, and associated metadata. Utilize it to gain a comprehensive understanding of your network's traffic routing and to ensure the correct configuration of your NAT rules.

## Examples

### List disabled NAT rules
Discover the segments that consist of disabled NAT rules. This can help you identify potential security loopholes in your network, thereby enhancing its safety and efficiency.

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
Explore the configuration of Network Address Translation (NAT) rules for a specific virtual system. This would be particularly useful for network administrators seeking to understand and manage the routing of network traffic within their system.

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
Explore which NAT rules are active for a particular device group in a Panorama setup. This can help in identifying potential security risks or troubleshooting network issues.

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
Identify the number of Network Address Translation (NAT) rules associated with each unique group tag. This can be useful for monitoring and managing network traffic routing configurations.

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
Explore which NAT rules are tagged with an administrative marker of yellow color. This is useful for identifying specific configurations that may require attention or follow a certain administrative pattern.

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
Uncover the details of NAT rules that facilitate packet transitions between distinct zones. This is useful in network management to identify potential areas of data flow and troubleshoot connectivity issues.

```sql
select
  name,
  source_zones,
  destination_zone
from
  panos_nat_rule
where
  not (source_zones ? destination_zone);
```

### List NAT rules which translate to unknown addresses
Determine the areas in which Network Address Translation (NAT) rules are translating to unidentified addresses. This is useful for identifying potential misconfigurations or security risks within your network infrastructure.

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