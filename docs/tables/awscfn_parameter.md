# Table: awscfn_parameter

Parameters enable you to input custom values to your template each time you create or update a stack.

## Examples

### Basic info

```sql
select
  name,
  type,
  default_value,
  path
from
  awscfn_parameter;
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
  r.name as resource_map_name,
  r.type as resource_type,
  r.properties_src ->> 'BucketName' as bucket_reference,
  p.default_value as referred_value
from
  awscfn_resource as r,
  awscfn_parameter as p
where
  p.name = properties_src -> 'BucketName' ->> 'Ref'
  and r.type = 'AWS::S3::Bucket';
```

```sh
+-------------------+-----------------+--------------------------+----------------+
| resource_map_name | resource_type   | bucket_reference         | referred_value |
+-------------------+-----------------+--------------------------+----------------+
| DevBucket         | AWS::S3::Bucket | {"Ref": "WebBucketName"} | TestWebBucket  |
+-------------------+-----------------+--------------------------+----------------+
```

### List parameters with no default value configured

```sql
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
