package main

import (
	"fmt"
	"os"
)

const version = "0.0.0-dev"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--version", "-v":
		fmt.Printf("familiar %s\n", version)
	case "--help", "-h", "help":
		printHelp()
	default:
		// TODO: implement the core NL → command loop
		// See README.md §"The core loop, precisely" for the full spec.
		fmt.Fprintf(os.Stderr, "familiar: core loop not yet implemented\nreceived: %q\n", os.Args[1])
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf(`familiar %s

A natural-language front-end to the shell.

Usage:
  familiar "<what you want>"     translate and confirm a shell command
  familiar explain "<text>"      translate only, never run
  familiar init                  first-run setup — install and pull the model
  familiar config get|set ...    manage configuration
  familiar doctor                check Ollama is up, model present, PATH sane

Flags:
  --dry-run, -n      translate and print, never execute (same as explain)
  --profile <name>   override model profile for this invocation (lite|default|pro)
  --yes, -y          skip confirm for non-destructive commands only
  --verbose          show raw model prompt/response
  --version          print version
  --help             print this help

`, version)
}
