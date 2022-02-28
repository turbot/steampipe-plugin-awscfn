![image](https://hub.steampipe.io/images/plugins/turbot/cloudformation-social-graphic.png)

# CloudFormation Plugin for Steampipe

Use SQL to query data from CloudFormation template files.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/cloudformation)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/cloudformation/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-cloudformation/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install cloudformation
```

Configure your [config file](https://hub.steampipe.io/plugins/turbot/cloudformation#configuration) to include directories with CloudFormation template files. If no directory is specified, the current working directory will be used.

Run steampipe:

```shell
steampipe query
```

Query all resources in your CloudFormation template files:

```sql
select
  name,
  type,
  jsonb_pretty(properties) as resource_properties
from
  cloudformation_resource;
```

```sh
> select name, type, jsonb_pretty(arguments) as args from cloudformation_resource;
+-----------------+----------------------+---------------------------------------------+
| name            | type                 | resource_properties                         |
+-----------------+----------------------+---------------------------------------------+
| myDynamoDBTable | AWS::DynamoDB::Table | {                                           |
|                 |                      |     "KeySchema": [                          |
|                 |                      |         {                                   |
|                 |                      |             "KeyType": "HASH",              |
|                 |                      |             "AttributeName": {              |
|                 |                      |                 "Ref": "HashKeyElementName" |
|                 |                      |             }                               |
|                 |                      |         }                                   |
|                 |                      |     ],                                      |
|                 |                      |     "AttributeDefinitions": [               |
|                 |                      |         {                                   |
|                 |                      |             "AttributeName": {              |
|                 |                      |                 "Ref": "HashKeyElementName" |
|                 |                      |             },                              |
|                 |                      |             "AttributeType": {              |
|                 |                      |                 "Ref": "HashKeyElementType" |
|                 |                      |             }                               |
|                 |                      |         }                                   |
|                 |                      |     ],                                      |
|                 |                      |     "ProvisionedThroughput": {              |
|                 |                      |         "ReadCapacityUnits": {              |
|                 |                      |             "Ref": "ReadCapacityUnits"      |
|                 |                      |         },                                  |
|                 |                      |         "WriteCapacityUnits": {             |
|                 |                      |             "Ref": "WriteCapacityUnits"     |
|                 |                      |         }                                   |
|                 |                      |     }                                       |
|                 |                      | }                                           |
+-----------------+----------------------+---------------------------------------------+
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-cloudformation.git
cd steampipe-plugin-cloudformation
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```shell
make
```

Configure the plugin:

```shell
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/cloudformation.spc
```

Try it!

```shell
steampipe query
> .inspect cloudformation
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-cloudformation/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [CloudFormation Plugin](https://github.com/turbot/steampipe-plugin-cloudformation/labels/help%20wanted)
