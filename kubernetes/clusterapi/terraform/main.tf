terraform {
  required_version = "~> 1.5.3"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"

  profile = var.profile

  # Tags to apply to all AWS resources by default
  default_tags {
    tags = {
      ManagedBy = "terraform"
    }
  }
}

locals {
  cluster_name = "clusterapi-handson"
}

data "aws_availability_zones" "available" {
  # exclude local zones
  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.0.0"

  name = "clusterapi-handson"
  cidr = "10.0.0.0/16"
  azs  = slice(data.aws_availability_zones.available.names, 0, 3)

  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.4.0/24", "10.0.5.0/24", "10.0.6.0/24"]

  enable_nat_gateway = true
  single_nat_gateway = true
  # Required in eks vpc requirements
  enable_dns_hostnames = true
}

resource "aws_eks_cluster" "handson" {
  name     = local.cluster_name
  role_arn = aws_iam_role.eks_cluster.arn
  vpc_config {
    endpoint_private_access = false
    endpoint_public_access  = true
    public_access_cidrs     = ["0.0.0.0/0"]
    subnet_ids              = slice(module.vpc.private_subnets, 0, 3)
  }
  version = "1.28"

  #  # Ensure that IAM Role permissions are created before and deleted after EKS Cluster handling.
  # Otherwise, EKS will not be able to properly delete EKS managed EC2 infrastructure such as Security Groups.
  depends_on = [
    aws_iam_role_policy_attachment.eks_cluster
  ]
}

# IAM Role
locals {
  aws_dns = data.aws_partition.current.dns_suffix
}

# Assume policy for eks cluster
data "aws_iam_policy_document" "eks_assume_role_policy" {
  statement {
    sid     = "EKSClusterAssumeRole"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["eks.${data.aws_partition.current.dns_suffix}"]
    }
  }
}

resource "aws_iam_role" "eks_cluster" {
  name               = local.cluster_name
  assume_role_policy = data.aws_iam_policy_document.eks_assume_role_policy.json
}

# Attach managed EKS Policy to eks role
resource "aws_iam_role_policy_attachment" "eks_cluster" {
  for_each = toset(["arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"])

  policy_arn = each.value
  role       = aws_iam_role.eks_cluster.name
}

data "aws_partition" "current" {}
