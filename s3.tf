resource "aws_s3_bucket" "s3_instance_01" {
  bucket = var.s3_bucket_name

  tags = {
    Name        = "Distribution Portal"
    Environment = "Dev"
  }
}

resource "aws_s3_bucket_acl" "s3_instance_01_acl" {
  bucket = aws_s3_bucket.s3_instance_01.id
  acl    = "private"
}