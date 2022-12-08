resource "aws_dynamodb_table" "distribution-portal-table" {
  name = "distribution-portal"
  billing_mode = "PROVISIONED"
  read_capacity = "30"
  write_capacity = "30"
  attribute {
    name = "name"
    type = "S"
  }

  hash_key = "name"
}

resource "aws_dynamodb_table" "installation-service-table" {
  name = "installation-service"
  billing_mode = "PROVISIONED"
  read_capacity = "30"
  write_capacity = "30"
  
  attribute {
    name = "org-project-workspace"
    type = "S"
  }

  attribute {
      name = "pkg-name"
      type = "S"
  }

  attribute {
    name = "pkg-name-version"
    type = "S"
  }

  attribute {
    name = "pkg-name-version-status"
    type = "S"
  }

  attribute {
    name = "id"
    type = "S"
  }

  global_secondary_index {
    name = "id-index"
    hash_key          = "id"
    projection_type    = "ALL"
    write_capacity     = 10
    read_capacity      = 10
    non_key_attributes = []
  }

  local_secondary_index {
    name               = "pkg-name-version-index"
    range_key          = "pkg-name-version"
    projection_type    = "ALL"
    non_key_attributes = []
  }

  local_secondary_index {
    name               = "pkg-name-index"
    range_key          = "pkg-name"
    projection_type    = "ALL"
    non_key_attributes = []
  }

  range_key = "pkg-name-version-status"
  hash_key = "org-project-workspace"
}