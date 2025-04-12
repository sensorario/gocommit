# gocommit

A simple Go CLI tool that stages all changes, commits them with a message, and pushes to your Git repository.

## 📦 Installation

### Requirements

- [Go](https://go.dev/dl/) installed (`go version`)
- `make` and `git` available on your system

### Steps

Clone the repository and install the tool globally:

```bash
git clone https://github.com/YOUR-USERNAME/gocommit.git
cd gocommit
make
```

This will:

- Run `go mod tidy` to install dependencies
- Compile the Go program
- Move the resulting binary to `/usr/local/bin/gocommit`
- Overwrite any existing command with the same name

## 🧼 Clean

To remove the compiled binary from the current directory:

```bash
make clean
```

## 🚀 Usage

Navigate to any Git repository and run:

```bash
gocommit
```

You’ll be prompted to enter a commit message. The tool executes:

1. `git add .`
2. `git commit -m "your message"`
3. `git push`

## 🙋‍♂️ Example

```bash
$ gocommit
Enter commit message: fix: correct typo in README
[main 123abcd] fix: correct typo in README
 1 file changed, 1 insertion(+), 1 deletion(+)
To github.com:your-user/your-repo.git
```

## 🛠 Development

To install or update dependencies manually:

```bash
go mod tidy
```

To build the binary manually:

```bash
go build -o gocommit main.go
```

## 📝 License

MIT
