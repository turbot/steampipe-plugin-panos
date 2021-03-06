![image](https://hub.steampipe.io/images/plugins/turbot/panos-social-graphic.png)

# PAN-OS Plugin for Steampipe

Use SQL to query firewalls, security policies and more from [PAN-OS](https://docs.paloaltonetworks.com/pan-os).

- **[Get started →](https://hub.steampipe.io/plugins/turbot/panos)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/panos/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-panos/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install panos
```

Run steampipe:

```shell
steampipe query
```

Run a query:

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

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-panos.git
cd steampipe-plugin-panos
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```shell
make
```

Configure the plugin:

```shell
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/panos.spc
```

Try it!

```shell
steampipe query
> .inspect panos
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-prometheus/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [PAN-OS Plugin](https://github.com/turbot/steampipe-plugin-panos/labels/help%20wanted)
