# folio-md

Convert Markdown files to PDF via a Docker container running Pandoc + LuaLaTeX.

The output is styled to loosely match GitHub's Markdown rendering: TeX Gyre Heros font, GitHub-inspired colors for links, code blocks, and blockquotes, with full emoji support via Noto Color Emoji.

## Why

In the era of coding agents, Markdown has become the default output format. Generating legible PDF reports from Markdown shouldn't require uploading source to opaque web tools. This project is a self-contained solution built entirely on open-source tools — no external services, no hidden code.

## How to give the agent this tool

### Option 1 — CLAUDE.md/ AGENT.md snippet (simplest)

Add the following to your project's `CLAUDE.md`:

````markdown
To export a Markdown file as a styled PDF, run:

```sh
/path/to/folio-md/convert.sh <input.md>
```

The PDF is written alongside the input file with the same base name. Docker must be running and the folio-md image must be built first (run `./build_image.sh` once in the folio-md directory).
````

### Option 2 — Claude Code slash command

Copy `.claude/commands/to-pdf.md` from this repo into your project's `.claude/commands/` directory, then update the path inside it to point to `convert.sh` on your machine.

After that, you or the agent can invoke it as:

```sh
/to-pdf report.md
```

## Requirements

- Docker

## Setup

Build the Docker image once:

```sh
./build_image.sh
```

## Usage

```sh
./convert.sh <input.md>
```

The PDF is written to the same directory as the input file, with the same base name (e.g. `notes.md` → `notes.pdf`).

## Demo

```sh
./convert.sh test/test.md
```

## License

MIT — see [LICENSE](LICENSE).

## How it works

- `Dockerfile` — extends `pandoc/extra` with Noto Color Emoji fonts and the TeX Gyre font family
- `convert.sh` — mounts the input file's directory and the repo directory into the container, then runs Pandoc with `lualatex` as the PDF engine
- `emoji_header.tex` — LaTeX preamble injected at build time; configures fonts, page layout, colors, typography, headings, links, code blocks, and blockquotes

```

```

```

```

```

```
