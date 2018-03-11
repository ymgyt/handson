provider "aws" {
  region = "ap-northeast-1"
}

# resource configuration
resource "aws_instance" "hello-instance" {
  ami           = "ami-ceafcba8"
  instance_type = "t2.micro"

  tags {
    Name = "hello-update-instance"
  }
}
