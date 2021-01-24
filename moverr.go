/*
Moverr moves shows from content managers to new locations and creates a symbolic
link from the old path to the new.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/martinohansen/moverr/radarr"
	"github.com/martinohansen/moverr/show"
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
		fmt.Println("Usage: moverr radarr -a <key> -t <tag> -d <destination> [-hpsv]")
		fmt.Println("\nOptions:")
		radarrCmd.PrintDefaults()
	}

	radarrAPIKey := radarrCmd.String("a", "", "Radarr API key (required)")
	radarrAuthority := radarrCmd.String("h", "http://localhost:7878", "Radarr host")
	radarrDestination := radarrCmd.String("d", "", "Destination (required)")
	radarrPrefixPath := radarrCmd.String("p", "", "Prefix movie paths with this path")
	radarrSymlinkPath := radarrCmd.String("s", "", "Override symlink path")
	radarrTag := radarrCmd.String("t", "", "Tag to move (required)")
	radarrVerbose := radarrCmd.Bool("v", false, "Verbose output")

	switch os.Args[1] {
	case "radarr":
		radarrCmd.Parse(os.Args[2:])

		// Check for required flags
		if *radarrAPIKey == "" {
			radarrCmd.Usage()
			os.Exit(1)
		}
		if *radarrDestination == "" {
			radarrCmd.Usage()
			os.Exit(1)
		}
		if *radarrTag == "" {
			radarrCmd.Usage()
			os.Exit(1)
		}
	case "version":
		fmt.Println("moverr version 0.2.1")
		os.Exit(0)
	default:
		flag.Usage()
	}

	if radarrCmd.Parsed() {
		radarrConn := &radarr.Connection{
			APIKey:    *radarrAPIKey,
			Authority: *radarrAuthority,
		}

		movies := radarr.Show(*radarrTag, *radarrConn)

		for _, movie := range movies {
			if *radarrPrefixPath != "" {
				movie.Directory = path.Join(*radarrPrefixPath, movie.Directory)
			}

			movable, err := movie.Movable()
			if err != nil {
				log.Fatalf("%s failed to check if movable: %s", movie.Directory, err)
			}

			switch movable {
			case true:
				if *radarrSymlinkPath == "" {
					*radarrSymlinkPath = path.Clean(*radarrDestination)
				} else {
					*radarrSymlinkPath = path.Clean(*radarrSymlinkPath)
				}
				log.Printf("%s is not moved, moving to: %s and creating symlink to: %s", movie.Directory, *radarrDestination, *radarrSymlinkPath)
				err := show.Move(movie, *radarrDestination, *radarrSymlinkPath)
				if err != nil {
					log.Fatalf("%s failed to move: %s", movie.Directory, err)
				}
				log.Printf("%s finished moving and created symlink", movie.Directory)
			case false:
				if *radarrVerbose {
					log.Printf("%s is already moved, skipping...", movie.Directory)
				}
			}
		}
	}
}
