---
organization: Turbot
category: ["network"]
icon_url: "/images/plugins/turbot/panos.svg"
brand_color: "#FA582D"
display_name: PAN-OS
name: panos
description: Steampipe plugin to query PAN-OS firewalls, security policies and more.
og_description: Query PAN-OS with SQL! Open source CLI. No DB required.
og_image: "/images/plugins/turbot/panos-social-graphic.png"
---

# PAN-OS + Steampipe

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

[PAN-OS](https://docs.paloaltonetworks.com/pan-os) is the operating system for Palo Alto Networks NGFWs and Panorama.

Example query:
```sql
select
  name,
  value,
  description
from
  panos_address_object
order by
  name
```

```
+----------+-----------------+-------------------------+
| name     | value           | description             |
+----------+-----------------+-------------------------+
| localnet | 192.168.80.0/24 | The 192.168.80 network. |
+----------+-----------------+-------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/panos/tables)**

## Get started

### Install

Download and install the latest PAN-OS plugin:

```bash
steampipe plugin install panos
```

### Configuration

Installing the latest panos plugin will create a config file (`~/.steampipe/config/panos.spc`) with a single connection named `panos`:

```hcl
connection "panos" {
  plugin   = "panos"
  hostname = "127.0.0.1"
  api_key  = "secret"
}
```

* `hostname` - The hostname / IP address of PAN-OS.
* `api_key` - The API key for the firewall.

Environment variables are also available as an alternate configuration method:
* `PANOS_HOSTNAME`
* `PANOS_API_KEY`

## Get involved

* Open source: https://github.com/turbot/steampipe-plugin-panos
* Community: [Slack Channel](https://steampipe.io/community/join)