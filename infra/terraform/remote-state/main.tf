variable "stage" {
    type = string
    default = "dev"
}

variable "app_name" {
  type = string
  default = "byd"
}

provider "aws" {
    region = "ap-southeast-2"
}

resource "aws_s3_bucket" "terraform_state" {
    bucket = "${var.app_name}-${var.stage}-tfstate"

    lifecycle {
      prevent_destroy = true
    }
}

resource "aws_s3_bucket_versioning" "this" {
    bucket = aws_s3_bucket.terraform_state.id

    versioning_configuration {
      status = "Enabled"
    }
}

resource "aws_dynamodb_table" "terraform_state_lock" {
    name = "${var.app_name}-${var.stage}-tfstate"
    read_capacity = 1
    write_capacity = 1
    hash_key = "LockID"

    attribute {
      name = "LockID"
      type = "S"
    }
}

output "remote_state_s3_bucket_name" {
  value = aws_s3_bucket.terraform_state.bucket
}

output "remote_state_ddb_table_name" {
  value = aws_dynamodb_table.terraform_state_lock.name
}