# Simple Task Tracker - Go

_Small app for learning Golang, idea was taken from the Kubernetes/Docker repo, but done in a more code involved way._

Application integrates with AWS DynamoDB for reading, updating, and deleting tasks.

The old application stored tasks in JSON format to the file provided, this may have been simpler to do with a map data-type but provided
a great excerise in handling JSON and using structs. The old code has been placed into `old_main.go` for learning and documentation purposes. The `.exe` can also be added to your environment variables by referencing the `.bat` file, this means you can invoke the app to run from the command line using `%ttgo%`, this is what I chose, or whatever you choose to name the variable. Although the location of the `.exe` must be changed first _(this is no longer required with the updated app since it uses AWS for storage.)_.

