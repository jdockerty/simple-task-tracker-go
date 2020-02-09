variable "aws_region" {
  default = "eu-west-2"
}

variable "aws_amis" {
  default = {
      eu-west-2 = "ami-09120dc9a0ea0cdea" # My AMI with Go + AWS CLI configured
  }
}