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