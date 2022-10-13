connection "awscfn" {
  plugin = "awscfn"

  # Paths is a list of locations to search for CloudFormation template files
  # All paths are resolved relative to the current working directory (CWD)
  # Wildcard based searches are supported, including recursive searches

  # Defaults to CWD
  paths = ["*.template", "*.yaml", "*.yml", "*.json"]
}