# Moverr

Moverr is taking movies with a particular tag in Radarr, moving them to the specified destination and creates a symbolic link from the old path to the new.

## Installation

Install [Go](https://golang.org/doc/install) and run `go get` to install package:

```bash
go get github.com/martinohansen/moverr
```

## Usage

```bash
$ moverr
Usage: moverr (radarr | version) [<args>]

$ moverr radarr
Usage: moverr radarr -t <tag> -d <destination> [-ahps]

Options:
  -a string
        Radarr API key
  -d string
        Destination (required)
  -h string
        Radarr host (default "http://localhost:7878")
  -p string
        Prefix paths with this path
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
