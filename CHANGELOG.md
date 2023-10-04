## v0.6.1 [2023-10-04]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#37](https://github.com/turbot/steampipe-plugin-awscfn/pull/37))

## v0.6.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#33](https://github.com/turbot/steampipe-plugin-awscfn/pull/33))
- Recompiled plugin with Go version `1.21`. ([#33](https://github.com/turbot/steampipe-plugin-awscfn/pull/33))

## v0.5.0 [2023-06-20]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.5.0/CHANGELOG.md#v550-2023-06-16) which significantly reduces API calls and boosts query performance, resulting in faster data retrieval. This update significantly lowers the plugin initialization time of dynamic plugins by avoiding recursing into child folders when not necessary. ([#24](https://github.com/turbot/steampipe-plugin-awscfn/pull/24))

## v0.4.0 [2023-04-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#22](https://github.com/turbot/steampipe-plugin-awscfn/pull/22))

## v0.3.1 [2022-12-09]

_Bug fixes_

- Cleanup examples in `docs/index.md`.

## v0.3.0 [2022-11-17]

_What's new?_

- Added support for retrieving AWS CloudFormation template files from remote Git repositories and S3 buckets. For more information, please see [Supported Path Formats](https://hub.steampipe.io/plugins/turbot/awscfn#supported-path-formats). ([#15](https://github.com/turbot/steampipe-plugin-awscfn/pull/15))
- Added file watching support for files included in the `paths` config argument. ([#15](https://github.com/turbot/steampipe-plugin-awscfn/pull/15))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.0.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v500-2022-11-16) which includes support for fetching remote files with go-getter and file watching. ([#15](https://github.com/turbot/steampipe-plugin-awscfn/pull/15))

## v0.2.0 [2022-09-28]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#14](https://github.com/turbot/steampipe-plugin-awscfn/pull/14))
- Recompiled plugin with Go version `1.19`. ([#14](https://github.com/turbot/steampipe-plugin-awscfn/pull/14))

## v0.1.1 [2022-04-28]

_Bug fixes_

- Fixed the `invalid image` type error when trying to install or update the plugin ([#10](https://github.com/turbot/steampipe-plugin-awscfn/pull/10))

## v0.1.0 [2022-04-27]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#7](https://github.com/turbot/steampipe-plugin-awscfn/pull/7))
- Added support for native Linux ARM and Mac M1 builds. ([#8](https://github.com/turbot/steampipe-plugin-awscfn/pull/8))

## v0.0.2 [2022-03-25]

_Enhancements_

- Added column `map` to the `awscfn_mapping` table
- The `name` column in the `awscfn_mapping` table now contains the key name from each name-value pair
- The `value` column in the `awscfn_mapping` table now contains the value from each name-value pair

## v0.0.1 [2022-03-24]

_What's new?_

- New tables added
  - [awscfn_mapping](https://hub.steampipe.io/plugins/turbot/awscfn/tables/awscfn_mapping)
  - [awscfn_output](https://hub.steampipe.io/plugins/turbot/awscfn/tables/awscfn_output)
  - [awscfn_parameter](https://hub.steampipe.io/plugins/turbot/awscfn/tables/awscfn_parameter)
  - [awscfn_resource](https://hub.steampipe.io/plugins/turbot/awscfn/tables/awscfn_resource)
