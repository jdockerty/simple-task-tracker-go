package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Task struct {
	TaskName    string
	TaskSummary string
	CompletedBy string
}

type Tasks []Task

const (
	currentRuntime string = runtime.GOOS
	jsonFileName string = `MyTasks.json`
)

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
		taskNameIn := readUserInput()

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

func jsonFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}


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

func jsonFormatToString(i interface{}) string {
	jsonData, _ := json.MarshalIndent(i, "", "\t")
	return string(jsonData) + "\n"
}

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

// Could maybe use this function with viewAllTasks too? fmt.print(jsonformattostring(readjson())) ??
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

func getIndex(taskList Tasks, taskName string) int {
	for i := range taskList {
		if taskList[i].TaskName == taskName {
			return i
		}
	}
	return -1
}
func deleteTasks() {
	allTasks := readJSONToTasks()
	fmt.Print("Enter the task name to delete: ")
	taskToDelete := readUserInput()
	fmt.Println("Deleting: ", jsonFormatToString(allTasks[getIndex(allTasks, taskToDelete)]))
	allTasks = append(allTasks[:getIndex(allTasks, taskToDelete)], allTasks[getIndex(allTasks, taskToDelete)+1:]...)
	writeToJSONFile(allTasks, true)
}

func taskMenu() {
	fmt.Println("Task Tracker - Go")
	fmt.Println("1 - Add new task.\n2 - View current tasks.\n3 - Delete completed tasks.\nExit - Closes the application.")
	for {
		fmt.Print("Select an option menu value: ")
		switch strings.ToLower(readUserInput()) {
		case "1":
			addNewTasks()
		case "2":
			viewAllTasks()
		case "3":
			deleteTasks()
		case "exit":
			exitProgram()
		}
	}
}

func exitProgram() {
	os.Exit(3)
}

func main() {
	taskMenu()
}
