# Table: awscfn_mapping

The Mappings section matches a key to a corresponding set of named values. For example, if you want to set values based on a region, you can create a mapping that uses the region name as a key and contains the values you want to specify for each specific region.

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

```sql
select
  map,
  key,
  name,
  value,
  path
from
  awscfn_mapping;
```

### Get the HVM64 AMI ID in us-east-1

```sql
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

```sql
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
