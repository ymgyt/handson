variable "vpc_id" {}
variable "subnet_id" {}
variable "name" {}

variable "extra_sgs" {
  default = []
}

variable "environment" {
  default = "dev"
}

variable "extra_packages" {
}

variable "external_nameserver" {}

variable "instance_type" {
  type = "map"

  default = {
    dev  = "t2.micro"
    test = "t2.medium"
    prod = "t2.large"
  }
}

resource "aws_security_group" "allow_http" {
  name        = "${var.name} allow_http"
  description = "Allow HTTP traffic"
  vpc_id      = "${var.vpc_id}"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "app-server" {
  ami                    = "ami-ceafcba8"
  instance_type          = "${lookup(var.instance_type, var.environment)}"
  subnet_id              = "${var.subnet_id}"
  vpc_security_group_ids = ["${concat(var.extra_sgs,aws_security_group.allow_http.*.id)}"]
  user_data              = "${data.template_file.user_data.rendered}"
  key_name = "develop2"
  associate_public_ip_address = true

  tags {
    Name = "${var.name}"
  }

  lifecycle {
  }
}

data "template_file" "user_data" {
  template = "${file("${path.module}/user_data.sh.tpl")}"

  vars {
    packages   = "${var.extra_packages}"
    nameserver = "${var.external_nameserver}"
  }
}

output "hostname" {
  value = "${aws_instance.app-server.private_dns}"
}
