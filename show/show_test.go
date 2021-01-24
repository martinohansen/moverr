package show

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/facebookarchive/symwalk"
)

var (
	testRoot   string = "/tmp/moverr"
	testSrc    string = fmt.Sprintf("%s/src", testRoot)
	testDst    string = fmt.Sprintf("%s/dst", testRoot)
	showPath   string = fmt.Sprintf("%s/movie", testSrc)
	showFile   string = fmt.Sprintf("%s/movie.mp4", showPath)
	seriesPath string = fmt.Sprintf("%s/series", testSrc)
	seasonPath string = fmt.Sprintf("%s/season 1", seriesPath)
	seriesFile string = fmt.Sprintf("%s/episode 1.mp4", seasonPath)
)

func cleanupShow(t *testing.T) {
	err := os.RemoveAll(testRoot)
	if err != nil {
		t.Errorf("failed to cleanup files: %s", err)
	}
}

func TestMove(t *testing.T) {
	defer cleanupShow(t)
	x := func() {
		os.MkdirAll(showPath, 0755)
		os.MkdirAll(testDst, 0755)
		content := []byte("i'm a movie")
		ioutil.WriteFile(showFile, content, 0755)
	}
	moveCompare(t, x, showPath)
	tree(testRoot, t)
}

func TestMoveEmptyDir(t *testing.T) {
	defer cleanupShow(t)
	x := func() {
		os.MkdirAll(showPath, 0755)
		os.MkdirAll(testDst, 0755)
	}
	moveCompare(t, x, showPath)
}

func TestAlreadyMoved(t *testing.T) {
	defer cleanupShow(t)
	show := Show{Directory: showPath}

	os.MkdirAll(showPath, 0755)
	os.MkdirAll(testDst, 0755)
	content := []byte("i'm a movie whos already moved")
	ioutil.WriteFile(showFile, content, 0755)
	_ = Move(show, testDst, testDst)

	moveable, _ := show.Movable()
	if moveable {
		t.Errorf("failed to see that movie was already moved")
	}
}

func TestMoveSeries(t *testing.T) {
	defer cleanupShow(t)
	x := func() {
		os.MkdirAll(seasonPath, 0755)
		os.MkdirAll(testDst, 0755)
		content := []byte("i'm a episode")
		ioutil.WriteFile(seriesFile, content, 0755)
	}
	moveCompare(t, x, seriesPath)
	tree(testRoot, t)
}

func moveCompare(t *testing.T, init func(), path string) {
	init()

	show := Show{Directory: path}
	size, err := dirSize(path)
	if err != nil {
		t.Errorf("failed to get show size: %s", err)
	}

	moveable, err := show.Movable()
	if moveable {
		err = Move(show, testDst, testDst)
		if err != nil {
			t.Errorf("failed to move show: %s", err)
		}
	} else {
		t.Errorf("show not moveable even thus is should be: %s", err)
	}

	newSize, err := dirSize(path)
	if err != nil {
		t.Errorf("failed to get show size: %s", err)
	}

	if size != newSize {
		t.Errorf("dir size after move doesn't match: %v vs %v bytes", size, newSize)
	}
}

func dirSize(path string) (int64, error) {
	var size int64
	err := symwalk.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// tree checks if the paths are looking correct
func tree(rootPath string, t *testing.T) {
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				t.Error(err)
			}
			if info.Mode()&os.ModeSymlink != 0 {
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}
				checkLink(path, link, t)
			} else {
				t.Logf(path)
			}
			return nil
		})
	if err != nil {
		t.Error(err)
	}
}

// checkLink compares if a path and its symlink have the overall same path
func checkLink(path, link string, t *testing.T) {
	xPath := strings.ReplaceAll(path, "src", "x")
	xLink := strings.ReplaceAll(link, "dst", "x")
	if xPath != xLink {
		t.Errorf("paths are not looking correct:")
	}
	t.Logf("%s -> %s", path, link)
}
