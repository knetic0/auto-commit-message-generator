# Auto Commit Message Generator

A lightweight CLI tool written in Go that leverages Google’s Gemini generative AI to automatically craft meaningful Git commit messages based on your uncommitted changes. Simply run the tool in your repository, review the suggested message, and press Enter to commit or Esc to cancel.

# Features

- **AI-Powered Commit Messages**
Uses Google Gemini (via HTTP request) to generate clear, conventional commit style messages from your diffs.

- **Zero External Dependencies**
Relies only on Go’s standard library (net/http, encoding/json, os/exec, etc.).

- **Interactive Confirmation**
Presents the generated message and waits for Enter (to proceed) or Esc (to abort).

- **Handles Initial Commit**
Detects if your repo has only one commit and uses git diff-tree accordingly.

- **Cross-Platform**
Works on macOS, Linux, and Windows (with Go and Git installed).

# Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/knetic0/auto-commit-message-generator.git
    ```

2. Build the binary:
   ```bash
   cd auto-commit-message-generator
   go build -o generate-commit ./cmd/app
   ```

3. Move the binary to your PATH:
   ```bash
   Mac/Linux:
   sudo mv generate-commit /usr/local/bin/
    ```

# Configuration

```bash
export GEMINI_API_KEY="your_api_key_here"
```

# Usage

Run the tool in your Git repository:

```bash
cd path/to/your/git/repo
generate-commit
```

# Project Structure

```
auto-commit-message-generator/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── git/
│   │   ├── gitsys.go
│   └── gemini/
│       ├── client.go
│       └── types.go
├── go.mod
├── go.sum
└── README.md
```
