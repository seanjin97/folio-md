# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

folio-md converts Markdown files to styled PDFs via a Docker container running Pandoc + LuaLaTeX. The output is styled to loosely match GitHub's Markdown rendering (TeX Gyre Heros font, GitHub-inspired colors for links, code blocks, blockquotes, full emoji support via Noto Color Emoji).

## CLI usage

### Build the image

```sh
./cli/build_image.sh
# equivalent to: docker build -f cli/Dockerfile . -t folio-md:latest
```

### Convert a Markdown file to PDF

Run from the directory containing the file:

```sh
docker run --rm -v "$(pwd):/data" folio-md <input.md>
```

Output is written alongside the input file with the same base name (`notes.md` → `notes.pdf`).

Demo using the included test file:

```sh
docker run --rm -v "$(pwd):/data" folio-md test/test.md
```

## Server usage

### Run locally (dev)

```sh
cd server
docker compose up --build
# POST http://localhost:8080/convert  (multipart: file=<input.md>)
# Response: PDF binary
```

### Build the server image

Build context must be the repo root so both `style.tex` and `server/` are accessible:

```sh
docker build -f server/Dockerfile . -t folio-md-server:latest
```

## Architecture

| File                         | Role                                                                                                                                |
| ---------------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `style.tex`                  | LaTeX preamble baked into both images; controls fonts, page layout, colors, typography, headings, links, code blocks, and blockquotes |
| `cli/Dockerfile`             | Extends `pandoc/extra:3.9.0.2-debian`; installs fonts; bakes in `entrypoint.sh`                                                    |
| `cli/entrypoint.sh`          | Runs `pandoc` with `lualatex` and `style.tex` — no flags needed at runtime                                                          |
| `cli/build_image.sh`         | Wrapper around `docker build` for the CLI image                                                                                     |
| `server/server.py`           | FastAPI app; single `POST /convert` endpoint; calls pandoc as subprocess; returns PDF binary                                        |
| `server/Dockerfile`          | Extends `pandoc/extra:3.9.0.2-debian`; adds uv + Python deps; runs uvicorn                                                         |
| `server/docker-compose.yml`  | Local dev setup; mounts `server.py` for hot reload via `--reload`                                                                   |
| `server/pyproject.toml`      | Pinned Python dependencies managed by uv                                                                                            |
| `.devcontainer/devcontainer.json` | Opens VS Code inside the server container for IDE support without local Python                                                 |
| `.claude/commands/to-pdf.md` | Claude Code slash command (`/to-pdf <file>`) that wraps the CLI docker run invocation                                               |

Style changes go in `style.tex`. CLI runtime behavior lives in `cli/entrypoint.sh`. Server logic lives in `server/server.py`. Dependency changes (fonts, TeX packages) go in the respective `Dockerfile`.
