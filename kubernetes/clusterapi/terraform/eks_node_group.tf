data "aws_ssm_parameter" "eks_ami_release_version" {
  name = "/aws/service/eks/optimized-ami/${aws_eks_cluster.handson.version}/amazon-linux-2/recommended/release_version"
}

# Node Group
resource "aws_eks_node_group" "group1" {
  cluster_name  = aws_eks_cluster.handson.name
  node_role_arn = aws_iam_role.eks_node_group.arn
  scaling_config {
    desired_size = 2
    max_size     = 2
    min_size     = 1
  }
  subnet_ids = slice(module.vpc.private_subnets, 0, 3)
  version    = aws_eks_cluster.handson.version

  # https://docs.aws.amazon.com/eks/latest/APIReference/API_Nodegroup.html#AmazonEKS-Type-Nodegroup-amiType
  ami_type        = "AL2_x86_64"
  release_version = data.aws_ssm_parameter.eks_ami_release_version.value
  instance_types  = ["t3.medium"]

  depends_on = [
    aws_iam_role_policy_attachment.eks_node_group
  ]

  lifecycle {
    ignore_changes = [
      scaling_config[0].desired_size
    ]
  }
}

# IAM
data "aws_iam_policy_document" "eks_node_group_assume_role_policy" {
  statement {
    sid     = "EKSNodeGroupAssumeRole"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.${local.aws_dns}"]
    }
  }
}

resource "aws_iam_role" "eks_node_group" {
  name               = "${local.cluster_name}-nodegroup"
  assume_role_policy = data.aws_iam_policy_document.eks_node_group_assume_role_policy.json
}

resource "aws_iam_role_policy_attachment" "eks_node_group" {
  # https://docs.aws.amazon.com/eks/latest/userguide/create-node-role.html
  for_each = toset([
    "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
    "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
  ])

  policy_arn = each.value
  role       = aws_iam_role.eks_node_group.name
}

