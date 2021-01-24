# Moverr

Moverr is taking movies with a particular tag in Radarr or Sonarr, moving them
to the specified destination and creates a symbolic link from the old path to
the new.

## Installation

Install [Go](https://golang.org/doc/install) and run `go get` to install package:

```bash
go get github.com/martinohansen/moverr
```

## Usage

```bash
$ moverr
Usage: moverr (radarr | sonarr | version) [<args>]

$ moverr radarr # or sonarr
Usage: moverr radarr -a <key> -t <tag> -d <destination> [-hpsv]

Options:
  -a string
        Radarr API key (required)
  -d string
        Destination (required)
  -h string
        Radarr host (default "http://localhost:7878")
  -p string
        Prefix movie paths with this path
  -s string
        Override symlink path
  -t string
        Tag to move (required)
  -v    Verbose output
```

For example:

```bash
moverr radarr \
-a <API Key> \
-d <destination> \
-t <tag>
```

## Contributing

Pull requests are welcome.
