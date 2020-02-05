# Task Tracker - Go

_Small app for learning Golang, idea was taken from the Kubernetes/Docker repo, but done in a more code involved way._

Application integrates with AWS DynamoDB for reading, updating, and deleting tasks. The core aspect of the application comes from the `html/template` and `net/http` in-built packages for Go, providing dynamic HTML pages and HTTP routing respectively.

The previous iterations of this app were made for Go learning purposes, these are placed into the `legacycode` folder. The first instance of the application was primarily writing to a JSON file, the second was an interactive CLI that integrated with AWS.

![mainmenu](https://github.com/jdockerty/simpletasktrackergo/blob/master/images/menu.png)
