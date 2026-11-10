package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
	data, err := os.ReadFile("pyq_os.json")
	if err != nil {
		fmt.Println(err)
		panic("error in reading pyq_os.json")
	}
	var Pyqs []PyqEntry
	var PyqEmbeds []PyqEmbed
	err = json.Unmarshal(data, &Pyqs)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	//fmt.Println(len(Pyqs))
	//for _, pyq := range Pyqs {
	//	fmt.Println(pyq, '\n')
	//}
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
}
func embed(texts []string) ([][]float64, error) {
	body, err := json.Marshal(embedRequest{Model: model, Input: texts})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(ollamaURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("is Ollama running? %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var out embedResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Embeddings, nil
}
func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
