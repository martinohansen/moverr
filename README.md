# Moverr

Moverr is taking movies with a particular tag in Radarr, moving them to the specified destination and creates a symbolic link from the old path to the new.

## Installation

Install [Go](https://golang.org/doc/install) and run `go get` to install package:

```bash
$ go version
$ go get github.com/martinohansen/moverr

$ moverr
  -apiKey string
        Radarr API key (required)
  -authority string
        Radarr host (default "http://localhost:7878")
  -destination string
        Movie destination after move (default "./sample/symbolic")
  -prefixPath string
        Prefix for Radarr paths (default "./sample")
  -symlinkPath string
        Symlink path (default "../symbolic")
  -tag string
        Radarr tag (default "test")
```

## Usage

```bash
$ moverr \
  -apiKey <Your Radarr API Key>
```

## Contributing

Pull requests are welcome.
