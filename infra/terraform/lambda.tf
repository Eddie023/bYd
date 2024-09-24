data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

// defines an IAM role that allows Lambda to access resources in your AWS account.
resource "aws_iam_role" "iam_for_lambda" {
  name               = "byd_iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

// Attaches AWS managed policy that allows your lambda function to write to CloudWatch logs
resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "archive_file" "lambda_zip" {
  type = "zip"

  source_file = "../../bin/rest-api/bootstrap"
  output_path = "rest-api-lambda-handler"
}

resource "aws_lambda_function" "restapi_handler" {
  function_name = "${local.name}-restapi-lambda-handler"
  filename      = "rest-api-lambda-handler"

  role    = aws_iam_role.iam_for_lambda.arn
  handler = "bootstrap"

  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  timeout = 60
  runtime = "provided.al2"

  lifecycle {
    ignore_changes = [source_code_hash]
  }

  environment {
    variables = {
      DB_CONNECTION_URI = var.db_connection_uri
    }
  }
}

data "archive_file" "post_signup_lambda_zip" {
  type = "zip"

  source_file = "../../bin/post-signup/bootstrap"
  output_path = "cognito-post-signup-handler"
}


resource "aws_lambda_function" "cognito_post_signup" {
  function_name = "${local.name}-cognito_post_signup-handler"
  filename      = "cognito-post-signup-handler"

  role    = aws_iam_role.iam_for_lambda.arn
  handler = "bootstrap"

  source_code_hash = data.archive_file.post_signup_lambda_zip.output_base64sha256

  runtime = "provided.al2"

  timeout = 60

  environment {
    variables = {
      DB_CONNECTION_URI = var.db_connection_uri
    }
  }
}

resource "aws_lambda_permission" "cognito_lambda" {
  statement_id  = "AllowCognitoInvokeLambda"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.cognito_post_signup.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.main.arn
}


resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.restapi_handler.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.main.execution_arn}/*"
}

resource "aws_cloudwatch_log_group" "restapi_handler_log_group" {
  name              = "${local.name}-restapi-handler-logs"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "post_signup_log_group" {
  name              = "${local.name}-post-signup-logs"
  retention_in_days = 7
}
