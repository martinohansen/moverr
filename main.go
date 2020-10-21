package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/martinohansen/moverr/internal/copy"
)

var config Config

type Config struct {
	APIKey      string
	PrefixPath  string
	Authority   string
	Destination string
	SymlinkPath string
	Tag         string
}

type Tag struct {
	Label    string `json:"label"`
	MovieIds []int  `json:"movieIds"`
	ID       int    `json:"id"`
}

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

func (tag Tag) Movies() []Movie {
	var movies []Movie
	for _, id := range tag.MovieIds {
		movie := GetMovie(id)
		movies = append(movies, movie)
	}
	return movies
}

func GetTag(label string) Tag {
	url := fmt.Sprintf("%s/api/v3/tag/detail?APIKey=%s", config.Authority, config.APIKey)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var tags []Tag
	json.Unmarshal(body, &tags)

	for _, tag := range tags {
		if tag.Label == label {
			return tag
		}
	}
	emptyTag := new(Tag)
	return *emptyTag
}

func GetMovie(id int) Movie {
	url := fmt.Sprintf("%s/api/v3/movie/%v?APIKey=%s", config.Authority, id, config.APIKey)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var movie Movie
	json.Unmarshal(body, &movie)
	return movie
}

func (movie Movie) IsSymlinked() bool {
	var locations []string
	locations = append(locations, fmt.Sprintf("%s%s", config.PrefixPath, movie.Path))
	locations = append(locations, fmt.Sprintf("%s%s", config.PrefixPath, movie.MovieFile.Path))

	for _, location := range locations {
		fi, _ := os.Lstat(location)
		if fi.Mode()&os.ModeSymlink != 0 {
			return true
		}
	}
	return false
}

func (movie Movie) Move() error {
	folderName := strings.Split(movie.FolderName, "/")
	src := fmt.Sprintf("%s%s", config.PrefixPath, movie.Path)
	dst := fmt.Sprintf("%s/%s", config.Destination, folderName[len(folderName)-1])
	symlink := fmt.Sprintf("%s/%s", config.SymlinkPath, folderName[len(folderName)-1])
	log.Printf("moving %s from %s to %s with symlink %s", movie.Title, src, dst, symlink)
	err := copy.Directory(src, dst)
	if err != nil {
		os.RemoveAll(dst)
		return err
	}
	err = os.RemoveAll(src)
	if err != nil {
		return err
	}
	err = os.Symlink(symlink, src)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.StringVar(&config.PrefixPath, "prefixPath", "./sample", "Prefix for Radarr paths")
	flag.StringVar(&config.APIKey, "apiKey", "", "Radarr API key (required)")
	flag.StringVar(&config.Authority, "authority", "http://localhost:7878", "Radarr host")
	flag.StringVar(&config.Destination, "destination", "./sample/symbolic", "Movie destination after move")
	flag.StringVar(&config.SymlinkPath, "symlinkPath", "../symbolic", "Symlink path")
	flag.StringVar(&config.Tag, "tag", "test", "Radarr tag")
	flag.Parse()

	if config.APIKey == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	tag := GetTag(config.Tag)
	for _, movie := range tag.Movies() {
		if !movie.IsSymlinked() {
			err := movie.Move()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
