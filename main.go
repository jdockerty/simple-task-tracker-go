package main

import (
"fmt"
"encoding/json"
"io/ioutil"
"bufio"

)

type Task struct {
	TaskName string `json:"TaskName"`

}

func main() {
	fmt.Println("Task Tracker - Go")
	fmt.Println("--- Menu ---")
	fmt.Println("1 - Add new task.\n2 - View current tasks.")
	fmt.Println("Enter a value to select menu options...")

	scanner := bufio.NewScanner(os.Stdin)
}