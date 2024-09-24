resource "aws_cognito_user_pool" "main" {
  name                     = "${local.name}-user-pool"
  username_attributes      = ["email"]
  auto_verified_attributes = ["email"]

  schema {
    name                     = "given_name"
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    required                 = true

    string_attribute_constraints {
      min_length = 1
      max_length = 50
    }
  }

  schema {
    name                     = "family_name"
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    required                 = true

    string_attribute_constraints {
      min_length = 1
      max_length = 50
    }
  }

  lambda_config {
    post_confirmation = aws_lambda_function.cognito_post_signup.arn
  }

  depends_on = [aws_lambda_function.cognito_post_signup]
}

# Uncomment the following section if you want to include google as your identity provider for cognito. 
# 
# resource "aws_cognito_identity_provider" "google" {
#   user_pool_id = aws_cognito_user_pool.main.id 
#   provider_name = "Google"
#   provider_type = "Google"

#   provider_details = {
#     client_id = var.google_client_id
#     client_secret = var.google_client_secret
#     authorize_scopes = "profile email openid"
#   }

#   attribute_mapping = {
#     email = "email"
#     username = "sub"
#     given_name = "given_name"
#     family_name = "family_name"
#     email_verified = "email_verified"
#   }

#   lifecycle {
#     ignore_changes = [ user_pool_id, provider_name ]
#   }
# }

resource "aws_cognito_resource_server" "main" {
  name         = "${local.name}-resource-server"
  identifier   = "https://${var.api_custom_domain_name}"
  user_pool_id = aws_cognito_user_pool.main.id

  scope {
    scope_name        = "all"
    scope_description = "Get access to all API Gateway endpoints"
  }
}

resource "aws_cognito_user_pool_domain" "main" {
  domain       = local.name
  user_pool_id = aws_cognito_user_pool.main.id
}

resource "aws_cognito_user_pool_client" "main" {
  name                                 = "${local.name}-user-pool-client"
  user_pool_id                         = aws_cognito_user_pool.main.id
  generate_secret                      = false
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code", "implicit"]
  supported_identity_providers         = ["COGNITO"]
  # NOTE: uncomment the following if you are using Google as your identity provider
  # supported_identity_providers         = ["COGNITO", "Google"]
  allowed_oauth_scopes                 = ["email", "openid"]
  // ADD your Production URL callback and logout URLs 
  callback_urls                        = ["http://localhost:3000/login"] 
  logout_urls = [ "http://localhost:3000/login?signout=true" ]

  depends_on = [aws_cognito_user_pool.main, aws_cognito_resource_server.main]
}

