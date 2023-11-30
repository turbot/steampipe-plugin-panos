---
title: "Steampipe Table: panos_address_object - Query OCI Panorama Address Objects using SQL"
description: "Allows users to query Panorama Address Objects, specifically the details of each address object, providing insights into network configurations."
---

# Table: panos_address_object - Query OCI Panorama Address Objects using SQL

A Panorama Address Object in OCI is a resource that represents an IP address, IP address range, or FQDN that you can use anywhere in a policy for source or destination address. It provides a way to group multiple IP addresses or FQDNs under a single name to streamline policy creation and management. These objects can be used to configure security policies and NAT rules in Panorama.

## Table Usage Guide

The `panos_address_object` table provides insights into address objects within OCI Panorama. As a network engineer, explore object-specific details through this table, including IP addresses, FQDNs, and associated descriptions. Utilize it to uncover information about objects, such as those with specific IP ranges, the association of multiple IPs under a single name, and the verification of policy configurations.

## Examples

### List all address objects
This query helps you explore all address objects in your network, allowing you to gain insights into the various network entities and their configurations. This is particularly useful in network management and security, where understanding the structure and organization of your network can aid in troubleshooting and threat mitigation.

```sql
select
  *
from
  panos_address_object;
```

### List all address objects ordered by object name
Explore all address objects in an organized manner, sorted by their respective names. This helps in efficient management and easy identification of individual address objects.

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