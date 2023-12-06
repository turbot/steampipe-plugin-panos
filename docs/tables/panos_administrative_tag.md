---
title: "Steampipe Table: panos_administrative_tag - Query Panorama Administrative Tags using SQL"
description: "Allows users to query Panorama Administrative Tags, providing insights into the administrative tags assigned to various resources within the Palo Alto Networks Panorama environment."
---

# Table: panos_administrative_tag - Query Panorama Administrative Tags using SQL

The Panorama Administrative Tags are a feature of Palo Alto Networks Panorama that allows administrators to assign tags to various resources within the environment. These tags can be utilized to filter and sort resources, making it easier to manage large and complex environments. They can also be used in policy rules to enforce security measures based on the assigned tags.

## Table Usage Guide

The `panos_administrative_tag` table provides insights into the Administrative Tags within the Palo Alto Networks Panorama. As a network administrator, explore tag-specific details through this table, including the tag name and color. Utilize it to uncover information about tags, such as those assigned to specific resources, allowing for better management and organization of the Panorama environment.

**Important Notes**
- For a full reference on mapping `color` codes to `color` names, please refer to: https://registry.terraform.io/providers/PaloAltoNetworks/panos/latest/docs/resources/administrative_tag.

## Examples

### Basic info
Explore which administrative tags are currently in use. This can help in managing and organizing system resources more effectively.

```sql
select
  name,
  color,
  comment
from
  panos_administrative_tag;
```

### List administrative tags containing a specific color
Explore which administrative tags are associated with a specific color. This can be useful in organizing and categorizing your tags based on color coding for easy identification and management.

```sql
select
  name,
  color,
  comment
from
  panos_administrative_tag
where
  color = 'color4';
```

### List administrative tags for a specific `vsys`
Explore the administrative tags associated with a specific system. This helps in identifying and organizing resources in a network for better management and control.

```sql
select
  name,
  color,
  comment
from
  panos_administrative_tag
where
  vsys = 'vsys1';
```

### List administrative tags for a **Panorama** `device group`
Explore the administrative tags associated with a specific device group in Panorama, allowing you to understand and manage the group's settings and comments.

```sql
select
  name,
  color,
  comment
from
  panos_administrative_tag
where
  device_group = 'group1';
```