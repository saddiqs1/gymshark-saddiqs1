variable "image_tag" {
  description = "Immutable ECR image tag deployed to Lambda"
  type        = string

  validation {
    condition     = length(trimspace(var.image_tag)) > 0
    error_message = "image_tag must not be empty."
  }
}
