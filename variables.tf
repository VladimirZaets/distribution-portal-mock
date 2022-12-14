variable "s3_bucket_name" {
  type = string
  description = "The unique name for S3 bucket"
}

variable "aws_settings" {
  type = object({
    region    = string
    access_key = string
    secret_key = string
    token = string
  })
  default = {
    region = "us-east-1"
  }
  sensitive = true
}