variable "aws_region" {
  description = "The region to create the resources into"
  type        = string
  default     = "ap-southeast-2" // Sydney
}

variable "app" {
  type        = string
  description = "The name for the application"
  default     = "byd"
}

variable "stage" {
  description = "The stage or environment for example prod, dev"
  type        = string
  default     = "dev"
}

variable "api_custom_domain_name" {
  type    = string
  default = "api.byd.com"
}

variable "s3_lambda_bucket_name" {
  type    = string
  default = "byd-lambda-zip"
}

variable "db_connection_uri" {
  type    = string
}

# variable "google_client_id" {
#   type = string
#   description = "Google Client ID used for cognito identity provider"
# }

# variable "google_client_secret" {
#   type = string 
#   description = "Google Client Secret used for cognito identity provider"
# }
