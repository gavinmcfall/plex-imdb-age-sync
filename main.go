package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/jrudio/go-plex-client"
)

func main() {

	//Authenticate with Plex
	plexConnection, err := plex.New("http://10.90.3.204:32400", os.Getenv("PLEX_TOKEN"))
	if err != nil {
		fmt.Println("Error connecting to plex", err)
		return
	}
	// Test your connection to your Plex server
	_, err = plexConnection.Test()
	if err != nil {
		fmt.Println("Error testing connection", err)
		return
	}

	//Call Plex and get a list of Libraries: GetLibraries
	PlexLibraries, err := plexConnection.GetLibraries()
	if err != nil {
		fmt.Println("Error testing connection", err)
		return
	}

	//Call Plex and get a list of Libraries: GetLibraries
	allContent, err := AssemblingPlexLibraries(PlexLibraries, plexConnection)
	if err != nil {
		fmt.Println("Error gettig All Content", err)
		return
	}

	// // New Rating
	// newRatings := []Library{}
	// for _, Library := range allContent {
	// 	Library.Content = pullRatings(Library.Content)
	// 	newRatings = append(newRatings, Library)
	// }

	// //Pull Ratings

	// //For each library, get the content
	// for _, Library := range newRatings {
	// 	fmt.Println(Library.Name, Library.Key, len(Library.Content))
	// 	for _, Media := range Library.Content {
	// 		fmt.Println(Media.Title, Media.RatingKey, Media.ContentRating)
	// 	}
	// }

	//Testing DatabaseIDs
	for _, Library := range allContent {
		for _, Media := range Library.Content {
			fmt.Println(Media.Title, Media.Type, GetDatabaseID(Media, "imdb"))
		}
	}
}

// Take the IMDB ID from the metadata ✓

// Given the IMDB ID scrape the Age rating from IMDB ✕

// Take that IMDB Age Rating and use it to update the metadata.contentRating ✕

type Library struct {
	Name    string
	Key     string
	Content []plex.Metadata
}

// Function that takes a plex library and returns a struct of all the libraries and their content
func AssemblingPlexLibraries(Libraries plex.LibrarySections, plexConnection *plex.Plex) ([]Library, error) {
	var results = []Library{}
	for _, Directory := range Libraries.MediaContainer.Directory {
		if Directory.Type == "show" || Directory.Type == "movie" {
			Library := Library{
				Name: Directory.Title,
				Key:  Directory.Key,
			}
			searchResults, err := plexConnection.GetLibraryContent(Directory.Key, "")
			if err != nil {
				return nil, err
			}
			Library.Content = searchResults.MediaContainer.Metadata
			results = append(results, Library)
		}
	}
	return results, nil
}

// Function takes a provider (imdb, tmdb, tvdb) and returns that providers unique ID for that Movie/Show
func GetDatabaseID(metadata plex.Metadata, provider string) string {
	fmt.Println("Grabing Provider for "+metadata.Title, "from "+provider)
	for i, DatabaseID := range metadata.AltGUIDs {
		fmt.Println("Index: ", i)
		fmt.Println("DatabaseID: ", DatabaseID)

	}

	return ""
}

// Function that returns an array of plex metadata to update all the ratings
func pullRatings(metadata []plex.Metadata) []plex.Metadata {
	results := []plex.Metadata{}
	// Take Database ID from GetDatabaseID and pass it to imdbScraper
	for _, Media := range metadata {
		fmt.Println(Media.Title, Media.Type, GetDatabaseID(Media, "imdb"), "Old "+Media.ContentRating)
		if Media.Type == "movie" {
			Media.ContentRating = imdbScraper(GetDatabaseID(Media, "imdb"))
		}
		fmt.Println("New " + Media.ContentRating)
		results = append(results, Media)
	}
	return results
}

// Function that takes a provided IMDB ID and Scrapes the Age Rating
func imdbScraper(titleID string) string {
	// Extract IMDb ID from the URL
	var imdbURL = "https://www.imdb.com/title/" + titleID + "/parentalguide"

	fmt.Println("IMDB ID "+titleID, imdbURL)

	// creating a new Colly instance
	c := colly.NewCollector()

	rating := ""

	// scraping logic
	c.OnHTML("section#certificates", func(e *colly.HTMLElement) {
		e.ForEach("ul.ipl-inline-list li.ipl-inline-list__item a[href*=\"/search/title?certificates=NZ:\"]", func(_ int, elem *colly.HTMLElement) {
			// Extract text from the <a> element
			text := elem.Text
			println("Text: " + text)

			// Split the text by colon
			parts := strings.Split(text, ":")

			// Ensure there are two parts (country and rating)
			if len(parts) == 2 {
				rating = strings.TrimSpace(parts[1])
				fmt.Println("Rating: " + rating)
			}
		})
	})

	// visiting the target page
	err := c.Visit(imdbURL)
	if err != nil {
		fmt.Println("Error scraping IMDB", err)
		return ""
	}
	return rating
}

//For loop over all the directories
//For each directory, extract the directories name, and key and put it into our Libraries struct
//Then for each directory, we are going to pull out the Library contents and put the results of that into our content struct
//Then return that struct
//Which happens inside the AssemblingPlexLibraries function

//Make a function that can take a metadata and return the IMDB ID
//Make a function that can take the IMDB ID and return the Age Rating
//Make a function that take a Metadata and an Age Rating and update the metadata

//function to determine if something already has an age rating

//Massive ass stretch goal
// Conditionally look for new things that have not yet been updated.

//convention in golang is to use a builder container to build the binary
//and to distribute using a scratch container
