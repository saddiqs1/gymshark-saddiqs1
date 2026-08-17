resource "aws_dynamodb_table" "pack_sizes" {
  name         = "gymshark-saddiqs1-pack-sizes"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "size"

  attribute {
    name = "size"
    type = "N"
  }

  point_in_time_recovery {
    enabled = true
  }

  server_side_encryption {
    enabled = true
  }
}

output "pack_sizes_table_name" {
  description = "DynamoDB table containing the configured pack sizes"
  value       = aws_dynamodb_table.pack_sizes.name
}

output "pack_sizes_table_arn" {
  description = "ARN of the DynamoDB pack sizes table"
  value       = aws_dynamodb_table.pack_sizes.arn
}
