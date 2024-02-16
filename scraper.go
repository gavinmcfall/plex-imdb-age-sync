package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

// initializing a data structure to keep the scraped data
type videoAgeRating struct {
	Country string
	Rating  string
	IMDbID  string
}

func imdbScraper() {
	// initializing the slice of structs to store the data to scrape
	var videoAgeRatings []videoAgeRating

	// Extract IMDb ID from the URL
	imdbURL := "https://www.imdb.com/title/tt15314262/parentalguide"
	imdbIDParts := strings.Split(imdbURL, "/")
	var imdbID string
	for i, part := range imdbIDParts {
		if part == "title" && i+1 < len(imdbIDParts) {
			imdbID = imdbIDParts[i+1]
			break
		}
	}

	// creating a new Colly instance
	c := colly.NewCollector()

	// scraping logic
	c.OnHTML("section#certificates", func(e *colly.HTMLElement) {
		e.ForEach("ul.ipl-inline-list li.ipl-inline-list__item a[href*=\"/search/title?certificates=NZ:\"]", func(_ int, elem *colly.HTMLElement) {
			// Extract text from the <a> element
			text := elem.Text

			// Split the text by colon
			parts := strings.Split(text, ":")

			// Ensure there are two parts (country and rating)
			if len(parts) == 2 {
				country := strings.TrimSpace(parts[0])
				rating := strings.TrimSpace(parts[1])

				// Create a new videoAgeRating struct
				videoAgeRating := videoAgeRating{
					Country: country,
					Rating:  rating,
					IMDbID:  imdbID,
				}

				// Append the struct to the slice
				videoAgeRatings = append(videoAgeRatings, videoAgeRating)
			}
		})
	})

	// visiting the target page
	c.Visit(imdbURL)

	// opening the CSV file
	file, err := os.Create("ratings.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// initializing a file writer
	writer := csv.NewWriter(file)

	// writing the CSV headers
	headers := []string{
		"IMDbID",
		"Country",
		"Rating",
	}
	writer.Write(headers)

	// writing each videoAgeRating as a CSV row
	for _, videoAgeRating := range videoAgeRatings {
		// converting a videoAgeRating to an array of strings
		record := []string{
			videoAgeRating.IMDbID,
			videoAgeRating.Country,
			videoAgeRating.Rating,
		}

		// adding a CSV record to the output file
		writer.Write(record)
	}
	defer writer.Flush()
}
