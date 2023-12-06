---
title: "Steampipe Table: panos_security_rule - Query Panorama Security Rules using SQL"
description: "Allows users to query Panorama Security Rules, specifically the rules that control network access by defining the source and destination addresses, application, and action (allow, deny, or drop)."
---

# Table: panos_security_rule - Query Panorama Security Rules using SQL

Panorama Security Rules are a feature within that control network access by defining the source and destination addresses, application, and action (allow, deny, or drop). These rules are critical for managing network security and ensuring that only authorized traffic can access certain resources. Panorama Security Rules can be configured in a variety of ways to meet the specific needs of your network.

## Table Usage Guide

The `panos_security_rule` table provides insights into Panorama Security Rules within. As a network administrator, explore rule-specific details through this table, including source and destination addresses, application, and action. Utilize it to manage network security and ensure that only authorized traffic can access certain resources.

## Examples

### Basic ingress rule info
Explore the specifics of your network's security rules, including the type, action, and details about source and destination zones and addresses. This can help you understand how your network's security is configured and identify potential vulnerabilities.

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
Uncover the details of inactive security rules to understand the system's potential vulnerabilities and areas for improvement.

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
Determine the number of security rules associated with each group to understand the level of security measures applied. This can help identify areas where security may be lacking or overly stringent.

```sql
select
  case
    when group_tag is null then 'none'
    else group_tag
  end as group_tag,
  count(*)
from
  panos_security_rule
group by
  group_tag;
```

### List security rules having public access to specific tagged addresses
Determine the areas in which security rules allow public access to addresses tagged with high impact. This can help identify potential security risks and tighten access controls where necessary.

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
Discover the segments that lack an 'application' tag in your security rules. This can help you identify potential gaps in your rule tagging, thereby improving your security rule management.

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
Explore which security rules contain administrative tags marked in yellow. This can be useful for quickly identifying and reviewing specific security policies in your network that are flagged with this color for administrative purposes.

```sql
with yellow_tags as (
  select
    name
  from
    panos_administrative_tag
  where
    color = 'color4'; -- yellow
)
select
  panos_security_rule.name,
  panos_security_rule.type,
  panos_security_rule.description
from
  panos_security_rule
  join yellow_tags on panos_security_rule.tags ? yellow_tags.name;
```