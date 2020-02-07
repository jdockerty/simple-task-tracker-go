variable "aws_region" {
  default = "eu-west-2"
}

variable "aws_amis" {
  default = {
      eu-west-2 = "ami-003ff27f7edfb65e9" # My AMI with Go + dependencies + AWS CLI configured
  }
}