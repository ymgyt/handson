terraform {
  required_version = ">= 1.6.4"

  required_providers {
    aws = {
      source  = "registry.terraform.io/hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  # Tags to apply to all AWS resources by default
  default_tags {
    tags = {
      Project = "twnana"
    }
  }
}
