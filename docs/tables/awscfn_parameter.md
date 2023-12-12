---
title: "Steampipe Table: awscfn_parameter - Query AWS CloudFormation Parameters using SQL"
description: "Allows users to query AWS CloudFormation Parameters, providing insights into the parameters used in the AWS CloudFormation service."
---

# Table: awscfn_parameter - Query AWS CloudFormation Parameters using SQL

AWS CloudFormation is a service that helps you model and set up your Amazon Web Services resources so you can spend less time managing those resources and more time focusing on your applications that run in AWS. You create a template that describes all the AWS resources that you want (like Amazon EC2 instances or Amazon RDS DB instances), and AWS CloudFormation takes care of provisioning and configuring those resources for you. You don't need to individually create and configure AWS resources and figure out what's dependent on what; AWS CloudFormation handles all of that.

## Table Usage Guide

The `awscfn_parameter` table provides insights into the parameters used in the AWS CloudFormation service. As a Cloud Engineer or DevOps professional, you can explore parameter-specific details through this table, including default values, descriptions, and types. Utilize it to understand the configuration and dependencies of your AWS resources, and to ensure that the parameters used in your AWS CloudFormation templates are correctly configured and secure.

## Examples

### Basic info
Discover the segments that utilize different AWS CloudFormation parameters, such as their names and types, to gain insights into their default values and the path where they're stored. This is useful in understanding the configuration and usage of different parameters within your AWS CloudFormation service.

```sql+postgres
select
  name,
  type,
  default_value,
  path
from
  awscfn_parameter;
```

```sql+sqlite
select
  name,
  type,
  default_value,
  path
from
  awscfn_parameter;
```

### List S3 buckets with BucketName properties that reference a parameter
Determine the areas in which S3 bucket properties are referencing a parameter. This can be useful in managing and organizing your AWS resources, by allowing you to identify any dependencies or links between your S3 buckets and other AWS parameters.
For instance, if a CloudFormation template is defined as:

```yaml
Parameters:
  WebBucketName:
    Type: String
    Default: 'TestWebBucket'
Resources:
  DevBucket:
    Type: "AWS::S3::Bucket"
    Condition: CreateDevBucket
    Properties:
      AccessControl: PublicRead
      BucketName: !Ref WebBucketName
      WebsiteConfiguration:
        IndexDocument: index.html
```


```sql+postgres
select
  r.name as resource_name,
  r.type as resource_type,
  r.properties_src ->> 'BucketName' as bucket_name_src,
  p.default_value as bucket_name
from
  awscfn_resource as r,
  awscfn_parameter as p
where
  p.name = properties_src -> 'BucketName' ->> 'Ref'
  and r.type = 'AWS::S3::Bucket';
```

```sql+sqlite
select
  r.name as resource_name,
  r.type as resource_type,
  json_extract(r.properties_src, '$.BucketName') as bucket_name_src,
  p.default_value as bucket_name
from
  awscfn_resource as r,
  awscfn_parameter as p
where
  p.name = json_extract(json_extract(r.properties_src, '$.BucketName'), '$.Ref')
  and r.type = 'AWS::S3::Bucket';
```

```sh
+---------------+-----------------+--------------------------+----------------+
| resource_name | resource_type   | bucket_name_src          | bucket_name    |
+---------------+-----------------+--------------------------+----------------+
| DevBucket     | AWS::S3::Bucket | {"Ref": "WebBucketName"} | TestWebBucket  |
+---------------+-----------------+--------------------------+----------------+
```

### List parameters with no default value configured
Determine the areas in which parameters are lacking a default setting. This is useful to identify potential areas of concern or oversight in your configuration.

```sql+postgres
select
  name,
  type,
  description,
  path
from
  awscfn_parameter
where
  default_value is null;
```

```sql+sqlite
select
  name,
  type,
  description,
  path
from
  awscfn_parameter
where
  default_value is null;
```