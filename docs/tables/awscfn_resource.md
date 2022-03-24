# Table: awscfn_resource

Each resource block describes one or more AWS resources that you want to include in the stack, such as Amazon EC2 instances, DynamoDB tables, or Amazon S3 buckets.

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

```sql
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

```sql
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

```sql
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

### Get S3 bucket BucketName property value

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

```sql
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

```sh
+---------------+-----------------+--------------------------+----------------+
| resource_name | resource_type   | bucket_name_src          | bucket_name    |
+---------------+-----------------+--------------------------+----------------+
| DevBucket     | AWS::S3::Bucket | {"Ref": "WebBucketName"} | TestWebBucket  |
+---------------+-----------------+--------------------------+----------------+
```
