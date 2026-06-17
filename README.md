# Familiar

> A natural-language front-end to the shell for people who never left it.

`familiar` turns plain English into the exact shell command you meant — then
shows it to you and waits. It runs a local model, ships as a single installable
CLI, and never sends your filesystem anywhere.

**Status:** v0 spec (pre-code). Last updated 2026-06-16.

---

## The problem it actually solves

The weak pitch is "natural language for the terminal." That's a toy — anyone who
lives in the CLI already knows `ls -S` and doesn't need an LLM to guess it for
them, lossily and nondeterministically.

The real problem is **inverse discoverability**: you know exactly what you want,
but not the incantation. You lose minutes, daily, to questions like —

- "find every file under here I touched in the last week that's over 10 MB"
- "the `tar` invocation that *excludes* `.git` and follows symlinks"
- "the `ffmpeg` line to strip audio and re-encode to 720p"
- "rename all these `IMG_*.JPEG` to lowercase `.jpg`"

You know `find`, `tar`, and `ffmpeg` exist. You don't remember the flag soup, so
you Google it, or you stop and read `man` for the hundredth time. Familiar
collapses that loop: describe the intent, get the correct command back.

## The one design decision everything hangs on

**Familiar is a translator, not an autonomous agent.** It prints the command it
proposes and waits for you to confirm before anything runs.

```
$ familiar "files over 10MB I changed this week, under the current dir"

  find . -type f -size +10M -mtime -7

  [Enter] run   [e] edit   [c] copy   [q] cancel
```

This is deliberate, and it answers the obvious objection that an LLM is "lossy
and indeterministic." We don't hide the model's output behind side effects — we
*surface* it. You see the command before it touches the disk. The
nondeterminism becomes harmless, and the byproduct is that you actually learn
the flag instead of outsourcing it forever. Familiar is a teacher that happens
to also press Enter for you.

Destructive operations (anything that moves, renames, deletes, or overwrites)
**always** require explicit confirmation and can never be `--yes`-bypassed in v0.

---

## v0 scope

### What it operates on

The **filesystem**, rooted at the current working directory and its subtree.
Nothing else in v0. Specifically:

| Category    | Examples                                                        | v0?            |
|-------------|-----------------------------------------------------------------|----------------|
| Inspect     | find, list by size/date/type, count, `du`, "what's the biggest" | ✅ yes          |
| Navigate    | locate a path, print it, `cd` helper, open in `$EDITOR`/Finder  | ✅ yes          |
| Move/rename | `mv`, batch rename, lowercase extensions, flatten a dir         | ✅ yes (confirm)|
| Edit        | open the right file in your editor at the right spot            | ✅ yes (opens only) |
| Archive     | `tar`/`zip`/`unzip` with the right flags                        | ⏳ v0.1         |
| Code-aware  | "where is `parseConfig` defined" (semantic / AST)               | ❌ later        |
| Notes / RAG | search Obsidian, embeddings over markdown                       | ❌ later        |
| Media       | `ffmpeg`, image conversion                                      | ❌ later        |

**Explicit non-goals for v0:** no semantic search, no embeddings, no indexing,
no code understanding, no notes. Just files, names, sizes, dates, and the
commands that manipulate them. Tight scope = a clean demo = something you can
actually finish.

In-scope content editing is limited to *opening the right file* — Familiar will
not generate file contents in v0. "Edit my config" means "open the config in
`$EDITOR`," not "write new config for me."

### Which model drives it

**Default: `qwen2.5-coder:7b` via [Ollama](https://ollama.com).**

Reasoning: qwen2.5-coder is specifically tuned for code and command generation,
it's strong at small sizes, it's permissively licensed, and Ollama gives us a
zero-pain local runtime on Apple Silicon. The 7B quantized weights are ~4–5 GB
and run comfortably on an M-series Mac.

We do **not** bundle weights inside the binary. Shipping multi-gigabyte weights
in a Homebrew formula is a maintenance nightmare. Instead Familiar depends on
Ollama and pulls the model on first run (`familiar init`).

#### Model profiles ("branches" by size)

You wanted size variants tied to features. Implement them as **profiles**, not
git branches — one codebase, a config flag that selects the model:

| Profile   | Model                  | Footprint | Use case                                        |
|-----------|------------------------|-----------|-------------------------------------------------|
| `lite`    | `qwen2.5-coder:1.5b`   | ~1 GB     | Old hardware, instant latency, simple `find`/`ls` |
| `default` | `qwen2.5-coder:7b`     | ~4.5 GB   | The recommended daily driver                    |
| `pro`     | `qwen2.5-coder:14b`    | ~9 GB     | Hairier multi-step pipelines, better flag recall |

Switch with `familiar config set profile pro`. This is the cleanest way to honor
"different sizes for different features" without fragmenting the project.

### CLI command shape

The headline form is bare natural language as the first argument:

```
familiar "<what you want, in plain English>"
```

Everything else is a subcommand under a verb namespace so the NL path stays
unambiguous:

```
familiar "show the largest file under ~/Applications"   # the core loop
familiar run "<text>"        # explicit alias for the core loop
familiar explain "<text>"    # translate + explain the command, never run it
familiar init                # install/pull the model, first-run setup
familiar config get|set ...  # profile, editor, confirm behavior
familiar doctor              # check Ollama is up, model present, PATH sane
familiar --version | --help
```

Global flags:

```
--dry-run / -n     translate and print only, never execute (same as `explain`)
--profile <name>   override model profile for this one invocation
--yes / -y         skip confirm for NON-destructive commands only
--verbose          show the raw model prompt/response (for the demo + debugging)
```

#### The core loop, precisely

1. Take the NL string.
2. Build a prompt: system instructions + the current working directory + a short
   schema of allowed command families + the user's text.
3. Ask the local model for a single shell command (or a short pipeline).
4. Parse and **classify** it: read-only vs. mutating.
5. Render the command. Read-only with `--yes` may auto-run; mutating always
   prompts.
6. On confirm, execute in a subprocess, stream stdout/stderr, return the exit
   code.

#### A safety rail worth building early

Before showing a command, run it past a small **denylist / classifier**: refuse
or hard-confirm on `rm -rf`, `:(){ :|:& };:`, writes outside the CWD subtree,
`sudo`, `curl ... | sh`, and redirects that clobber. The model is not trusted to
be safe; the wrapper enforces safety. This is also a great thing to *show* in a
talk — "here's how I sandbox an LLM's shell output."

---

## Suggested tech stack (to decide next)

Not locked, but a recommendation so the next session can just start:

- **Language: Go.** Single static binary, trivial `brew install`, fast startup
  (matters for a CLI you invoke constantly), good subprocess and TTY handling.
  Rust is the equally-good alternative; Python makes distribution harder and
  startup slower for a tool like this.
- **Distribution: Homebrew tap.** `brew install nigelhughey/tap/familiar`.
- **Model runtime: Ollama**, talked to over its local HTTP API (`localhost:11434`).
- **Config: `~/.config/familiar/config.toml`.**

## Roadmap

- **v0 (this spec):** confirm-first NL → shell for filesystem inspect / navigate
  / move / open. One model profile working end-to-end.
- **v0.1:** archive/compress commands; `lite` and `pro` profiles; `doctor`.
- **v0.2:** multi-step pipelines, history (`familiar !!` to re-explain the last
  command), shell completion.
- **Later:** code-aware search (AST/ripgrep hybrid), notes/RAG, media (`ffmpeg`).

## Open questions for next session

1. Go or Rust? (Recommendation: Go.)
2. Do we want a REPL mode (`familiar` with no args drops into a prompt), or
   strictly one-shot invocations for v0? (Recommendation: one-shot only.)
3. How strict is the CWD-subtree sandbox — hard block, or confirm-with-warning
   for paths outside it?

---

*Familiar is a learning/portfolio project. The point is a sharp, finishable v0
that demos well, not feature completeness.*
