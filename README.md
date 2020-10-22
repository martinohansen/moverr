# Moverr

Moverr is taking movies with a particular tag in Radarr, moving them to the specified destination and creates a symbolic link from the old path to the new.

## Installation

Install [Go](https://golang.org/doc/install) and run `go get` to install package:

```bash
$ go get github.com/martinohansen/moverr

$ moverr
Usage: moverr --apiKey <key> <tag> <destination>

see --help for more
```

## Usage

```bash
$ moverr \
  --apiKey <Radarr API key> \
  <Radarr tag> \
  <destination>
```

## Contributing

Pull requests are welcome.
