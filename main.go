package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
type PyqEmbed struct {
	Id        int       `json:"id"`
	PyqEntry  PyqEntry  `json:"pyqEntry"`
	Embedding []float64 `json:"embedding"`
}
type embedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type embedResponse struct {
	Embeddings [][]float64 `json:"embeddings"`
}

const (
	ollamaURL = "http://localhost:11434/api/embed"
	model     = "nomic-embed-text"
)

var Pyqs []PyqEntry
var PyqEmbeds []PyqEmbed

func main() {
	data, err := os.ReadFile("pyq_os.json")
	if err != nil {
		fmt.Println(err)
		panic("error in reading pyq_os.json")
	}

	err = json.Unmarshal(data, &Pyqs)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	input := make([]string, len(Pyqs))

	for i, pyq := range Pyqs {
		input[i] = pyq.Question
	}

	vecs, err := embed(input)
	if err != nil {
		fmt.Println(err)
	}
	for i, pyq := range Pyqs {
		entry := PyqEmbed{
			Id:        pyq.ID,
			PyqEntry:  pyq,
			Embedding: vecs[i],
		}
		PyqEmbeds = append(PyqEmbeds, entry)
	}
	fmt.Println("pyqs embed is completed successfully")
	//for _, pyq := range PyqEmbeds {
	//	fmt.Println(pyq)
	//}
	mux := http.NewServeMux()
	mux.HandleFunc("/search", searchHandler)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}
