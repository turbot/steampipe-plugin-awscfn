# Table: awscfn_resource

Each resource block describes one or more AWS resources that you want to include in the stack, such as Amazon EC2 instances, DynamoDB tables, or Amazon S3 buckets.

## Examples

### Basic info

```sql
select
  name,
  type,
  properties,
  path
from
  awscfn_resource;
```

### List AWS IAM users

```sql
select
  name,
  type,
  properties,
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
  and properties -> 'KMSKeyId' is null;
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
