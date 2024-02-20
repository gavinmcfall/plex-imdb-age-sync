package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gavinmcfall/go-plex-client"
	"github.com/gocolly/colly"
)

func main() {

	//Ratings Fallback
	//US PG-13 == NZ R13
	//US R == NZ R16

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

	// //Call Plex and get a list of Libraries: GetLibraries
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
	newRatings := []Library{}
	for _, Library := range allContent {
		Library.Content, err = pullRatings(plexConnection, Library.Content)
		if err != nil {
			fmt.Println("Error getting Ratings", err)
			return
		}
		newRatings = append(newRatings, Library)
	}

	// //Pull Ratings

	// //For each library, get the content
	for _, Library := range newRatings {
		fmt.Println(Library.Name, Library.Key, len(Library.Content))
		for _, Media := range Library.Content {
			fmt.Println(Media.Title, Media.RatingKey) //Add getting the IMDB ID in here
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
		if Directory.Type == "movie" || Directory.Type == "show" {
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
func GetDatabaseID(plexConnection *plex.Plex, metadata plex.Metadata, provider string) (string, error) {

	// fmt.Println("Grabing Provider for "+metadata.Title, "from "+provider)

	result, err := plexConnection.GetMetadata(metadata.RatingKey)
	if err != nil {
		return "", err
	}
	for _, Provider := range result.MediaContainer.Metadata[0].AltGUIDs {
		DatabaseID := strings.Split(Provider.ID, "://")
		if DatabaseID[0] == provider {
			return DatabaseID[1], nil
		}
		fmt.Println(Provider.ID)
	}

	return "", fmt.Errorf("provider not found %s %s", result.MediaContainer.Metadata[0].Title, result.MediaContainer.Metadata[0].RatingKey)
}

// Function that returns an array of plex metadata to update all the ratings
func pullRatings(plexConnection *plex.Plex, metadata []plex.Metadata) ([]plex.Metadata, error) {
	results := []plex.Metadata{}
	// Take Database ID from GetDatabaseID and pass it to imdbScraper
	for _, Media := range metadata {
		if Media.Type == "movie" {
			DBID, err := GetDatabaseID(plexConnection, Media, "imdb")
			if err != nil {
				return nil, err
			}
			Media.ContentRating = imdbScraper(DBID)
		}
		results = append(results, Media)
	}
	return results, nil
}

// Function that takes a provided IMDB ID and Scrapes the Age Rating
func imdbScraper(titleID string) string {
	// Extract IMDb ID from the URL
	var imdbURL = "https://www.imdb.com/title/" + titleID + "/parentalguide"

	// fmt.Println("IMDB ID "+titleID, imdbURL)

	// creating a new Colly instance
	c := colly.NewCollector()

	nzrating := ""
	usrating := ""

	// scraping logic
	c.OnHTML("section#certificates", func(e *colly.HTMLElement) {
		e.ForEach("ul.ipl-inline-list li.ipl-inline-list__item a[href*=\"/search/title?certificates=NZ:\"]", func(_ int, elem *colly.HTMLElement) {
			// Extract text from the <a> element
			text := elem.Text
			// println("Text: " + text)

			// Split the text by colon
			parts := strings.Split(text, ":")

			// Ensure there are two parts (country and rating)
			if len(parts) == 2 {
				nzrating = strings.TrimSpace(parts[1])
				fmt.Println("NZ Rating: " + nzrating)
			}
		})
		e.ForEach("ul.ipl-inline-list li.ipl-inline-list__item a[href*=\"/search/title?certificates=US:\"]", func(_ int, elem *colly.HTMLElement) {
			// Extract text from the <a> element
			text := elem.Text
			// println("Text: " + text)

			// Split the text by colon
			parts := strings.Split(text, ":")

			// Ensure there are two parts (country and rating)
			if len(parts) == 2 {
				usrating = strings.TrimSpace(parts[1])
				// fmt.Println("US Rating: " + usrating)
			}
		})
	})

	// visiting the target page
	err := c.Visit(imdbURL)
	if err != nil {
		// fmt.Println("Error scraping IMDB", err)
		return ""
	}
	if nzrating == "" {
		nzrating = ratingFallback(usrating)
		// fmt.Println("Overriding Missing NZ Rating with: " + nzrating)
	}
	// println("Final Rating: " + nzrating)
	return nzrating
}

func ratingFallback(usrating string) string {
	if strings.HasPrefix(usrating, "TV-Y") {
		return "G"
	}
	switch usrating {
	case "PG-13", "TV-PG", "TV-14":
		return "R13"
	case "TV-Y", "TV-G":
		return "G"
	case "TV-MA", "R":
		return "R16"
	default:
		return "R18"
	}
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
