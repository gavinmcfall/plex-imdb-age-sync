package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the configuration structure
type Config struct {
	PlexToken string `yaml:"plex_token"`
}

// Response struct to match the JSON structure
type Response struct {
	MediaContainer struct {
		Directory []struct {
			Key   string `json:"key"`
			Title string `json:"title"`
		} `json:"Directory"`
	} `json:"MediaContainer"`
}

func main() {
	// Open the config file
	configFile, err := os.Open("vars.yaml")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer configFile.Close()

	// Decode the config YAML
	var config Config
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config YAML:", err)
		return
	}

	// Make HTTP GET request to the API
	url := "http://10.90.3.204:32400/library/sections"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set required headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Token", config.PlexToken)

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var data Response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Create a CSV file
	file, err := os.Create("plex_libraries.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header to CSV file
	header := []string{"Title", "Key"}
	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error writing header to CSV file:", err)
		return
	}

	// Write data to CSV file
	for _, item := range data.MediaContainer.Directory {
		err := writer.Write([]string{item.Title, item.Key})
		if err != nil {
			fmt.Println("Error writing data to CSV file:", err)
			return
		}
	}

	fmt.Println("Data has been written to plex_libraries.csv")
}
