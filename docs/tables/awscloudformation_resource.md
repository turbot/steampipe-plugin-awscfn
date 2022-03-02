# Table: awscloudformation_resource

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
  awscloudformation_resource;
```

### List AWS IAM users

```sql
select
  name,
  type,
  properties,
  path
from
  awscloudformation_resource
where
  type = 'AWS::IAM::User';
```

### List AWS CloudTrail trails that are not encrypted

```sql
select
  name,
  path
from
  awscloudformation_resource
where
  type = 'AWS::CloudTrail::Trail'
  and properties -> 'KMSKeyId' is null;
```
