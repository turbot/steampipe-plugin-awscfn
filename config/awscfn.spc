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
