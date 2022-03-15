# Table: awscfn_mapping

The Mappings section matches a key to a corresponding set of named values. For example, if you want to set values based on a region, you can create a mapping that uses the region name as a key and contains the values you want to specify for each specific region.

## Examples

### Basic info

```sql
select
  name,
  key,
  value,
  path
from
  awscfn_mapping;
```

### Get the mapped HVM64 AMI in us-east-1

```sql
select
  name,
  key,
  value ->> 'HVM64' as image_id,
  path
from
  awscfn_mapping
where
  name = 'AWSRegionArch2AMI'
  and key = 'us-east-1';
```
