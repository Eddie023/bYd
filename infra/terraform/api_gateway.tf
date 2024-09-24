resource "aws_api_gateway_rest_api" "main" {
  name        = "${local.name}-rest-api"
  description = "Rest API endpoint for ${local.name}"
}

resource "aws_api_gateway_resource" "main" {
  rest_api_id = aws_api_gateway_rest_api.main.id
  parent_id   = aws_api_gateway_rest_api.main.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_authorizer" "cognito" {
  name          = "${local.name}-authorizer"
  type          = "COGNITO_USER_POOLS"
  rest_api_id   = aws_api_gateway_rest_api.main.id
  provider_arns = [aws_cognito_user_pool.main.arn]
}

resource "aws_api_gateway_method" "main" {
  rest_api_id = aws_api_gateway_rest_api.main.id
  resource_id = aws_api_gateway_resource.main.id
  http_method = "ANY"

  authorization = "COGNITO_USER_POOLS"
  authorizer_id = aws_api_gateway_authorizer.cognito.id
}

resource "aws_api_gateway_integration" "main" {
  rest_api_id             = aws_api_gateway_rest_api.main.id
  resource_id             = aws_api_gateway_resource.main.id
  http_method             = aws_api_gateway_method.main.http_method
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.restapi_handler.invoke_arn
  integration_http_method = "POST"
}

resource "aws_api_gateway_deployment" "main" {
  rest_api_id = aws_api_gateway_rest_api.main.id
  stage_name  = var.stage

  depends_on = [aws_api_gateway_integration.main]
}

resource "aws_cloudwatch_log_group" "gateway_log" {
  name              = "API-Gateway-Execution-Logs_${aws_api_gateway_rest_api.main.id}/${var.stage}"
  retention_in_days = 7
}
