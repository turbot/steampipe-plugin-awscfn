![image](https://hub.steampipe.io/images/plugins/turbot/awscloudformation-social-graphic.png)

# AWS CloudFormation Plugin for Steampipe

Use SQL to query data from AWS CloudFormation template files.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/awscloudformation)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/awscloudformation/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-awscloudformation/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install awscloudformation
```

Configure your [config file](https://hub.steampipe.io/plugins/turbot/awscloudformation#configuration) to include directories with AWS CloudFormation template files. If no directory is specified, the current working directory will be used.

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
git clone https://github.com/turbot/steampipe-plugin-awscloudformation.git
cd steampipe-plugin-awscloudformation
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```shell
make
```

Configure the plugin:

```shell
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/awscloudformation.spc
```

Try it!

```shell
steampipe query
> .inspect awscloudformation
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-awscloudformation/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [AWS CloudFormation Plugin](https://github.com/turbot/steampipe-plugin-awscloudformation/labels/help%20wanted)
