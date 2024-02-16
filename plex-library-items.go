package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

// LibAuthVar represents the LibAuthVaruration structure
type LibAuthVar struct {
	PlexToken string `yaml:"plex_token"`
}

// plexLibraryItemsData struct to match the JSON structure
type plexLibraryItemsData struct {
	MediaContainer struct {
		Metadata []struct {
			RatingKey string `json:"ratingKey"`
			Title     string `json:"title"`
		} `json:"Metadata"`
	} `json:"MediaContainer"`
}

func plexLibraryItems() {
	// Open the LibAuthVar file
	LibAuthVarFile, err := os.Open("vars.yaml")
	if err != nil {
		fmt.Println("Error opening LibAuthVar file:", err)
		return
	}
	defer LibAuthVarFile.Close()

	// Decode the LibAuthVar YAML
	var LibAuthVar LibAuthVar
	decoder := yaml.NewDecoder(LibAuthVarFile)
	err = decoder.Decode(&LibAuthVar)
	if err != nil {
		fmt.Println("Error decoding LibAuthVar YAML:", err)
		return
	}

	// Make HTTP GET request to the API
	url := "http://10.90.3.204:32400/library/sections/14/all"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set required headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Plex-Token", LibAuthVar.PlexToken)

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON plexLibraryItemsData
	var data plexLibraryItemsData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Create a CSV file
	file, err := os.Create("plex_library-items.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header to CSV file
	header := []string{"Title", "RatingKey"}
	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error writing header to CSV file:", err)
		return
	}

	// Write data to CSV file
	for _, item := range data.MediaContainer.Metadata {
		err := writer.Write([]string{item.Title, item.RatingKey})
		if err != nil {
			fmt.Println("Error writing data to CSV file:", err)
			return
		}
	}

	fmt.Println("Data has been written to plex_library_items.csv")
}
