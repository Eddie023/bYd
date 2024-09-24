locals {
  region = var.aws_region
  name   = "${var.app}-${var.stage}"
}