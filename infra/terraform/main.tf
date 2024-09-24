terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }

  backend "s3" {
    // app_name+stage-tfstate
    bucket         = "byd-dev-tfstate"
    dynamodb_table = "byd-dev-tfstate"
    key            = "terraform.tfstate"
    region         = "ap-southeast-2"
  }

  required_version = ">= 1.3.7"
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Environment = var.stage
      Name        = var.app
    }
  }
}
