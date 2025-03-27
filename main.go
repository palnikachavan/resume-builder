package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

const githubAPI = "https://api.github.com/users/"

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
}

// FetchRepositories fetches public repositories of a user
func FetchRepositories(username string, token string) ([]Repo, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "token "+token).
		SetHeader("User-Agent", "Go-GitHub-Extractor").
		Get(githubAPI + username + "/repos")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("GitHub API error: %s", resp.String())
	}

	var repos []Repo
	err = json.Unmarshal(resp.Body(), &repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	githubUsername := "palnikachavan"
	githubToken := os.Getenv("GITHUB_TOKEN") // Get token from .env

	if githubToken == "" {
		log.Fatalf("GITHUB_TOKEN not found in .env")
	}

	repos, err := FetchRepositories(githubUsername, githubToken)
	if err != nil {
		log.Fatalf("Error fetching repos: %v", err)
	}

	for _, repo := range repos {
		fmt.Printf("Project: %s\nDescription: %s\nLanguage: %s\nURL: %s\n\n",
			repo.Name, repo.Description, repo.Language, repo.HTMLURL)
	}
}
