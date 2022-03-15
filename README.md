![image](https://hub.steampipe.io/images/plugins/turbot/awscfn-social-graphic.png)

# AWS CloudFormation Plugin for Steampipe

Use SQL to query data from AWS CloudFormation template files.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/awscfn)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/awscfn/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-awscfn/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install awscfn
```

Configure your [config file](https://hub.steampipe.io/plugins/turbot/awscfn#configuration) to include directories with AWS CloudFormation template files. If no directory is specified, the current working directory will be used.

Run steampipe:

```shell
steampipe query
```

Query all resources in your AWS CloudFormation template files:

```sql
select
  name,
  type,
  jsonb_pretty(properties) as resource_properties
from
  awscfn_resource;
```

```sh
> select name, type, jsonb_pretty(arguments) as args from awscfn_resource;
+-----------+-----------------+---------------------------------------+---------------------------------------+
| name      | type            | jsonb_pretty                          | jsonb_pretty                          |
+-----------+-----------------+---------------------------------------+---------------------------------------+
| DevBucket | AWS::S3::Bucket | {                                     | {                                     |
|           |                 |     "BucketName": "TestWebBucket",    |     "BucketName": {                   |
|           |                 |     "AccessControl": "PublicRead",    |         "Ref": "WebBucketName"        |
|           |                 |     "WebsiteConfiguration": {         |     },                                |
|           |                 |         "IndexDocument": "index.html" |     "AccessControl": "PublicRead",    |
|           |                 |     }                                 |     "WebsiteConfiguration": {         |
|           |                 | }                                     |         "IndexDocument": "index.html" |
|           |                 |                                       |     }                                 |
|           |                 |                                       | }                                     |
+-----------+-----------------+---------------------------------------+---------------------------------------+
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-awscfn.git
cd steampipe-plugin-awscfn
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```shell
make
```

Configure the plugin:

```shell
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/awscfn.spc
```

Try it!

```shell
steampipe query
> .inspect awscfn
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-awscfn/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [AWS CloudFormation Plugin](https://github.com/turbot/steampipe-plugin-awscfn/labels/help%20wanted)
