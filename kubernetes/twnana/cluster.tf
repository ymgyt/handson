data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "master" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = "t3.medium"
  subnet_id                   = aws_subnet.public["ap-northeast-1a"].id
  key_name                    = aws_key_pair.ssh.key_name
  associate_public_ip_address = true
  vpc_security_group_ids      = [aws_security_group.ssh.id]

  tags = {
    Name = "master"
  }
}

resource "aws_instance" "worker" {
  for_each = {
    "worker-1" = {
      az = "ap-northeast-1c"
    },
    "worker-2" = {
      az = "ap-northeast-1d"
    }
  }

  ami                         = data.aws_ami.ubuntu.id
  instance_type               = "t3.medium"
  subnet_id                   = aws_subnet.public["${each.value.az}"].id
  key_name                    = aws_key_pair.ssh.key_name
  associate_public_ip_address = true
  vpc_security_group_ids      = [aws_security_group.ssh.id]

  tags = {
    Name = "${each.key}"
  }
}

resource "aws_key_pair" "ssh" {
  key_name   = "twnana"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJBBWTQ7KgWoz4V+nrDKf/PLg53VwpQzba/FUhbnyWqR ymgyt"
}

resource "aws_security_group" "ssh" {
  name        = "allow_ssh"
  description = "Allow ssh connection"
  vpc_id      = aws_vpc.kubernetes.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "allow_ssh"
  }
}

