# familiar — Claude Code session guide

This file is read automatically by Claude Code at the start of every session.
Keep it up to date as the project evolves so each session picks up context fast.

## What this project is

`familiar` is a local-first CLI that translates natural language into shell
commands, shows the command, and waits for confirmation before running anything.
It talks to a local Ollama model (no data leaves the machine).

Full spec: **README.md** — read it before starting any session.

## Current state

**Scaffold only.** The project compiles but does nothing useful yet.

- `go.mod` declares the module; no external dependencies added yet
- `cmd/familiar/main.go` prints help and version — the core loop is a stub
- `internal/ollama/`, `internal/classify/`, `internal/confirm/` are stubs with
  documented TODOs

## Locked decisions

| Decision | Choice | Why |
|---|---|---|
| Language | Go | Single static binary, fast startup, good subprocess/TTY handling |
| Model runtime | Ollama (HTTP API) | Zero-pain local runtime on Apple Silicon, no bundled weights |
| Default model | `qwen2.5-coder:7b` | Code-tuned, strong at small sizes, permissive license |
| CLI shape | One-shot (no REPL) for v0 | Tight scope; REPL is a v0.2 consideration |
| Config file | `~/.config/familiar/config.toml` | XDG convention |
| Distribution | Homebrew tap (post-v0) | `brew install nhughey/tap/familiar` |

## Open questions (decide before implementing)

1. CWD sandbox strictness — hard block or confirm-with-warning for paths outside?
2. Cobra or hand-rolled flag parsing? (Cobra recommended for subcommand ergonomics)

## Building

```sh
# Install Go first (once)
brew install go

# Verify the scaffold compiles
go build ./...

# Run the CLI
go run ./cmd/familiar/ --help
go run ./cmd/familiar/ --version
```

## Next implementation steps (in order)

1. `go get github.com/spf13/cobra` — wire up subcommands properly
2. `internal/ollama/client.go` — `Generate(model, prompt string) (string, error)`
3. `internal/classify/classify.go` — `ReadOnly(cmd string) bool`, `Dangerous(cmd string) bool`
4. `internal/confirm/confirm.go` — raw TTY prompt with Enter/e/c/q keys
5. `familiar init` — check Ollama is running, pull the model
6. `familiar doctor` — print status of Ollama, model, PATH
7. Core loop in `cmd/familiar/main.go` — wire all three internal packages together
8. `familiar explain` — dry-run mode (translate, never execute)
9. `familiar config` — get/set profile, editor

## Package layout rationale

```
cmd/familiar/main.go    entry point; stays thin — just wires up subcommands
internal/ollama/        Ollama HTTP client; internal so it's not importable as a lib
internal/classify/      command safety classifier; isolated so it's easy to test
internal/confirm/       TTY prompt; isolated so it can be tested with a mock reader
```

`internal/` prevents external packages from importing these; everything is wired
through `cmd/familiar/main.go`. This is idiomatic Go for CLI tools.
