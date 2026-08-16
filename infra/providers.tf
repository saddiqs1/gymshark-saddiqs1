terraform {
  required_version = ">= 1.10.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }

  backend "s3" {}
}

provider "aws" {
  region = "eu-west-1"

  default_tags {
    tags = {
      Project   = "gymshark-saddiqs1"
      ManagedBy = "Terraform"
    }
  }
}
