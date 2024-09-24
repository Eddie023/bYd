resource "aws_s3_bucket" "lambda_bucket" {
  bucket = var.s3_lambda_bucket_name

  force_destroy = true
}

resource "aws_s3_bucket_ownership_controls" "lambda_bucket" {
  bucket = aws_s3_bucket.lambda_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "lambda_bucket_acl" {
  depends_on = [aws_s3_bucket_ownership_controls.lambda_bucket]
  bucket     = aws_s3_bucket.lambda_bucket.id

  acl = "private"
}
