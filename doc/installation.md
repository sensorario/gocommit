# Installation

## Requirements

- [Go](https://go.dev/dl/) installed (`go version`)
- `make` and `git` available on your system

### Go installs

- MacOS: `brew install go`

## Steps

Clone the repository and install the tool globally:

```bash
git clone https://github.com/sensorario/gocommit.git
cd gocommit
make
```

This will:

- Run `go mod tidy` to install dependencies
- Compile the Go program
- Move the resulting binary to `/usr/local/bin/gocommit`
- Overwrite any existing command with the same name

## Clean

To remove the compiled binary from the current directory:

```bash
make clean
```