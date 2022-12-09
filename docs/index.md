---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/awscfn.svg"
brand_color: "#FF9900"
display_name: "AWS CloudFormation"
short_name: "awscfn"
description: "Steampipe plugin to query data from AWS CloudFormation template files."
og_description: "Query AWS CloudFormation template files with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/awscfn-social-graphic.png"
---

# AWS CloudFormation + Steampipe

An [AWS CloudFormation template file](https://aws.amazon.com/cloudformation/resources/templates/) is used to declare resources, variables, modules, and more.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query data using SQL.

Query all resources in your AWS CloudFormation files:

```sql
select
  name,
  type,
  jsonb_pretty(properties) as properties
from
  awscfn_resource;
```

```sh
+------------+------------------+---------------------------------------+
| name       | type             | properties                            |
+------------+------------------+---------------------------------------+
| CFNUser    | AWS::IAM::User   | {                                     |
|            |                  |     "Path": "/steampipe/"             |
|            |                  | }                                     |
| DevBucket  | AWS::S3::Bucket  | {                                     |
|            |                  |     "BucketName": "TestWebBucket",    |
|            |                  |     "AccessControl": "PublicRead",    |
|            |                  |     "WebsiteConfiguration": {         |
|            |                  |         "IndexDocument": "index.html" |
|            |                  |     }                                 |
|            |                  | }                                     |
| TestVolume | AWS::EC2::Volume | {                                     |
|            |                  |     "Iops": 100,                      |
|            |                  |     "Size": 100,                      |
|            |                  |     "Tags": [                         |
|            |                  |         {                             |
|            |                  |             "Key": "poc",             |
|            |                  |             "Value": "turbot"         |
|            |                  |         }                             |
|            |                  |     ],                                |
|            |                  |     "Encrypted": false,               |
|            |                  |     "VolumeType": "io1",              |
|            |                  |     "AutoEnableIO": false,            |
|            |                  |     "AvailabilityZone": "",           |
|            |                  |     "MultiAttachEnabled": false       |
|            |                  | }                                     |
+------------+------------------+---------------------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/awscfn/tables)**

## Get started

### Install

Download and install the latest AWS CloudFormation plugin:

```bash
steampipe plugin install awscfn
```

### Credentials

No credentials are required.

### Configuration

Installing the latest awscfn plugin will create a config file (`~/.steampipe/config/awscfn.spc`) with a single connection named `awscfn`:

```hcl
connection "awscfn" {
  plugin = "awscfn"

  # Paths is a list of locations to search for CloudFormation template files
  # Paths can be configured with a local directory, a remote Git repository URL, or an S3 bucket URL
  # Wildcard based searches are supported, including recursive searches
  # Local paths are resolved relative to the current working directory (CWD)

  # For example:
  #  - "*.template" matches all CloudFormation template files in the CWD
  #  - "**/*.template" matches all CloudFormation template files in the CWD and all sub-directories
  #  - "../*.template" matches all CloudFormation template files in the CWD's parent directory
  #  - "ELB*.template" matches all CloudFormation template files starting with "ELB" in the CWD
  #  - "/path/to/dir/*.template" matches all CloudFormation template files in a specific directory
  #  - "/path/to/dir/main.template" matches a specific file

  # If paths includes "*", all files (including non-CloudFormation template files) in
  # the CWD will be matched, which may cause errors if incompatible file types exist

  # Defaults to CWD
  paths = ["*.template", "*.yaml", "*.yml", "*.json"]
}
```

### Supported Path Formats

The `paths` config argument is flexible and can search for AWS CloudFormation template files from several different sources, e.g., local directory paths, Git, S3.

The following sources are supported:

- [Local files](#configuring-local-file-paths)
- [Remote Git repositories](#configuring-remote-git-repository-urls)
- [S3](#configuring-s3-urls)

Paths may [include wildcards](https://pkg.go.dev/path/filepath#Match) and support `**` for recursive matching. For example:

```hcl
connection "awscfn" {
  plugin = "awscfn"

  paths = [
    "*.template",
    "~/*.template",
    "github.com/awslabs/aws-cloudformation-templates//*.template",
    "github.com/awslabs/aws-cloudformation-templates//aws/services/ElasticLoadBalancing//*.yaml",
    "gitlab.com/versioncontrol1/cloudformationproject//substacks//*.yaml",
    "s3::https://demo-integrated-2022.s3.ap-southeast-1.amazonaws.com/template_examples//*.yaml"
  ]
}
```

**Note**: If any path matches on `*` without a valid AWS CloudFormation template file extension (i.e. `.template`, `.yaml` etc.), all files (including non-CloudFormation template files) in the directory will be matched, which may cause errors if incompatible file types exist.

#### Configuring Local File Paths

You can define a list of local directory paths to search for AWS CloudFormation template files. Paths are resolved relative to the current working directory. For example:

- `*.template` matches all CloudFormation template files in the CWD.
- `**/*.template` matches all CloudFormation template files in the CWD and all sub-directories.
- `../*.template` matches all CloudFormation template files in the CWD's parent directory.
- `ELB*.template` matches all CloudFormation template files starting with "ELB" in the CWD.
- `/path/to/dir/*.template` matches all CloudFormation template files in a specific directory. For example:
  - `~/*.template` matches all CloudFormation template files in the home directory.
  - `~/**/*.template` matches all CloudFormation template files recursively in the home directory.
- `/path/to/dir/main.template` matches a specific file.

```hcl
connection "awscfn" {
  plugin = "awscfn"

  paths = [ "*.template", "~/*.template", "/path/to/dir/main.template" ]
}
```

**NOTE:** If paths includes `*`, all files (including non-CloudFormation template files) in the CWD will be matched, which may cause errors if incompatible file types exist.

#### Configuring Remote Git Repository URLs

You can also configure `paths` with any Git remote repository URLs, e.g., GitHub, BitBucket, GitLab. The plugin will then attempt to retrieve any AWS CloudFormation template files from the remote repositories.

For example:

- `github.com/awslabs/aws-cloudformation-templates//*.template` matches all top-level CloudFormation template files in the specified github repository.
- `github.com/awslabs/aws-cloudformation-templates//**/*.yaml` matches all CloudFormation template files in the specified github repository and all sub-directories.
- `github.com/awslabs/aws-cloudformation-templates?ref=fix_7677//**/*.template` matches all CloudFormation template files in the specific tag of github repository.
- `github.com/awslabs/aws-cloudformation-templates//aws/services/ElasticLoadBalancing//*.yaml` matches all CloudFormation template files in the specified folder path.

You can specify a subdirectory after a double-slash (`//`) if you want to download only a specific subdirectory from a downloaded directory.

```hcl
connection "awscfn" {
  plugin = "awscfn"

  paths = [ "github.com/awslabs/aws-cloudformation-templates//aws/services/ElasticLoadBalancing//*.yaml" ]
}
```

Similarly, you can define a list of GitLab and BitBucket URLs to search for AWS CloudFormation template files:

```hcl
connection "awscfn" {
  plugin = "awscfn"

  paths = [
    "github.com/awslabs/aws-cloudformation-templates//**/*.template",
    "github.com/awslabs/aws-cloudformation-templates//aws/services/ElasticLoadBalancing//**/*.yaml",
    "gitlab.com/versioncontrol1/cloudformationproject//substacks//*.yaml",
    "gitlab.com/versioncontrol1/cloudformationproject//**/*.yaml"
  ]
}
```

#### Configuring S3 URLs

You can also query all AWS CloudFormation template files stored inside an S3 bucket (public or private) using the bucket URL.

##### Accessing a Private Bucket

In order to access your files in a private S3 bucket, you will need to configure your credentials. You can use your configured AWS profile from local `~/.aws/config`, or pass the credentials using the standard AWS environment variables, e.g., `AWS_PROFILE`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`.

We recommend using AWS profiles for authentication.

**Note:** Make sure that `region` is configured in the config. If not set in the config, `region` will be fetched from the standard environment variable `AWS_REGION`.

You can also authenticate your request by setting the AWS profile and region in `paths`. For example:

```hcl
connection "awscfn" {
  plugin = "awscfn"

  paths = [
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com//*.yaml?aws_profile=<AWS_PROFILE>",
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com/test_folder//*.template?aws_profile=<AWS_PROFILE>"
  ]
}
```

**Note:**

In order to access the bucket, the IAM user or role will require the following IAM permissions:

- `s3:ListBucket`
- `s3:GetObject`
- `s3:GetObjectVersion`

If the bucket is in another AWS account, the bucket policy will need to grant access to your user or role. For example:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ReadBucketObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::123456789012:user/YOUR_USER"
      },
      "Action": ["s3:ListBucket", "s3:GetObject", "s3:GetObjectVersion"],
      "Resource": ["arn:aws:s3:::test-bucket1", "arn:aws:s3:::test-bucket1/*"]
    }
  ]
}
```

##### Accessing a Public Bucket

Public access granted to buckets and objects through ACLs and bucket policies allows any user access to data in the bucket. We do not recommend making S3 buckets public, but if there are specific objects you'd like to make public, please see [How can I grant public read access to some objects in my Amazon S3 bucket?](https://aws.amazon.com/premiumsupport/knowledge-center/read-access-objects-s3-bucket/).

You can query any public S3 bucket directly using the URL without passing credentials. For example:

```hcl
connection "awscfn" {
  plugin = "awscfn"

  paths = [
    "s3::https://bucket-1.s3.us-east-1.amazonaws.com/test_folder//*.yaml",
    "s3::https://bucket-2.s3.us-east-1.amazonaws.com/test_folder//**/*.template"
  ]
}
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-awscfn
- Community: [Slack Channel](https://steampipe.io/community/join)
