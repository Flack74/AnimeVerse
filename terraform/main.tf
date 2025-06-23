terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

# Security Group
resource "aws_security_group" "animeverse_sg" {
  name_prefix = "animeverse-"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8000
    to_port     = 8000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "animeverse-sg"
  }
}

# EC2 Instance
resource "aws_instance" "animeverse_server" {
  ami           = "ami-0c02fb55956c7d316" # Amazon Linux 2023
  instance_type = "t2.micro"

  vpc_security_group_ids = [aws_security_group.animeverse_sg.id]

  user_data = <<-EOF
    #!/bin/bash
    yum update -y
    yum install -y docker
    systemctl start docker
    systemctl enable docker
    usermod -a -G docker ec2-user
    
    # Pull and run the application
    docker pull flack74621/animeverse:latest
    docker run -d -p 8000:8000 --name animeverse-app flack74621/animeverse:latest
  EOF

  tags = {
    Name = "animeverse-server"
  }
}

# Outputs
output "ec2_public_ip" {
  value = aws_instance.animeverse_server.public_ip
}

output "application_url" {
  value = "http://${aws_instance.animeverse_server.public_ip}:8000"
}