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
order by name;
```

```sh
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
  plugin = "panos"
  
  # hostname to connect to 
  # hostname = "127.0.0.1"
  
  # api key to use for connection
  # api_key = "secret"
  
  # Username/Password combination to use for the connection. Ignored if 'api_key' is set
  # username = "username"
  # password = "password"
  
  # Request timeout (in seconds) for calls to the endpoint. Defaults to 10. Increase this if the endpoint may return
  # a high number of resources
  # timeout = 10
}
```

Environment variables are also available as an alternate configuration method:

- `PANOS_HOSTNAME`
- `PANOS_API_KEY`
- `PANOS_USERNAME`
- `PANOS_PASSWORD`

> Note: If `api_key` or `PANOS_API_KEY` is used, then `username / PANOS_USERNAME` and `password / PANOS_PASSWORD` are ignored.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-panos
- Community: [Slack Channel](https://steampipe.io/community/join)
