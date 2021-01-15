/*
Package show implements types and interfaces for moving shows
*/
package show

import (
	"fmt"
	"os"
	"path"

	"github.com/martinohansen/moverr/internal/copy"
)

// Mover is the interface for things that move
type Mover interface {
	Source() string
	Movable() (bool, error)
}

// Show is a folder with a title and X number of media files
type Show struct {
	Directory, Title string
	Media
}

// Media is one or more file paths
type Media []struct {
	File string
}

func (s Show) String() string {
	return fmt.Sprintf("%s", s.Title)
}

// Source returns the source for a show
func (s Show) Source() string {
	return s.Directory
}

// Movable returns true if the show is alright to move
func (s Show) Movable() (bool, error) {
	locations := []string{s.Directory}

	for _, location := range locations {
		fi, err := os.Lstat(location)
		if err != nil {
			return false, err
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return false, fmt.Errorf("%s is already a symlink", fi.Name())
		}
	}
	return true, nil
}

// Move a show and symbolic links the from source to destination
func Move(m Mover, dst, sym string) error {
	movable, err := m.Movable()
	if movable {
		src := path.Clean(m.Source())
		dir, _ := path.Split(src)
		dst = path.Join(path.Clean(dst), dir)
		sym = path.Join(path.Clean(sym), dir)

		err := copy.CopyDir(src, dst)
		if err != nil {
			os.RemoveAll(dst)
			return fmt.Errorf("failed to copy: %s", err)
		}
		err = os.RemoveAll(src)
		if err != nil {
			return fmt.Errorf("failed to remove: %s", err)
		}
		err = os.Symlink(sym, src)
		if err != nil {
			return fmt.Errorf("failed to create symlink: %s", err)
		}
		return nil
	}
	return fmt.Errorf("show not movable: %s", err)
}
