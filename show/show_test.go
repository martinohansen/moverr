package show

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/facebookarchive/symwalk"
)

var (
	testRoot string = "/tmp/moverr"
	testSrc  string = fmt.Sprintf("%s/src", testRoot)
	testDst  string = fmt.Sprintf("%s/dst", testRoot)
	showPath string = fmt.Sprintf("%s/movie", testSrc)
	showFile string = fmt.Sprintf("%s/movie.mp4", showPath)
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
	moveCompare(t, x)
	t.Logf("Directory and file list after move:")
	tree(testRoot, t)
}

func TestMoveEmptyDir(t *testing.T) {
	defer cleanupShow(t)
	x := func() {
		os.MkdirAll(showPath, 0755)
		os.MkdirAll(testDst, 0755)
	}
	moveCompare(t, x)
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

func moveCompare(t *testing.T, init func()) {
	init()

	show := Show{Directory: showPath}
	size, err := dirSize(showPath)
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

	newSize, err := dirSize(showPath)
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

// tree prints a list of directories and files from rootPath for human
// verification
func tree(rootPath string, t *testing.T) error {
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			t.Logf(path)
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}
