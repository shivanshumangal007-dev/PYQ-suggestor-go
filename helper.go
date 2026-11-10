package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

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

func cosineSearch(a []float64, b []float64) float64 {
	var dot, na, nb float64
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}
