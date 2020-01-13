package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

type Tasks struct {
	TaskName    string 
	TaskSummary string 
	CompletedBy  string
}

func optionOne() {
	fmt.Println("Option one selected.")
}

func printJSON(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s) + "\n"
}

func optionTwo() {
	var myTasks Tasks

	jsonTasks, err := os.Open("MyTasks.json")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("JSON file opened...")

	defer jsonTasks.Close()

	byteVals, _ := ioutil.ReadAll(jsonTasks)

	json.Unmarshal(byteVals, &myTasks)
	fmt.Println(printJSON(myTasks))	
	}


func taskMenu() {
	fmt.Println("1 - Add new task.\n2 - View current tasks.")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Select an option value: ")
		userInput, _ := reader.ReadString('\n')
		if runtime.GOOS == "windows" {
			userInput = strings.TrimSuffix(userInput, "\r\n")
		} else {
			userInput = strings.TrimSuffix(userInput, "\n")
		}

		switch userInput {
		case "1":
			optionOne()
		case "2":
			optionTwo()
		}
	}
}

func main() {
	fmt.Println("Task Tracker - Go")
	taskMenu()
}
