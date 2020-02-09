# Task Tracker - Go

_Small app for learning Golang, idea was taken from the Kubernetes/Docker repo, but done in a more code involved way._

Application integrates with AWS DynamoDB for adding, viewing, updating, and deleting tasks. The core aspect of the application comes from the `html/template` and `net/http` in-built packages for Go, providing dynamic HTML pages and HTTP routing respectively.

The previous iterations of this app were made for Go learning purposes, these are placed into the `legacycode` folder. The first instance of the application was primarily writing to a JSON file, the second was an interactive CLI that integrated with AWS.

A simple menu provides basic navigation to various pages.

![mainmenu](https://github.com/jdockerty/simpletasktrackergo/blob/master/images/menu.png)


Tasks are viewed in a table which is dynmically generated upon visiting the page, the data is pulled from the AWS DynamoDB table after clicking the `View Tasks` menu option, this calls the relevant Go function.

![viewtasks](https://github.com/jdockerty/simpletasktrackergo/blob/master/images/viewtasks.png)


Terraform was used to deploy the web application onto a custom AMI, although this was not particularly necessary once utilising `user data` in AWS, as Golang and Git could be installed on the machine at boot-time. The Terraform file creates a VPC for deploying the instances, with an Elastic Load Balancer that listens on port 8080, this forwards the incoming traffic to the main subnet, where the instances are serving the web application. The load balancer is the internet facing portion of the architecture; the instances are not directly accessible via HTTP, they only allow traffic from `10.0.0.0/16`, which is the CIDR block allocated to the VPC. SSH is enabled from anywhere for testing purposes, although it is no longer necessary as `user data` provides the relevant configuration to the instances.
