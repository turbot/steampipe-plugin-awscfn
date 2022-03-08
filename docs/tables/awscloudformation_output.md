# Table: awscloudformation_output

The Outputs section declares output values that you can import into other stacks (to create cross-stack references), return in response (to describe stack calls), or view on the AWS CloudFormation console. For example, you can output the S3 bucket name for a stack to make the bucket easier to find.

## Examples

### Basic info

```sql
select
  name,
  description,
  value,
  path
from
  awscloudformation_output;
```

### List sensitive outputs

```sql
with output_table as (
  select
    name,
    description,
    split_part(substring(value from '\w*Ref:*\w*'), ':', 2) as parameter_reference,
    path
  from
    awscloudformation_output
  where
    value like '%Ref:%'
)
select
  o.name,
  o.description,
  o.path
from
  output_table as o
  left join awscloudformation_parameter as p on p.name = o.parameter_reference and o.path = p.path
where
  p.no_echo;
```
