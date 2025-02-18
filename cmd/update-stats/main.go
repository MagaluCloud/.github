package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"
)

type RepoStats struct {
	Stars int `json:"stargazers_count"`
	Forks int `json:"forks_count"`
}

type TemplateData struct {
	Mgccli   RepoStats
	Provider RepoStats
	Examples RepoStats
	Sdk      RepoStats
}

func fetchRepoStats(repo string) (RepoStats, error) {
	url := fmt.Sprintf("https://api.github.com/repos/MagaluCloud/%s", repo)
	resp, err := http.Get(url)
	if err != nil {
		return RepoStats{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RepoStats{}, err
	}

	var stats RepoStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return RepoStats{}, err
	}

	return stats, nil
}

func main() {
	data := TemplateData{}
	mgccli, err := fetchRepoStats("mgccli")
	if err != nil {
		fmt.Printf("Error fetching mgccli stats: %v\n", err)
		return
	}
	data.Mgccli = mgccli

	provider, err := fetchRepoStats("terraform-provider-mgc")
	if err != nil {
		fmt.Printf("Error fetching provider stats: %v\n", err)
		return
	}
	data.Provider = provider

	examples, err := fetchRepoStats("terraform-examples")
	if err != nil {
		fmt.Printf("Error fetching examples stats: %v\n", err)
		return
	}
	data.Examples = examples

	sdk, err := fetchRepoStats("mgc-sdk-go")
	if err != nil {
		fmt.Printf("Error fetching SDK stats: %v\n", err)
		return
	}
	data.Sdk = sdk

	tmpl, err := template.ParseFiles("cmd/update-stats/README.tmpl")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	output, err := os.Create("profile/README.md")
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer output.Close()

	if err := tmpl.Execute(output, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	fmt.Println("README.md has been successfully updated!")
}
