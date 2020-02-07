# Specify the provider and access details
provider "aws" {
  region = "${var.aws_region}"
}

# Create a VPC to launch our instances into
resource "aws_vpc" "tasktrackVPC" {
  cidr_block = "10.0.0.0/16"
}

# Create an internet gateway to give our subnet access to the outside world
resource "aws_internet_gateway" "tasktrackIGW" {
  vpc_id = "${aws_vpc.tasktrackVPC.id}"
}

# Grant the VPC internet access on its main route table
resource "aws_route" "internet_access" {
  route_table_id         = "${aws_vpc.tasktrackVPC.main_route_table_id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.tasktrackIGW.id}"
}

# Create a subnet to launch our instances into
resource "aws_subnet" "tasktrackSubnet" {
  vpc_id                  = "${aws_vpc.tasktrackVPC.id}"
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
}

# A security group for the ELB so it is accessible via the web
resource "aws_security_group" "elb" {
  name        = "terraform_example_elb"
  description = "Used in the terraform"
  vpc_id      = "${aws_vpc.tasktrackVPC.id}"

  # HTTP access from anywhere
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # outbound internet access
  egress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Our default security group to access
# the instances over SSH and HTTP
resource "aws_security_group" "tasktrackSG" {
  name        = "task_tracker"
  description = "Used in terraform"
  vpc_id      = "${aws_vpc.tasktrackVPC.id}"

  # SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP access from the VPC
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  # outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_elb" "tasktrackerELB" {
  name = "task-tracker-elb"

  subnets         = ["${aws_subnet.tasktrackSubnet.id}"]
  security_groups = ["${aws_security_group.elb.id}"]
  instances       = ["${aws_instance.tasktracker.id}", "${aws_instance.tasktrackerTwo.id}"]

  listener {
    instance_port     = 8080
    instance_protocol = "http"
    lb_port           = 8080
    lb_protocol       = "http"
  }
}


resource "aws_instance" "tasktracker" {
  # The connection block tells our provisioner how to
  # communicate with the resource (instance)
  connection {
    # The default username for our AMI
    user = "ec2-user"
    host = "${self.public_ip}"
    # The connection will use the local SSH agent for authentication.
  }

  instance_type = "t2.micro"

  # Lookup the correct AMI based on the region
  # we specified
  ami = "${lookup(var.aws_amis, var.aws_region)}"

  key_name = "BusinessInfra"

  # Our Security group to allow HTTP and SSH access
  vpc_security_group_ids = ["${aws_security_group.tasktrackSG.id}"]

  # We're going to launch into the same subnet as our ELB. In a production
  # environment it's more common to have a separate private subnet for
  # backend instances.
  subnet_id = "${aws_subnet.tasktrackSubnet.id}"
}


resource "aws_instance" "tasktrackerTwo" {
  # The connection block tells our provisioner how to
  # communicate with the resource (instance)
  connection {
    # The default username for our AMI
    user = "ec2-user"
    host = "${self.public_ip}"

    # The connection will use the local SSH agent for authentication.
  }

  instance_type = "t2.micro"

  # Lookup the correct AMI based on the region
  # we specified
  ami = "${lookup(var.aws_amis, var.aws_region)}"

  key_name = "BusinessInfra"

  # Our Security group to allow HTTP and SSH access
  vpc_security_group_ids = ["${aws_security_group.tasktrackSG.id}"]

  # We're going to launch into the same subnet as our ELB. In a production
  # environment it's more common to have a separate private subnet for
  # backend instances.
  subnet_id = "${aws_subnet.tasktrackSubnet.id}"

            
}

output "ELB_IP" {
  value = "${aws_elb.tasktrackerELB.dns_name}"
}