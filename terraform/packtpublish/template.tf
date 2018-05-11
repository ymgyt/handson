variable "region" {
  description = "AWS region"
  default     = "ap-northeast-1"
}

variable "environment" {
  default = "dev"
}

variable "external_nameserver" {
  default = "8.8.8.8"
}

variable "extra_packages" {
  default = "wget bind-utils"
}

provider "aws" {
  region = "${var.region}"
}

# resource configuration

resource "aws_vpc" "my_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_internet_gateway" "this" {
  vpc_id = "${aws_vpc.my_vpc.id}"
}

resource "aws_route_table" "this" {
  vpc_id = "${aws_vpc.my_vpc.id}"
}

resource "aws_route" "this" {
  route_table_id = "${aws_route_table.this.id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id = "${aws_internet_gateway.this.id}"
}

resource "aws_subnet" "public" {
  vpc_id     = "${aws_vpc.my_vpc.id}"
  cidr_block = "10.0.1.0/24"
}

resource "aws_security_group" "default" {
  name        = "Default SG"
  description = "Allow SSH access"
  vpc_id      = "${aws_vpc.my_vpc.id}"

  ingress {
    from_port  = 22
    to_port    = 22
    protocol   = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

module "mighty_trousers" {
  source              = "./modules/application"
  vpc_id              = "${aws_vpc.my_vpc.id}"
  subnet_id           = "${aws_subnet.public.id}"
  name                = "MightyTrousers"
  environment         = "${var.environment}"
  extra_sgs           = ["${aws_security_group.default.id}"]
  extra_packages      = "${var.extra_packages}"
  external_nameserver = "${var.external_nameserver}"
}
