/*
Package radarr implements types and functions to work on media from Radarr.
*/
package radarr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/martinohansen/moverr/internal/copy"
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

// Movies returns a slice of all movies with tag
func (tag Tag) Movies(conn Connection) []Movie {
	var movies []Movie
	for _, id := range tag.MovieIds {
		movie, _ := NewMovie(id, conn)
		movies = append(movies, *movie)
	}
	return movies
}

// NewTag return the tag with label from conn
func NewTag(label string, conn Connection) (*Tag, error) {
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
func NewMovie(id int, conn Connection) (*Movie, error) {
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

// Symlinked checks movie paths and returns true if its symlinked
func (movie Movie) Symlinked() (bool, error) {
	locations := []string{movie.Path, movie.MovieFile.Path}

	for _, location := range locations {
		fi, err := os.Lstat(location)
		if err != nil {
			return false, err
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return true, nil
		}
	}
	return false, nil
}

// Move movie to dst and create a symlink from source to sym
func (movie Movie) Move(dst string, sym string, conn Connection) error {
	dir := path.Base(movie.FolderName)
	src := path.Clean(movie.Path)
	dst = path.Join(path.Clean(dst), dir)
	sym = path.Join(path.Clean(sym), dir)

	err := copy.Directory(src, dst)
	if err != nil {
		os.RemoveAll(dst)
		return err
	}
	err = os.RemoveAll(src)
	if err != nil {
		return err
	}
	err = os.Symlink(sym, src)
	if err != nil {
		return err
	}
	return nil
}
