package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var words = []string{
	"bonfire",
	"cardio",
	"case",
	"character",
	"bonsai",
}

type PrefixResponse struct {
	Keyword string `json:"keyword"`
	Status  string `json:"status"`
	Prefix  string `json:"prefix"`
}

func main() {
	http.HandleFunc("/prefixes", prefixHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func prefixHandler(w http.ResponseWriter, r *http.Request) {
	keywords := r.URL.Query().Get("keywords")
	keywordList := strings.Split(keywords, ",")
	response := []PrefixResponse{}

	for _, keyword := range keywordList {
		found, prefix := findPrefix(keyword)
		if found {
			response = append(response, PrefixResponse{
				Keyword: keyword,
				Status:  "found",
				Prefix:  prefix,
			})
		} else {
			response = append(response, PrefixResponse{
				Keyword: keyword,
				Status:  "not_found",
				Prefix:  "not_applicable",
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func findPrefix(keyword string) (bool, string) {
	for i := 1; i <= len(keyword); i++ {
		prefix := keyword[:i]
		count := 0
		for _, word := range words {
			if strings.HasPrefix(word, prefix) {
				count++
			}
		}
		if count == 1 {
			return true, prefix
		}
	}
	return false, ""
}
