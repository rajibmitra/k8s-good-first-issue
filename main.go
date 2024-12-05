package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Issue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"html_url"`
}

func fetchIssues(repo string) ([]Issue, error) {
	baseURL, err := url.Parse("https://api.github.com/repos/")
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	repoPath, err := url.Parse(fmt.Sprintf("%s/issues?labels=good%%20first%%20issue", repo))
	if err != nil {
		return nil, fmt.Errorf("failed to parse repo path: %v", err)
	}

	fullURL := baseURL.ResolveReference(repoPath).String()
	log.Printf("Request URL: %s", fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Response Status: %s", resp.Status)
	log.Printf("Response Headers: %v", resp.Header)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("Response Body: %s", body)

	var issues []Issue
	if err := json.Unmarshal(body, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

func issuesHandler(w http.ResponseWriter, r *http.Request) {
	repos := []string{"kubernetes/kubernetes", "kubernetes/dashboard"}
	var allIssues []Issue

	for _, repo := range repos {
		issues, err := fetchIssues(repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		allIssues = append(allIssues, issues...)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allIssues)
}

func main() {
	http.HandleFunc("/api/issues", issuesHandler)
	http.Handle("/", http.FileServer(http.Dir("./ui")))
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
