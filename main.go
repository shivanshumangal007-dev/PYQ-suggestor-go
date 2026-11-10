package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type PyqEntry struct {
	ID       int    `json:"id"`
	Year     int    `json:"year"`
	Exam     string `json:"exam"`
	Marks    int    `json:"marks"`
	Topic    string `json:"topic"`
	Question string `json:"question"`
}

func main() {
	data, err := os.ReadFile("pyq_os.json")
	if err != nil {
		fmt.Println(err)
		panic("error in reading pyq_os.json")
	}
	var Pyqs []PyqEntry

	err = json.Unmarshal(data, &Pyqs)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	//fmt.Println(len(Pyqs))
	//for _, pyq := range Pyqs {
	//	fmt.Println(pyq, '\n')
	//}
}
