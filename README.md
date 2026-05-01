# folio-md

Convert Markdown files to PDF via a Docker container running Pandoc + LuaLaTeX.

The output is styled to loosely match GitHub's Markdown rendering: TeX Gyre Heros font, GitHub-inspired colors for links, code blocks, and blockquotes, with full emoji support via Noto Color Emoji.

## Why

In the era of coding agents, Markdown has become the default output format. Generating legible PDF reports from Markdown shouldn't require uploading source to opaque web tools. This project is a self-contained solution built entirely on open-source tools — no external services, no hidden code.

## How to give the agent this tool

### Option 1 — CLAUDE.md / AGENT.md snippet (simplest)

Add the following to your project's `CLAUDE.md`:

````markdown
To export a Markdown file as a styled PDF, run from the directory containing the file:

```sh
docker run --rm -v "$(pwd):/data" folio-md <input.md>
```

Docker must be running and the folio-md image must be built first (`docker build . -t folio-md:latest` in the folio-md directory).
````

### Option 2 — Claude Code slash command

Copy `.claude/commands/to-pdf.md` from this repo into your project's `.claude/commands/` directory.

After that, you or the agent can invoke it as:

```
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

Run from the directory containing your file:

```sh
docker run --rm -v "$(pwd):/data" folio-md <input.md>
```

The PDF is written to the same directory as the input file, with the same base name (e.g. `notes.md` → `notes.pdf`).

## Demo

```sh
docker run --rm -v "$(pwd):/data" folio-md test/test.md
```

## License

MIT — see [LICENSE](LICENSE).

## How it works

- `Dockerfile` — extends `pandoc/extra` with Noto Color Emoji fonts and the TeX Gyre font family
- `entrypoint.sh` — baked into the image; runs Pandoc with `lualatex` and `style.tex` as defaults, so no flags are needed at runtime
- `style.tex` — LaTeX preamble baked into the image; configures fonts, page layout, colors, typography, headings, links, code blocks, and blockquotes
