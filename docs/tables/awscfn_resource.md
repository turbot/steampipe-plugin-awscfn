# Table: awscfn_resource

Each resource block describes one or more AWS resources that you want to include in the stack, such as Amazon EC2 instances, DynamoDB tables, or Amazon S3 buckets.

**Note:** Resource properties in AWS CloudFormation template **must** be configured with proper value reference and data types as per [AWS CloudFormation Template](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-security-group-rule-1.html#cfn-ec2-security-group-rule-fromport) documentation; otherwise the column `properties` will return null.

For example, This sample [AutoScalingScheduledAction](https://s3.amazonaws.com/cloudformation-templates-us-east-1/AutoScalingScheduledAction.template) template has defined following configuration to create a SecurityGroup resource

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

where the properties `FromPort` and `ToPort` values have been defined as **string**. But as per [AWS CloudFormation Template](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-security-group-rule-1.html#cfn-ec2-security-group-rule-fromport) documentation, the value should be of type **integer**.

## Examples

### Basic info

```sql
select
  name,
  type,
  case
    when properties is null then properties_src
    else  properties
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
    when properties is null then properties_src
    else  properties
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

### Get custom input value for S3 bucket

For instance, if a template is defined as:

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
  properties_src ->> 'BucketName' as bucket_reference,
  properties ->> 'BucketName' as calculated_value
from
  awscfn_resource
where
  path = '/path/to/testBucket.template'
  and type = 'AWS::S3::Bucket';
```

```sh
+-------------------+-----------------+--------------------------+------------------+
| resource_map_name | resource_type   | bucket_reference         | calculated_value |
+-------------------+-----------------+--------------------------+------------------+
| DevBucket         | AWS::S3::Bucket | {"Ref": "WebBucketName"} | TestWebBucket    |
+-------------------+-----------------+--------------------------+------------------+
```
