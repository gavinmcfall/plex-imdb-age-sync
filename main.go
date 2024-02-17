package main

import (
	"fmt"
	"os"

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

	allContent, err := AssemblingPlexLibraries(PlexLibraries, plexConnection)
	if err != nil {
		fmt.Println("Error gettig All Content", err)
		return
	}

	for _, Library := range allContent {
		fmt.Println(Library.Name, Library.Key, len(Library.Content))
		for _, Media := range Library.Content {
			fmt.Println(Media.Title, Media.RatingKey)
		}
	}

	//Take those libraries and get a list of all the media: GetLibraryContent

	//Iterate though the media and get the metadata: GetMetadata

	//Take the IMDB ID from the metadata and scrape the Age rating from IMDB

	//Take that IMDB Age Rating and use it to update the media metadata

}

type Library struct {
	Name    string
	Key     string
	Content []plex.Metadata
}

func AssemblingPlexLibraries(Libraries plex.LibrarySections, plexConnection *plex.Plex) ([]Library, error) {
	var results = []Library{}
	for _, Directory := range Libraries.MediaContainer.Directory {
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
	return results, nil
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
