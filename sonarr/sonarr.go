/*
Package sonarr implements the show.Mover interface
*/
package sonarr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/martinohansen/moverr/show"
)

// Connection is everything needed to use the Sonarr API
type Connection struct {
	APIKey    string
	Authority string
}

// Series from Sonarr API
type Series struct {
	Title         string    `json:"title"`
	OriginalTitle string    `json:"originalTitle"`
	SortTitle     string    `json:"sortTitle"`
	Path          string    `json:"path"`
	Monitored     bool      `json:"monitored"`
	CleanTitle    string    `json:"cleanTitle"`
	TitleSlug     string    `json:"titleSlug"`
	Tags          []int     `json:"tags"`
	Added         time.Time `json:"added"`
	ID            int       `json:"id"`
}

// Tag from Sonarr API
type Tag struct {
	Label string `json:"label"`
	ID    int    `json:"id"`
}

// Show returns a slice of movies that satisfy the Mover interface
func Show(tag string, conn Connection) ([]show.Show, error) {
	var shows []show.Show
	t, err := tagID(tag, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag %s: %s", tag, err)
	}
	series, err := t.Series(conn)
	if err != nil {
		return nil, err
	}
	for _, s := range series {
		var show show.Show
		show.Directory = s.Path
		show.Title = s.Title
		shows = append(shows, show)
	}
	return shows, nil
}

// Series returns a slice of all series with tag
func (tag Tag) Series(conn Connection) ([]Series, error) {
	url := fmt.Sprintf("%s/api/series?APIKey=%s", conn.Authority, conn.APIKey)

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

	var slice []Series
	var series []Series
	json.Unmarshal(body, &series)

	for _, s := range series {
		for _, t := range s.Tags {
			if t == tag.ID {
				slice = append(slice, s)
			}
		}
	}
	return slice, nil
}

// tagID returns the tag ID for a given label
func tagID(label string, conn Connection) (*Tag, error) {
	url := fmt.Sprintf("%s/api/tag?APIKey=%s", conn.Authority, conn.APIKey)

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
