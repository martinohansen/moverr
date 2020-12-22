package show

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var testRoot string = "/tmp/moverr"
var testSrc string = fmt.Sprintf("%s/src", testRoot)
var testDst string = fmt.Sprintf("%s/dst", testRoot)
var showPath string = fmt.Sprintf("%s/movie", testSrc)
var showFile string = fmt.Sprintf("%s/movie.mp4", showPath)

func initShow(t *testing.T) {
	os.MkdirAll(showPath, 0755)
	os.MkdirAll(testDst, 0755)
	content := []byte("i'm a movie")
	ioutil.WriteFile(showFile, content, 0755)
}

func cleanupShow(t *testing.T) {
	err := os.RemoveAll(testRoot)
	if err != nil {
		t.Errorf("failed to cleanup files: %s", err)
	}
}

func TestMove(t *testing.T) {
	initShow(t)
	defer cleanupShow(t)

	show := Show{Directory: showPath}
	oldShow, _ := ioutil.ReadFile(showFile)

	err := Move(show, testDst, testDst)
	if err != nil {
		t.Errorf("failed to move show: %s", err)
	}

	newShow, err := ioutil.ReadFile(showFile)
	if err != nil {
		t.Errorf("failed to read movie after move: %s", err)
	}

	if bytes.Compare(oldShow, newShow) != 0 {
		t.Errorf("move is not the same after move")
	}
}
