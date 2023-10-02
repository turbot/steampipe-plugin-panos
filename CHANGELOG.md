## v0.3.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#22](https://github.com/turbot/steampipe-plugin-panos/pull/22))
- Recompiled plugin with Go version `1.21`. ([#22](https://github.com/turbot/steampipe-plugin-panos/pull/22))

## v0.2.0 [2023-03-22]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#17](https://github.com/turbot/steampipe-plugin-panos/pull/17))

## v0.1.0 [2022-09-28]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#14](https://github.com/turbot/steampipe-plugin-panos/pull/14))
- Recompiled plugin with Go version `1.19`. ([#14](https://github.com/turbot/steampipe-plugin-panos/pull/14))

## v0.0.3 [2022-06-29]

_Bug fixes_

- Updated the plugin category from `network` to `security`. ([#10](https://github.com/turbot/steampipe-plugin-panos/pull/10))

## v0.0.2 [2022-06-03]

_Breaking changes_

- Removed the column `raw` from all tables since the data included in the column was redundant. ([#6](https://github.com/turbot/steampipe-plugin-panos/pull/6))

## v0.0.1 [2022-05-18]

_What's new?_

- New tables added

  - [panos_administrative_tag](https://hub.steampipe.io/plugins/turbot/panos/tables/panos_administrative_tag)
  - [panos_address_object](https://hub.steampipe.io/plugins/turbot/panos/tables/panos_address_object)
  - [panos_nat_rule](https://hub.steampipe.io/plugins/turbot/panos/tables/panos_nat_rule)
  - [panos_security_rule](https://hub.steampipe.io/plugins/turbot/panos/tables/panos_security_rule)
