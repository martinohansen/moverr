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
	"github.com/martinohansen/moverr/sonarr"
)

func main() {
	// Root command
	flag.Usage = func() {
		fmt.Println("Usage: moverr (radarr | sonarr | version) [<args>]")
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

	// Sonarr subcommand
	sonarrCmd := flag.NewFlagSet("sonarr", flag.ExitOnError)
	sonarrCmd.Usage = func() {
		fmt.Println("Usage: moverr sonarr -a <key> -t <tag> -d <destination> [-hpsv]")
		fmt.Println("\nOptions:")
		sonarrCmd.PrintDefaults()
	}

	sonarrAPIKey := sonarrCmd.String("a", "", "Sonarr API key (required)")
	sonarrAuthority := sonarrCmd.String("h", "http://localhost:8989", "Sonarr host")
	sonarrDestination := sonarrCmd.String("d", "", "Destination (required)")
	sonarrPrefixPath := sonarrCmd.String("p", "", "Prefix series paths with this path")
	sonarrSymlinkPath := sonarrCmd.String("s", "", "Override symlink path")
	sonarrTag := sonarrCmd.String("t", "", "Tag to move (required)")
	sonarrVerbose := sonarrCmd.Bool("v", false, "Verbose output")

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
	case "sonarr":
		sonarrCmd.Parse(os.Args[2:])

		// Check for required flags
		if *sonarrAPIKey == "" {
			sonarrCmd.Usage()
			os.Exit(1)
		}
		if *sonarrDestination == "" {
			sonarrCmd.Usage()
			os.Exit(1)
		}
		if *sonarrTag == "" {
			sonarrCmd.Usage()
			os.Exit(1)
		}
	case "version":
		fmt.Println("moverr version 0.3.1")
		os.Exit(0)
	default:
		flag.Usage()
	}

	// Radarr subcommand logic
	if radarrCmd.Parsed() {
		radarrConn := &radarr.Connection{
			APIKey:    *radarrAPIKey,
			Authority: *radarrAuthority,
		}

		movies := radarr.Show(*radarrTag, *radarrConn)

		for _, movie := range movies {
			log.SetPrefix(movie.Directory)
			if *radarrPrefixPath != "" {
				movie.Directory = path.Join(*radarrPrefixPath, movie.Directory)
			}

			movable, err := movie.Movable()
			if err != nil {
				log.Fatalf("failed to check if movable: %s", err)
			}

			switch movable {
			case true:
				if *radarrSymlinkPath == "" {
					*radarrSymlinkPath = path.Clean(*radarrDestination)
				} else {
					*radarrSymlinkPath = path.Clean(*radarrSymlinkPath)
				}
				log.Printf("is not moved, moving to: %s and creating symlink to: %s", *radarrDestination, *radarrSymlinkPath)
				err := show.Move(movie, *radarrDestination, *radarrSymlinkPath)
				if err != nil {
					log.Fatalf("failed to move: %s", err)
				}
				log.Printf("finished moving and created symlink")
			case false:
				if *radarrVerbose {
					log.Printf("is already moved, skipping...")
				}
			}
		}
	}

	// Sonar subcommand logic
	if sonarrCmd.Parsed() {
		sonarrConn := &sonarr.Connection{
			APIKey:    *sonarrAPIKey,
			Authority: *sonarrAuthority,
		}

		series, err := sonarr.Show(*sonarrTag, *sonarrConn)
		if err != nil {
			log.Fatalf("failed to any series: %s", err)
		}
		for _, s := range series {
			log.SetPrefix(s.Directory)
			if *sonarrPrefixPath != "" {
				s.Directory = path.Join(*sonarrPrefixPath, s.Directory)
			}

			movable, err := s.Movable()
			if err != nil {
				log.Fatalf("failed to check if movable: %s", err)
			}

			switch movable {
			case true:
				if *sonarrSymlinkPath == "" {
					*sonarrSymlinkPath = path.Clean(*sonarrDestination)
				} else {
					*sonarrSymlinkPath = path.Clean(*sonarrSymlinkPath)
				}
				log.Printf("is not moved, moving to: %s and creating symlink to: %s", *sonarrDestination, *sonarrSymlinkPath)
				err := show.Move(s, *sonarrDestination, *sonarrSymlinkPath)
				if err != nil {
					log.Fatalf("failed to move: %s", err)
				}
				log.Printf("finished moving and created symlink")
			case false:
				if *sonarrVerbose {
					log.Printf("is already moved, skipping...")
				}
			}
		}
	}
}
