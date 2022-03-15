---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/awscfn.svg"
brand_color: "#FF9900"
display_name: "AWS CloudFormation"
short_name: "awscfn"
description: "Steampipe plugin to query data from AWS CloudFormation template files."
og_description: "Query AWS CloudFormation template files with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/awscfn-social-graphic.png"
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
  awscfn_resource;
```

```sh
> select name, type, jsonb_pretty(properties) as args from awscfn_resource;
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

## Documentation

- **[Table definitions & examples →](/plugins/turbot/awscfn/tables)**

## Get started

### Install

Download and install the latest AWS CloudFormation plugin:

```bash
steampipe plugin install awscfn
```

### Credentials

No credentials are required.

### Configuration

Installing the latest awscfn plugin will create a config file (`~/.steampipe/config/awscfn.spc`) with a single connection named `awscfn`:

```hcl
connection "awscfn" {
  plugin = "awscfn"

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

- Open source: https://github.com/turbot/steampipe-plugin-awscfn
- Community: [Slack Channel](https://steampipe.io/community/join)
