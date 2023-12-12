---
title: "Steampipe Table: awscfn_mapping - Query AWS CloudFormation Mappings using SQL"
description: "Allows users to query Mappings in AWS CloudFormation, specifically the mapping key-value pairs defined in AWS CloudFormation templates, providing insights into the mapping function of AWS CloudFormation."
---

# Table: awscfn_mapping - Query AWS CloudFormation Mappings using SQL

AWS CloudFormation is a service that helps you model and set up your Amazon Web Services resources so you can spend less time managing those resources and more time focusing on your applications that run in AWS. You create a template that describes all the AWS resources that you want (like Amazon EC2 instances or Amazon RDS DB instances), and AWS CloudFormation takes care of provisioning and configuring those resources for you. The Mappings section in AWS CloudFormation templates enables you to create conditional parameter values during stack creation.

## Table Usage Guide

The `awscfn_mapping` table provides insights into Mappings within AWS CloudFormation. As a developer or system administrator, explore mapping-specific details through this table, including mapping key-value pairs defined in AWS CloudFormation templates. Utilize it to uncover information about mappings, such as those with specific conditions, the relationships between mappings, and the verification of mapping functions.

## Examples

For all examples below, assume we're using a CloudFormation template with the following `Mappings` section:

```yaml
Mappings:
  RegionMap:
    us-east-1:
      "HVM64": "ami-0ff8a91507f77f867"
    us-west-1:
      "HVM64": "ami-0bdb828fd58c52235"
    eu-west-1:
      "HVM64": "ami-047bb4163c506cd98"
    ap-southeast-1:
      "HVM64": "ami-08569b978cc4dfa10"
    ap-northeast-1:
      "HVM64": "ami-06cd52961ce9f0d85"
```

### Basic info
Explore the configuration details of AWS CloudFormation mapping to gain insights into key-value pairs and their paths. This can be useful for understanding the structure and relationships within your AWS CloudFormation templates.

```sql+postgres
select
  map,
  key,
  name,
  value,
  path
from
  awscfn_mapping;
```

```sql+sqlite
select
  map,
  key,
  name,
  value,
  path
from
  awscfn_mapping;
```

### List all HVM64 AMI IDs
Explore the specific regions where Amazon Machine Images (AMIs) with 64-bit hardware-assisted virtualization are being used. This is beneficial in managing resources and ensuring optimal performance across different regions.

```sql+postgres
select
  map,
  key,
  name,
  value as hvm64_ami_id,
  path
from
  awscfn_mapping
where
  map = 'RegionMap'
  and name = 'HVM64';
```

```sql+sqlite
select
  map,
  key,
  name,
  value as hvm64_ami_id,
  path
from
  awscfn_mapping
where
  map = 'RegionMap'
  and name = 'HVM64';
```

### Get the HVM64 AMI ID in us-east-1
Explore the specific Amazon Machine Image (AMI) identifier for 64-bit virtual machines in the US East (N. Virginia) region. This can help in understanding the resources available for cloud computing in that region.

```sql+postgres
select
  map,
  key,
  name,
  value as hvm64_ami_id,
  path
from
  awscfn_mapping
where
  map = 'RegionMap'
  and key = 'us-east-1'
  and name = 'HVM64';
```

```sql+sqlite
select
  map,
  key,
  name,
  value as hvm64_ami_id,
  path
from
  awscfn_mapping
where
  map = 'RegionMap'
  and key = 'us-east-1'
  and name = 'HVM64';
```

### Get the region whose HVM64 AMI ID is "ami-0bdb828fd58c52235"
Explore which regions are associated with a specific Amazon Machine Image (AMI) ID. This can be useful in identifying where certain resources are being utilized, aiding in resource allocation and management.

```sql+postgres
select
  map,
  key,
  name,
  value as hvm64_ami_id,
  path
from
  awscfn_mapping
where
  map = 'RegionMap'
  and name = 'HVM64'
  and value = 'ami-0bdb828fd58c52235';
```

```sql+sqlite
select
  map,
  key,
  name,
  value as hvm64_ami_id,
  path
from
  awscfn_mapping
where
  map = 'RegionMap'
  and name = 'HVM64'
  and value = 'ami-0bdb828fd58c52235';
```