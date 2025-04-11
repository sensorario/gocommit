# gocommit

A simple Go CLI tool that stages all changes, commits them with a message, and pushes to your Git repository.

## 📦 Installation

To build and install `gocommit` globally:

```bash
make
```

This will:

- Compile the Go program.
- Move the resulting binary to `/usr/local/bin/gocommit`.
- Overwrite any existing command with the same name.

## 🧼 Clean

To remove the compiled binary from the current directory:

```bash
make clean
```

## 🚀 Usage

Inside any Git repository:

```bash
gocommit
```

You’ll be prompted to enter a commit message. The tool runs:

1. `git add .`
2. `git commit -m "your message"`
3. `git push`
