/*
Package radarr implements the show.Mover interface
*/
package radarr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/martinohansen/moverr/show"
)

// Connection is everything needed to use the Radarr API
type Connection struct {
	APIKey    string
	Authority string
}

// Tag from Radarr API
type Tag struct {
	Label    string `json:"label"`
	MovieIds []int  `json:"movieIds"`
	ID       int    `json:"id"`
}

// Movie from Radarr API
type Movie struct {
	Title               string    `json:"title"`
	OriginalTitle       string    `json:"originalTitle"`
	SortTitle           string    `json:"sortTitle"`
	SizeOnDisk          int       `json:"sizeOnDisk"`
	Status              string    `json:"status"`
	HasFile             bool      `json:"hasFile"`
	Path                string    `json:"path"`
	Monitored           bool      `json:"monitored"`
	MinimumAvailability string    `json:"minimumAvailability"`
	IsAvailable         bool      `json:"isAvailable"`
	FolderName          string    `json:"folderName"`
	CleanTitle          string    `json:"cleanTitle"`
	ImdbID              string    `json:"imdbId"`
	TmdbID              int       `json:"tmdbId"`
	TitleSlug           string    `json:"titleSlug"`
	Tags                []int     `json:"tags"`
	Added               time.Time `json:"added"`
	MovieFile           MovieFile `json:"movieFile"`
	ID                  int       `json:"id"`
}

// MovieFile from Radarr API
type MovieFile struct {
	MovieID             int       `json:"movieId"`
	RelativePath        string    `json:"relativePath"`
	Path                string    `json:"path"`
	Size                int       `json:"size"`
	DateAdded           time.Time `json:"dateAdded"`
	IndexerFlags        int       `json:"indexerFlags"`
	QualityCutoffNotMet bool      `json:"qualityCutoffNotMet"`
	ID                  int       `json:"id"`
}

// Show returns a slice of movies that satisfy the Mover interface
func Show(tag string, conn Connection) []show.Show {
	var shows []show.Show
	t, err := newTag(tag, conn)
	if err != nil {
		log.Fatalf("failed to get tag %s: %s", tag, err)
	}
	movies := t.Movies(conn)
	for _, movie := range movies {
		var show show.Show
		show.Directory = movie.Path
		show.Title = movie.Title
		shows = append(shows, show)
	}
	return shows
}

// Movies returns a slice of all movies with tag
func (tag Tag) Movies(conn Connection) []Movie {
	var movies []Movie
	for _, id := range tag.MovieIds {
		movie, err := newMovie(id, conn)
		if err != nil {
			log.Fatalf("failed to get movie %v: %s", id, err)
		}
		movies = append(movies, *movie)
	}
	return movies
}

// NewTag return the tag with label from conn
func newTag(label string, conn Connection) (*Tag, error) {
	url := fmt.Sprintf("%s/api/v3/tag/detail", conn.Authority)
	if conn.APIKey != "" {
		url = fmt.Sprintf("%s/api/v3/tag/detail?APIKey=%s", conn.Authority, conn.APIKey)
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%s returned %s", conn.Authority, res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	var tags []Tag
	json.Unmarshal(body, &tags)

	for _, tag := range tags {
		if tag.Label == label {
			return &tag, nil
		}
	}
	return nil, fmt.Errorf("didn't find tag")
}

// NewMovie returns the movie with id from conn
func newMovie(id int, conn Connection) (*Movie, error) {
	url := fmt.Sprintf("%s/api/v3/movie/%v?APIKey=%s", conn.Authority, id, conn.APIKey)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	var movie Movie
	json.Unmarshal(body, &movie)
	return &movie, nil
}
