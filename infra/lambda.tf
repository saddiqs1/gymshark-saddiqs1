data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "lambda" {
  name               = "gymshark-saddiqs1-lambda-execution"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/gymshark-saddiqs1"
  retention_in_days = 14
}

data "aws_iam_policy_document" "lambda_logs" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["${aws_cloudwatch_log_group.lambda.arn}:*"]
  }
}

resource "aws_iam_role_policy" "lambda_logs" {
  name   = "cloudwatch-logs"
  role   = aws_iam_role.lambda.id
  policy = data.aws_iam_policy_document.lambda_logs.json
}

data "aws_iam_policy_document" "lambda_pack_sizes" {
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:DeleteItem",
      "dynamodb:PutItem",
      "dynamodb:Scan",
    ]
    resources = [aws_dynamodb_table.pack_sizes.arn]
  }
}

resource "aws_iam_role_policy" "lambda_pack_sizes" {
  name   = "pack-sizes-dynamodb"
  role   = aws_iam_role.lambda.id
  policy = data.aws_iam_policy_document.lambda_pack_sizes.json
}

resource "aws_lambda_function" "app" {
  function_name = "gymshark-saddiqs1"
  description   = "Calculates the pack combination for a customer order"
  role          = aws_iam_role.lambda.arn

  package_type = "Image"
  image_uri    = "${aws_ecr_repository.app.repository_url}:${var.image_tag}"
  architectures = [
    "arm64",
  ]

  memory_size = 128
  timeout     = 10

  environment {
    variables = {
      APP_ENV                      = "production"
      AWS_LWA_READINESS_CHECK_PATH = "/health"
      LOG_LEVEL                    = "info"
      PACK_SIZES_TABLE_NAME        = aws_dynamodb_table.pack_sizes.name
    }
  }

  depends_on = [
    aws_cloudwatch_log_group.lambda,
    aws_ecr_repository_policy.lambda,
    aws_iam_role_policy.lambda_logs,
    aws_iam_role_policy.lambda_pack_sizes,
  ]
}
