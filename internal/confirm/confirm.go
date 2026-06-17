// Package confirm renders the interactive TUI prompt shown after translation.
//
//   find . -type f -size +10M -mtime -7
//
//   [Enter] run   [e] edit   [c] copy   [q] cancel
//
// TODO: implement using golang.org/x/term for raw TTY input so we can read
// a single keypress without requiring the user to hit Enter.
package confirm
