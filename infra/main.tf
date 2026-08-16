resource "aws_ecr_repository" "app" {
  name                 = "gymshark-saddiqs1"
  image_tag_mutability = "IMMUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }
}

resource "aws_ecr_lifecycle_policy" "app" {
  repository = aws_ecr_repository.app.name

  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "Keep the latest 20 images"
      selection = {
        tagStatus   = "any"
        countType   = "imageCountMoreThan"
        countNumber = 20
      }
      action = {
        type = "expire"
      }
    }]
  })
}

data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "lambda_ecr" {
  statement {
    sid    = "LambdaECRImageRetrievalPolicy"
    effect = "Allow"
    actions = [
      "ecr:BatchGetImage",
      "ecr:GetDownloadUrlForLayer",
    ]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    resources = [aws_ecr_repository.app.arn]

    condition {
      test     = "ArnLike"
      variable = "aws:SourceArn"
      values   = ["arn:aws:lambda:eu-west-1:${data.aws_caller_identity.current.account_id}:function:gymshark-saddiqs1"]
    }
  }
}

resource "aws_ecr_repository_policy" "lambda" {
  repository = aws_ecr_repository.app.name
  policy     = data.aws_iam_policy_document.lambda_ecr.json
}

output "ecr_repository_url" {
  description = "URL used to tag and push application images"
  value       = aws_ecr_repository.app.repository_url
}
