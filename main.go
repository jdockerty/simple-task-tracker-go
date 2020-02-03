package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// The Task struct defines a particular task, this is used for the JSON representation.
type Task struct {
	TaskName    string
	TaskDetails string
	CompletedBy string
}

// The Tasks type is a slice of the Task struct.
type Tasks []Task

// Constants for the current OS runtime and name of the JSON file.
const (
	currentRuntime string = runtime.GOOS
	jsonFileName   string = `MyTasks.json`
)

func awsSetup() *dynamodb.DynamoDB {
	session, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-2")})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(session)
	return svc
}

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

func deleteItemDynamoTable(dbSession *dynamodb.DynamoDB) {
	fmt.Print("Enter the task name to delete: ")
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

func awsCalls(choice string) {
	mySession := awsSetup()

	switch choice {
	case "view":
		readDynamoTable(mySession)
	case "delete":
		deleteItemDynamoTable(mySession)
	case "add":
		addItemDynamoTable(mySession)
	}
}

// addNewTasks is used to create a number of tasks, that the user specifies, and then write these
// to the JSON file. The user is prompted for the appropriate input.
func addNewTasks() {
	fmt.Println("How many tasks would you like to add?")
	fmt.Print("Enter a value: ")
	taskCount, err := strconv.Atoi(readUserInput())
	if err != nil {
		fmt.Println("String conversion error:", err)
	}
	fmt.Printf("Creating %d tasks...\n", taskCount)
	var myNewTasks Tasks
	for i := 1; i <= taskCount; i++ {
		fmt.Print("Enter a task name: ")
		taskNameIn := strings.ToLower(readUserInput())

		fmt.Print("Enter a summary of the task: ")
		taskSummaryIn := readUserInput()

		fmt.Print("Enter the date that the task should be completed by: ")
		completedByIn := readUserInput()

		newTask := Task{taskNameIn, taskSummaryIn, completedByIn}
		myNewTasks = append(myNewTasks, newTask)
	}
	fmt.Println("Tasks added:", jsonFormatToString(myNewTasks))
	writeToJSONFile(myNewTasks, false)

}

// writeToJSONFile will write to the JSON file containing the tasks, or append to it if the appendToFile
// variable is set to true.
func writeToJSONFile(taskList Tasks, appendToFile bool) {
	if appendToFile {
		jsonFile, err := os.OpenFile(jsonFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error opening:", err)
			exitProgram()
		}
		defer jsonFile.Close()
		taskJSON, _ := json.Marshal(taskList)
		err = ioutil.WriteFile(jsonFileName, taskJSON, 0644)
		fmt.Println("Written to file...")
	} else {
		var newTask Tasks
		jsonFile, err := os.OpenFile(jsonFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error opening:", err)
			exitProgram()
		}
		defer jsonFile.Close()
		bytes, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println("Error reading:", err)
			exitProgram()
		}
		err = json.Unmarshal(bytes, &newTask)
		for _, task := range taskList {
			newTask = append(newTask, task)
		}
		taskJSON, _ := json.Marshal(newTask)
		err = ioutil.WriteFile(jsonFileName, taskJSON, 0644)
		fmt.Println("Written to file...")
	}
}

// jsonFormatToString returns string with the standard JSON indentation, this is used for printing JSON neatly to the console.
func jsonFormatToString(i interface{}) string {
	jsonData, _ := json.MarshalIndent(i, "", "\t")
	return string(jsonData) + "\n"
}

// viewAllTasks opens the JSON file and prints the location of the file for debugging purposes and the tasks contained within it.
func viewAllTasks() {
	var myTasks Tasks

	jsonTasks, err := os.Open(jsonFileName)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("JSON file location: '%s'\n", jsonTasks.Name())
	fmt.Println("All tasks...")
	defer jsonTasks.Close()

	byteVals, _ := ioutil.ReadAll(jsonTasks)

	json.Unmarshal(byteVals, &myTasks)
	fmt.Println(jsonFormatToString(myTasks))
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

// readJSONToTasks returns the Tasks slice. This opens the file and reads the contents of the file to a byte array.
// This byte array is then unmarshalled into the empty Tasks variable and returned.
func readJSONToTasks() Tasks {
	var currentTasks Tasks
	jsonFile, err := os.Open(jsonFileName)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(bytes, &currentTasks)
	return currentTasks
}

// getIndex returns an integer value for the index position of a given taskName in the slice of Tasks.
// This works by iterating over the tasks within the slice and matching the appropriate index value
// with the parameter of the task name which is being searched for.
func getIndex(taskList Tasks, taskName string) int {
	for i := range taskList {
		if taskList[i].TaskName == taskName {
			return i
		}
	}
	return -1
}

// deleteTasks utilises the readJSONToTasks() function to read all of the current tasks into a Tasks struct.
// Deletion is done through finding the index of the task specified, through getIndex(), and appending to the
// slice before the given index, then after the given index + 1 and above. This removes the element from the slice.
// This change is then written to the JSON file to reflect the deletion, passing the new slice with the relevant removal.
func deleteTasks() {
	allTasks := readJSONToTasks()
	fmt.Print("Enter the task name to delete: ")
	taskToDelete := strings.ToLower(readUserInput())
	taskIndex := getIndex(allTasks, taskToDelete)
	fmt.Println("Deleting: ", jsonFormatToString(allTasks[taskIndex]))
	allTasks = append(allTasks[:taskIndex], allTasks[taskIndex+1:]...)
	writeToJSONFile(allTasks, true)
}

// taskMenu calls the main option menu, this is continually looped
// for the user to select any options they wish to use.
func taskMenu() {
	fmt.Println("Task Tracker - Go")
	fmt.Println("1 - Add new task.\n2 - View current tasks.\n3 - Delete completed tasks.\n4 - Test AWS stuff.\nExit - Closes the application.")
	for {
		fmt.Print("Select an option menu value: ")
		switch strings.ToLower(readUserInput()) {
		case "1":
			awsCalls("add")
			//addNewTasks()
		case "2":
			//viewAllTasks()
			awsCalls("view")
		case "3":
			//deleteTasks()
			awsCalls("delete")
		// case "4":
		// 	testAWS()
		case "exit":
			exitProgram()
		}
	}
}

// exitProgram is used to call for the program to end.
func exitProgram() {
	os.Exit(2)
}

func main() {
	taskMenu()
}
