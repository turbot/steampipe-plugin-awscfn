---
title: "Steampipe Table: awscfn_output - Query AWS CloudFormation Outputs using SQL"
description: "Allows users to query Outputs from AWS CloudFormation Stacks, providing information about the outputs of each stack."
---

# Table: awscfn_output - Query AWS CloudFormation Outputs using SQL

AWS CloudFormation is a service that helps you model and set up your Amazon Web Services resources so you can spend less time managing those resources and more time focusing on your applications that run in AWS. Outputs in AWS CloudFormation provide a way to output values from a stack and make them easily accessible. They can be used to import and export values between different stacks, and can be used to manage and organize resources in your AWS environment.

## Table Usage Guide

The `awscfn_output` table provides insights into the outputs of AWS CloudFormation Stacks. As a DevOps engineer or Cloud Architect, you can explore output-specific details through this table, including stack names, output keys, and output values. This can be particularly useful for managing and organizing your AWS resources, as well as for troubleshooting and optimizing your AWS environment.

## Examples

### Basic info
Explore the key details of your AWS CloudFormation stack outputs. This query allows you to gain insights into the output values and their associated descriptions and paths, which can be beneficial in understanding your stack's configuration and performance.

```sql
select
  name,
  description,
  value,
  path
from
  awscfn_output;
```

### List outputs that return an EC2 instance public DNS name
Explore which CloudFormation outputs provide public DNS names for EC2 instances. This can be useful for identifying resources that are potentially exposed to the internet.

```sql
select
  name,
  value,
  description,
  path
from
  awscfn_output
where
  value like '%Fn::GetAtt:%PublicDnsName%';
```

### List outputs that show sensitive parameter values
Identify the areas in your AWS CloudFormation outputs that may be exposing sensitive parameter values. This can be useful in enhancing security by pinpointing potential areas of data leakage.

```sql
with output_table as (
  select
    name,
    description,
    split_part(substring(value from '\w*Ref:*\w*'), ':', 2) as parameter_reference,
    path
  from
    awscfn_output
  where
    value like '%Ref:%'
)
select
  o.name,
  o.description,
  o.path
from
  output_table as o
  left join awscfn_parameter as p on p.name = o.parameter_reference and o.path = p.path
where
  p.no_echo;
```