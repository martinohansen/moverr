package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/martinohansen/moverr/radarr"
)

func main() {
	// Root command
	flag.Usage = func() {
		fmt.Println("Usage: moverr (radarr | version) [<args>]")
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// Radarr subcommand
	radarrCmd := flag.NewFlagSet("radarr", flag.ExitOnError)
	radarrCmd.Usage = func() {
		fmt.Println("Usage: moverr radarr -t <tag> -d <destination> [-ahps]")
		fmt.Println("\nOptions:")
		radarrCmd.PrintDefaults()
	}

	radarrAPIKey := radarrCmd.String("a", "", "Radarr API key")
	radarrAuthority := radarrCmd.String("h", "http://localhost:7878", "Radarr host")
	radarrDestination := radarrCmd.String("d", "", "Destination (required)")
	radarrPrefixPath := radarrCmd.String("p", "", "Prefix paths with this path")
	radarrSymlinkPath := radarrCmd.String("s", "", "Override symlink path")
	radarrTag := radarrCmd.String("t", "", "Tag to move (required)")

	switch os.Args[1] {
	case "radarr":
		radarrCmd.Parse(os.Args[2:])

		// Check for required flags
		if *radarrDestination == "" {
			radarrCmd.Usage()
			os.Exit(1)
		}
		if *radarrTag == "" {
			radarrCmd.Usage()
			os.Exit(1)
		}
	case "version":
		fmt.Println("moverr version 0.1.0")
		os.Exit(0)
	default:
		flag.Usage()
	}

	if radarrCmd.Parsed() {
		radarrConn := &radarr.Connection{
			APIKey:    *radarrAPIKey,
			Authority: *radarrAuthority,
		}

		tag, err := radarr.NewTag(*radarrTag, *radarrConn)
		if err != nil {
			log.Fatalf("failed to get tag %s: %s", *radarrTag, err)
		}

		for _, movie := range tag.Movies(*radarrConn) {
			if *radarrPrefixPath != "" {
				movie.Path = path.Join(*radarrPrefixPath, movie.Path)
				movie.MovieFile.Path = path.Join(*radarrPrefixPath, movie.MovieFile.Path)
			}

			symlink, err := movie.Symlinked()
			if err != nil {
				log.Fatalf("%s failed to check for symlink: %s", movie.Title, err)
			}
			switch symlink {
			case true:
				log.Printf("%s is already moved, skipping...", movie.Title)
			case false:
				if *radarrSymlinkPath == "" {
					*radarrSymlinkPath = path.Clean(*radarrDestination)
				} else {
					*radarrSymlinkPath = path.Clean(*radarrSymlinkPath)
				}
				log.Printf("%s is not moved, moving to: %s and creating symlink to: %s", movie.Title, *radarrDestination, *radarrSymlinkPath)
				err := movie.Move(*radarrDestination, *radarrSymlinkPath, *radarrConn)
				if err != nil {
					log.Fatalf("%s failed to move: %s", movie.Title, err)
				}
				log.Printf("%s finished moving and created symlink", movie.Title)
			}
		}
	}
}
