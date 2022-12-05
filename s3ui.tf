resource "aws_s3_bucket" "s3ui" {
  bucket = "distribution-portal-ui"

  tags = {
    Name        = "Distribution Portal UI"
    Environment = "Dev"
  }

  website {
    index_document = "index.html"
    error_document = "404.html"
  }
}

resource "aws_s3_bucket_acl" "s3ui" {
  bucket = aws_s3_bucket.s3ui.id
  acl    = "public-read"
}

resource "aws_s3_bucket_policy" "s3ui" {
  bucket = aws_s3_bucket.s3ui.id
  policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                "Sid": "PublicReadGetObject",
                "Effect": "Allow",
                "Principal": "*",
                "Action": "s3:GetObject",
                "Resource": "arn:aws:s3:::distribution-portal-ui/*"
            }            
        ]
    })
}