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
	writeToJSONFile(myNewTasks)

}

func jsonFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func writeToJSONFile(taskList Tasks) {
	var newTask Tasks
	jsonFile, err := os.OpenFile("MyTasks.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
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
	err = ioutil.WriteFile("MyTasks.json", taskJSON, 0644)
	fmt.Println("Written to file...")
}

func jsonFormatToString(i interface{}) string {
	jsonData, _ := json.MarshalIndent(i, "", "\t")
	return string(jsonData) + "\n"
}


func viewAllTasks() {
	var myTasks Tasks

	jsonTasks, err := os.Open("MyTasks.json")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("JSON file: '%s' opened\n", jsonTasks.Name())
	fmt.Println("All tasks...")
	defer jsonTasks.Close()

	byteVals, _ := ioutil.ReadAll(jsonTasks)

	json.Unmarshal(byteVals, &myTasks)
	fmt.Println(jsonFormatToString(myTasks))
}

func readUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	if runtime.GOOS == "windows" {
		userInput = strings.TrimSuffix(userInput, "\r\n")
	} else {
		userInput = strings.TrimSuffix(userInput, "\n")
	}
	return userInput
}

func taskMenu() {
	fmt.Println("1 - Add new task.\n2 - View current tasks.")
	for {
		fmt.Print("Select an option value: ")
		switch readUserInput() {
		case "1":
			addNewTasks()
		case "2":
			viewAllTasks()
		case "exit":
			exitProgram()
		}
	}
}

func exitProgram() {
	os.Exit(3)
}

func main() {
	fmt.Println("Task Tracker - Go")
	taskMenu()
}
