package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/martinohansen/moverr/internal/copy"
)

// Config is populated from flags
type Config struct {
	APIKey      string
	PrefixPath  string
	Authority   string
	SymlinkPath string
	Verbose     bool
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

// Movies returns all movies on tag
func (tag Tag) Movies(cfg Config) []Movie {
	var movies []Movie
	for _, id := range tag.MovieIds {
		movie := NewMovie(id, cfg)
		movies = append(movies, movie)
	}
	return movies
}

// NewTag constructs a tag from label
func NewTag(label string, cfg Config) (Tag, error) {
	url := fmt.Sprintf("%s/api/v3/tag/detail?APIKey=%s", cfg.Authority, cfg.APIKey)
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
			return tag, nil
		}
	}
	emptyTag := new(Tag)
	return *emptyTag, errors.New("found no tag")
}

// NewMovie constructs a movie from ID
func NewMovie(id int, cfg Config) Movie {
	url := fmt.Sprintf("%s/api/v3/movie/%v?APIKey=%s", cfg.Authority, id, cfg.APIKey)
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

// IsSymlinked checks various paths to figure out if movie is already a symlink
func (movie Movie) IsSymlinked(cfg Config) (bool, error) {
	var locations []string
	locations = append(locations, path.Join(cfg.PrefixPath, movie.Path))
	locations = append(locations, path.Join(cfg.PrefixPath, movie.MovieFile.Path))

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

// Move movie from current path to dst and creates symlink
func (movie Movie) Move(dst string, cfg Config) error {
	dir := path.Base(movie.FolderName)
	fmt.Print(dir, movie.FolderName)
	src := path.Join(cfg.PrefixPath, movie.Path)
	dst = path.Join(dst, dir)
	sym := dst
	if cfg.PrefixPath != "" {
		sym = path.Join(cfg.SymlinkPath, dir)
	}

	if cfg.Verbose {
		log.Printf("moving %s from %s to %s with symlink %s...", movie.Title, src, dst, sym)
	}
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

	if cfg.Verbose {
		log.Printf("moved %s successfully", movie.Title)
	}
	return nil
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.PrefixPath, "prefixPath", "", "Prefix Radarr paths")
	flag.StringVar(&cfg.APIKey, "apiKey", "", "Radarr API key")
	flag.StringVar(&cfg.Authority, "authority", "http://localhost:7878", "Radarr host")
	flag.StringVar(&cfg.SymlinkPath, "symlinkPath", "", "Change symlink path (defaults to destination)")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Verbose output")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Printf("Usage: %s --apiKey <key> <tag> <destination>\n\n see --help for more", os.Args[0])
		os.Exit(1)
	}

	tag, err := NewTag(args[0], cfg)
	if err != nil {
		log.Printf("no tag with label %s found", args[0])
	}

	for _, movie := range tag.Movies(cfg) {
		symlink, err := movie.IsSymlinked(cfg)
		if err != nil {
			log.Fatalf("failed to check %s for symlink: %s", movie.Title, err)
		}
		switch symlink {
		case true:
			if cfg.Verbose {
				log.Printf("%s is already moved, skipping...", movie.Title)
			}
		case false:
			if cfg.Verbose {
				log.Printf("%s is not moved, going to...", movie.Title)
			}
			err := movie.Move(args[1], cfg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
