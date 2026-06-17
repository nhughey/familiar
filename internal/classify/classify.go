// Package classify decides whether a shell command is safe to auto-run
// or requires explicit confirmation.
package classify

// ReadOnly returns true if cmd is a pure inspection command (find, ls, du, cat…).
// When ReadOnly returns true and --yes is set, familiar may skip the prompt.
//
// TODO: implement using a keyword allowlist approach first, then refine.
// Conservative default: if unsure, return false (mutating).

// TODO: implement Dangerous — hard-block commands regardless of --yes.
// Initial denylist (see README §"A safety rail worth building early"):
//   rm -rf /   fork bomb   sudo   curl|sh   redirects that clobber outside CWD
