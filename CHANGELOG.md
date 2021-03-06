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
