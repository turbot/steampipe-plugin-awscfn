---
title: "Steampipe Table: awscfn_resource - Query AWS CloudFormation Resources using SQL"
description: "Allows users to query AWS CloudFormation Resources, specifically providing information about AWS resources in a stack, such as the logical and physical resource IDs and the type of resource."
---

# Table: awscfn_resource - Query AWS CloudFormation Resources using SQL

AWS CloudFormation is a service that helps you model and set up your Amazon Web Services resources. You can create a template that describes the AWS resources that you want to use. The service then takes care of provisioning and configuring those resources for you.

## Table Usage Guide

The `awscfn_resource` table provides insights into AWS resources in a stack. As a DevOps engineer, explore resource-specific details through this table, including the logical and physical resource IDs and the type of resource. Utilize it to uncover information about resources, such as their current status, stack ID, and the time when the resource was last updated.

The `properties_src` column contains the raw resource properties, while the `properties` column uses [AWS' goformation library](https://github.com/awslabs/goformation) to resolve CloudFormation instrinsic functions and references. In some cases, goformation is unable to parse the CloudFormation template or is unable to resolve property values.

For example, the sample [AutoScalingScheduledAction](https://s3.amazonaws.com/cloudformation-templates-us-east-1/AutoScalingScheduledAction.template) CloudFormation template includes a SecurityGroup resource:

```json
"InstanceSecurityGroup": {
  "Type": "AWS::EC2::SecurityGroup",
  "Properties": {
    "GroupDescription": "Enable SSH access and HTTP access on the configured port",
    "SecurityGroupIngress": [ {
      "IpProtocol": "tcp",
      "FromPort": "22",
      "ToPort": "22",
      "CidrIp": { "Ref": "SSHLocation" }
    }, {
      "IpProtocol": "tcp",
      "FromPort": "80",
      "ToPort": "80",
      "CidrIp": "0.0.0.0/0"
    } ],
    "VpcId": { "Ref" : "VpcId" }
  }
}
```

Because the `FromPort` and `ToPort` property values are of type `String` (which is valid per CloudFormation), but the [AWS::EC2::SecurityGroup Ingress schema](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-security-group-rule-1.html#cfn-ec2-security-group-rule-fromport) defines their type as `Integer`, goformation is unable to parse the CloudFormation template and `properties` will be returned as `null`:

```sql
select
  name,
  jsonb_pretty(properties_src) as properties_src,
  properties
from
  awscfn_resource
where
  name = 'InstanceSecurityGroup';
```

```sh
+-----------------------+-------------------------------------------------+------------+
| name                  | properties_src                                  | properties |
+-----------------------+-------------------------------------------------+------------+
| InstanceSecurityGroup | {                                               | <null>     |
|                       |     "VpcId": {                                  |            |
|                       |         "Ref": "VpcId"                          |            |
|                       |     },                                          |            |
|                       |     "GroupDescription": "Enable SSH access...", |            |
|                       |     "SecurityGroupIngress": [                   |            |
|                       |         {                                       |            |
|                       |             "CidrIp": {                         |            |
|                       |                 "Ref": "SSHLocation"            |            |
|                       |             },                                  |            |
|                       |             "ToPort": "22",                     |            |
|                       |             "FromPort": "22",                   |            |
|                       |             "IpProtocol": "tcp"                 |            |
|                       |         },                                      |            |
|                       |         {                                       |            |
|                       |             "CidrIp": "0.0.0.0/0",              |            |
|                       |             "ToPort": "80",                     |            |
|                       |             "FromPort": "80",                   |            |
|                       |             "IpProtocol": "tcp"                 |            |
|                       |         }                                       |            |
|                       |     ]                                           |            |
|                       | }                                               |            |
+-----------------------+-------------------------------------------------+------------+
```

## Examples

### Basic info
Analyze the settings of AWS CloudFormation resources to understand their types and configurations. This can be particularly useful to assess the elements within your infrastructure and their current status.

```sql+postgres
select
  name,
  type,
  case
    when properties is not null then properties
    else properties_src
  end as resource_properties,
  path
from
  awscfn_resource;
```

```sql+sqlite
select
  name,
  type,
  case
    when properties is not null then properties
    else properties_src
  end as resource_properties,
  path
from
  awscfn_resource;
```

### List AWS IAM users
Explore which AWS Identity and Access Management (IAM) users are active in your system. This provides a comprehensive view of user access, aiding in security and compliance management.

```sql+postgres
select
  name,
  type,
  case
    when properties is not null then properties
    else properties_src
  end as resource_properties,
  path
from
  awscfn_resource
where
  type = 'AWS::IAM::User';
```

```sql+sqlite
select
  name,
  type,
  case
    when properties is not null then properties
    else properties_src
  end as resource_properties,
  path
from
  awscfn_resource
where
  type = 'AWS::IAM::User';
```

### List AWS CloudTrail trails that are not encrypted
Explore which AWS CloudTrail trails lack encryption to enhance security measures. This helps in identifying potential vulnerabilities and ensuring compliance with security best practices.

```sql+postgres
select
  name,
  path
from
  awscfn_resource
where
  type = 'AWS::CloudTrail::Trail'
  and (
    ( properties is not null and properties -> 'KMSKeyId' is null )
    or properties_src -> 'KMSKeyId' is null
  );
```

```sql+sqlite
select
  name,
  path
from
  awscfn_resource
where
  type = 'AWS::CloudTrail::Trail'
  and (
    ( properties is not null and json_extract(properties, '$.KMSKeyId') is null )
    or json_extract(properties_src, '$.KMSKeyId') is null
  );
```

### Get S3 bucket BucketName property value
Determine the default name assigned to your AWS S3 bucket resources. This is useful for keeping track of your buckets and ensuring they are named according to your organizational standards.
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
  name as resource_map_name,
  type as resource_type,
  properties_src ->> 'BucketName' as bucket_name_src,
  default_value as bucket_name
from
  awscfn_resource
where
  type = 'AWS::S3::Bucket';
```

```sql+sqlite
select
  name as resource_map_name,
  type as resource_type,
  json_extract(properties_src, '$.BucketName') as bucket_name_src,
  default_value as bucket_name
from
  awscfn_resource
where
  type = 'AWS::S3::Bucket';
```

```sh
+---------------+-----------------+--------------------------+----------------+
| resource_name | resource_type   | bucket_name_src          | bucket_name    |
+---------------+-----------------+--------------------------+----------------+
| DevBucket     | AWS::S3::Bucket | {"Ref": "WebBucketName"} | TestWebBucket  |
+---------------+-----------------+--------------------------+----------------+
```