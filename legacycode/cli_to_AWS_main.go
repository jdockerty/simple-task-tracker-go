package legacycode

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"os"
	"strings"
	"runtime"
)

const (
	currentRuntime string = runtime.GOOS
)

// awsConnection is a helper function for creating an AWS Session and returns the DynamoDB client for use around the program.
func awsConnection() *dynamodb.DynamoDB {
	session, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-2")})
	if err != nil {
		panic(err)
	}

	dbInstance := dynamodb.New(session)
	return dbInstance
}

// readDynamoTable will print out the contents of the entire DynamoDB table to the console window.
func readDynamoTable(dbSession *dynamodb.DynamoDB) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("Task-Tracker"),
	}

	allData, err := dbSession.Scan(input)
	if err != nil {
		panic(err)
	}

	if len(allData.Items) == 0 {
		fmt.Println("Table is empty.")
	} else {
		for _, value := range allData.Items {
			fmt.Printf("\nTaskID: %s\nTask Name: %s\nTask Details: %s\nCompletion Date: %s\n\n",
				*value["TaskID"].S, *value["Task Name"].S, *value["Task Details"].S, *value["Completion Date"].S)
		}
	}

}

func modifyItemDynamoTable(dbSession *dynamodb.DynamoDB) {
	fmt.Print("Enter the ID of the task you wish to modify: ")
	modifyTaskID := readUserInput()

	fmt.Print("Enter a task name: ")
	taskName := readUserInput()

	fmt.Print("Enter the task details: ")
	taskDetails := readUserInput()

	fmt.Print("Enter the completion date: ")
	completeBy := readUserInput()

	itemInput := &dynamodb.PutItemInput{
		TableName: aws.String("Task-Tracker"),
		Item: map[string]*dynamodb.AttributeValue{
			"TaskID": {
				S: aws.String(modifyTaskID),
			},
			"Task Name": {
				S: aws.String(taskName),
			},
			"Task Details": {
				S: aws.String(taskDetails),
			},
			"Completion Date": {
				S: aws.String(completeBy),
			},
		},
	}

	_, err := dbSession.PutItem(itemInput)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nTask modified: \n\tTask Name = %s\n\tTask Details = %s\n\tCompletion Date = %s\n\n", taskName, taskDetails, completeBy)
}

// addItemDynamoTable will add the item to DynamoDB and generate a TaskID with a new UUID. This is displayed to the user for clarity.
func addItemDynamoTable(dbSession *dynamodb.DynamoDB) {
	newTaskID := uuid.New()

	fmt.Print("Enter a task name: ")
	taskName := readUserInput()

	fmt.Print("Enter the task details: ")
	taskDetails := readUserInput()

	fmt.Print("Enter the completion date: ")
	completeBy := readUserInput()

	itemInput := &dynamodb.PutItemInput{
		TableName: aws.String("Task-Tracker"),
		Item: map[string]*dynamodb.AttributeValue{
			"TaskID": {
				S: aws.String(newTaskID.String()),
			},
			"Task Name": {
				S: aws.String(taskName),
			},
			"Task Details": {
				S: aws.String(taskDetails),
			},
			"Completion Date": {
				S: aws.String(completeBy),
			},
		},
	}

	_, err := dbSession.PutItem(itemInput)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nTask sent: \n\tTaskID = %s\n\tTask Name = %s\n\tTask Details = %s\n\tCompletion Date = %s\n\n",
		newTaskID.String(), taskName, taskDetails, completeBy)
}

// Deletes the task with the relevant TaskID from the DynamoDB table, this function is best used alongside the readDynamoTable function as it provides
// the relevant TaskID to delete.
func deleteItemDynamoTable(dbSession *dynamodb.DynamoDB) {
	fmt.Print("Enter the TaskID to delete: ")
	itemChoice := readUserInput()

	itemDelete := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"TaskID": {
				S: aws.String(itemChoice),
			},
		},
		TableName: aws.String("Task-Tracker"),
	}

	_, err := dbSession.DeleteItem(itemDelete)
	if err != nil {
		panic(err)
	}
	fmt.Println("Task deleted.")
}

// awsCalls is a helper function for using the interactive command-line to pick a option from the menu and then call the appropriate function.
func awsCalls(choice string) {
	mySession := awsConnection()

	switch choice {
	case "view":
		readDynamoTable(mySession)
	case "delete":
		deleteItemDynamoTable(mySession)
	case "add":
		addItemDynamoTable(mySession)
	case "modify":
		modifyItemDynamoTable(mySession)
	}
}

// readUserInput returns string representation with what the user has entered via standard input.
// The trailing newlines are removed before the string is returned, this includes whitespace on Windows OS.
func readUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	if currentRuntime == "windows" {
		userInput = strings.TrimSuffix(userInput, "\r\n")
	} else {
		userInput = strings.TrimSuffix(userInput, "\n")
	}
	return userInput
}

// taskMenu calls the main option menu, this is continually looped
// for the user to select any options they wish to use. Any input which is not in the switch case will cause the menu to refresh 
// the options and is not accepted.
func taskMenu() {
	fmt.Println("1 - Add new task.\n2 - View current tasks.\n3 - Delete completed tasks.\n4 - Modify a task.\n5 - Closes the application.\n")
	for {

		fmt.Print("Select an option menu value: ")
		userChoice := readUserInput()
		switch userChoice {
		case "1":
			awsCalls("add")
		case "2":
			awsCalls("view")
		case "3":
			awsCalls("delete")
		case "4":
			awsCalls("modify")
		case "5":
			exitProgram()
		default:
			fmt.Println("\n-- Error: Enter a numeric value on the menu --\n")
			taskMenu()
		}
	}
}

// exitProgram is used to call for the program to end.
func exitProgram() {
	os.Exit(2)
}

func main() {
	fmt.Println("Task Tracker - Go")
	fmt.Println("Any input not in the menu will refresh the options.")
	taskMenu()
}
