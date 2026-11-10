package main

import (
	"encoding/json"
	"net/http"
	"sort"
)

type result struct {
	i     int
	score float64
}

type searchResult struct {
	Score float64  `json:"score"`
	Pyq   PyqEntry `json:"pyq"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Query string `json:"query"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	embeddedQuery, err := embed([]string{body.Query})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	results := make([]result, len(PyqEmbeds))
	for i, p := range PyqEmbeds {
		results[i] = result{i, cosineSearch(embeddedQuery[0], p.Embedding)}
	}
	sort.Slice(results, func(a, b int) bool {
		return results[a].score > results[b].score
	})
	top := min(5, len(results)) // avoids a panic if there are fewer than 5
	out := make([]searchResult, 0, top)
	for _, r := range results[:top] {
		out = append(out, searchResult{Score: r.score, Pyq: Pyqs[r.i]})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}
