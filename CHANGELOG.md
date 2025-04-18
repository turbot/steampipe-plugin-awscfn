## v1.1.1 [2025-04-18]

_Bug fixes_

- Fixed Linux AMD64 plugin build failures for `Postgres 14 FDW`, `Postgres 15 FDW`, and `SQLite Extension` by upgrading GitHub Actions runners from `ubuntu-20.04` to `ubuntu-22.04`.

## v1.1.0 [2025-04-17]

_Dependencies_

- Recompiled plugin with Go version `1.23.1`. ([#55](https://github.com/turbot/steampipe-plugin-awscfn/pull/55))
- Recompiled plugin with [steampipe-plugin-sdk v5.11.5](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.11.5/CHANGELOG.md#v5115-2025-03-31) that addresses critical and high vulnerabilities in dependent packages. ([#55](https://github.com/turbot/steampipe-plugin-awscfn/pull/55))

## v1.0.0 [2024-10-22]

There are no significant changes in this plugin version; it has been released to align with [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin adheres to [semantic versioning](https://semver.org/#semantic-versioning-specification-semver), ensuring backward compatibility within each major version.

_Dependencies_

- Recompiled plugin with Go version `1.22`. ([#53](https://github.com/turbot/steampipe-plugin-awscfn/pull/53))
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. ([#53](https://github.com/turbot/steampipe-plugin-awscfn/pull/53))

## v0.7.1 [2023-12-12]

_Bug fixes_

- Fixed the missing optional tag on `paths` config parameter.

## v0.7.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#48](https://github.com/turbot/steampipe-plugin-awscfn/pull/48))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#48](https://github.com/turbot/steampipe-plugin-awscfn/pull/48))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-awscfn/blob/main/docs/LICENSE). ([#48](https://github.com/turbot/steampipe-plugin-awscfn/pull/48))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#47](https://github.com/turbot/steampipe-plugin-awscfn/pull/47))

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
