---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/awscloudformation.svg"
brand_color: "#008000"
display_name: "AWS CloudFormation"
short_name: "awscloudformation"
description: "Steampipe plugin to query data from AWS CloudFormation template files."
og_description: "Query AWS CloudFormation template files with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/awscloudformation-social-graphic.png"
---

# AWS CloudFormation + Steampipe

An AWS CloudFormation template file is used to declare resources, variables, modules, and more.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query data using SQL.

Query all resources in your AWS CloudFormation files:

```sql
select
  name,
  type,
  jsonb_pretty(properties) as resource_properties
from
  awscloudformation_resource;
```

```sh
> select name, type, jsonb_pretty(arguments) as args from awscloudformation_resource;
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

## Documentation

- **[Table definitions & examples →](/plugins/turbot/awscloudformation/tables)**

## Get started

### Install

Download and install the latest AWS CloudFormation plugin:

```bash
steampipe plugin install awscloudformation
```

### Credentials

No credentials are required.

### Configuration

Installing the latest awscloudformation plugin will create a config file (`~/.steampipe/config/awscloudformation.spc`) with a single connection named `awscloudformation`:

```hcl
connection "awscloudformation" {
  plugin = "awscloudformation"

  # Paths is a list of locations to search for CloudFormation template files
  # All paths are resolved relative to the current working directory (CWD)
  # Wildcard based searches are supported, including recursive searches

  # For example:
  #  - "*.template" matches all CloudFormation template files in the CWD
  #  - "**/*.template" matches all CloudFormation template files in the CWD and all sub-directories
  #  - "../*.template" matches all CloudFormation template files in the CWD's parent directory
  #  - "ELB*.template" matches all CloudFormation template files starting with "ELB" in the CWD
  #  - "/path/to/dir/*.template" matches all CloudFormation template files in a specific directory
  #  - "/path/to/dir/main.template" matches a specific file

  # If paths includes "*", all files (including non-CloudFormation template files) in
  # the CWD will be matched, which may cause errors if incompatible file types exist

  # Defaults to CWD
  paths = ["*.template"]
}
```

- `paths` - A list of directory paths to search for AWS CloudFormation template files. Paths are resolved relative to the current working directory. Paths may [include wildcards](https://pkg.go.dev/path/filepath#Match) and also support `**` for recursive matching. Defaults to the current working directory.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-awscloudformation
- Community: [Slack Channel](https://steampipe.io/community/join)
