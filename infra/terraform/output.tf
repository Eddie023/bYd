output "api_gateway_endpoint" {
  value = aws_api_gateway_deployment.main.invoke_url
}

output "restapi_lambda_handler" {
  value = aws_lambda_function.restapi_handler.function_name
}

output "cognito_userpool_id" {
  value = aws_cognito_user_pool.main.id
}

output "cognito_userpool_client_id" {
  value = aws_cognito_user_pool_client.main.id
}

output "restapi_cloudwatch_loggroup" {
  value = aws_cloudwatch_log_group.restapi_handler_log_group.name
}

output "restapi_endpoint" {
  value = aws_route53_record.a_record.name
}

output "cognito_domain" {
  value = "https://${aws_cognito_user_pool_domain.main.domain}.auth.${var.aws_region}.amazoncognito.com"
}