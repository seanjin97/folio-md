# folio-md

Convert Markdown files to PDF via a Docker container running Pandoc + LuaLaTeX.

The output is styled to loosely match GitHub's Markdown rendering: TeX Gyre Heros font, GitHub-inspired colors for links, code blocks, and blockquotes, with full emoji support via Noto Color Emoji.

## Why

In the era of coding agents, Markdown has become the default output format. Generating legible PDF reports from Markdown shouldn't require uploading source to opaque web tools. This project is a self-contained solution built entirely on open-source tools — no external services, no hidden code.

## Requirements

- Docker

---

## CLI

Convert a Markdown file to PDF directly from your terminal.

### Setup

Build the image once:

```sh
./cli/build_image.sh
```

### Usage

Run from the directory containing your file:

```sh
docker run --rm -v "$(pwd):/data" folio-md <input.md>
```

The PDF is written to the same directory as the input file, with the same base name (e.g. `notes.md` → `notes.pdf`).

### Demo

```sh
docker run --rm -v "$(pwd):/data" folio-md test/test.md
```

### Give this tool to an agent

**Option 1 — CLAUDE.md / AGENT.md snippet (simplest)**

Add the following to your project's `CLAUDE.md`:

````markdown
To export a Markdown file as a styled PDF, run from the directory containing the file:

```sh
docker run --rm -v "$(pwd):/data" folio-md <input.md>
```

Docker must be running and the folio-md image must be built first (`./cli/build_image.sh` in the folio-md directory).
````

**Option 2 — Claude Code slash command**

Copy `.claude/commands/to-pdf.md` from this repo into your project's `.claude/commands/` directory.

After that, you or the agent can invoke it as:

```
/to-pdf report.md
```

---

## Server

A REST API that accepts a Markdown file and returns a PDF. Useful for integrating with web apps or other services.

### Run locally

```sh
cd server
docker compose up --build
```

### Convert a file

```sh
curl -F "file=@input.md" http://localhost:8080/convert -o output.pdf
```

### API

`POST /convert`

| Field | Type | Description |
| ----- | ---- | ----------- |
| `file` | multipart file | Input `.md` file |

Returns the PDF as `application/pdf`.

### Build the server image

Build context must be the repo root so both `style.tex` and `server/` are accessible:

```sh
docker build -f server/Dockerfile . -t folio-md-server:latest
```

---

## How it works

| File | Role |
| ---- | ---- |
| `style.tex` | LaTeX preamble baked into both images; controls fonts, page layout, colors, typography, headings, links, code blocks, and blockquotes |
| `cli/Dockerfile` | Extends `pandoc/extra`; installs fonts; bakes in `entrypoint.sh` |
| `cli/entrypoint.sh` | Runs Pandoc with `lualatex` and `style.tex` — no flags needed at runtime |
| `server/server.py` | FastAPI app; single `POST /convert` endpoint; calls pandoc as subprocess; returns PDF binary |
| `server/Dockerfile` | Extends `pandoc/extra`; adds uv + Python deps; overrides entrypoint to run uvicorn |

## License

MIT — see [LICENSE](LICENSE).
