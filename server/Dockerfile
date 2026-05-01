FROM pandoc/extra:3.9.0.2-debian

RUN apt-get update && apt-get install -y --no-install-recommends \
    fonts-noto-color-emoji \
    && rm -rf /var/lib/apt/lists/*

RUN tlmgr install tex-gyre \
    && luaotfload-tool -u

COPY --from=ghcr.io/astral-sh/uv:0.11.8-python3.14-trixie-slim /usr/local/bin/uv /bin/uv

COPY style.tex /style.tex

WORKDIR /app
COPY server/pyproject.toml .
RUN uv sync --no-install-project

COPY server/server.py .

ENV UV_NO_CACHE=1

EXPOSE 8080
ENTRYPOINT []
CMD ["uv", "run", "--", "uvicorn", "server:app", "--host", "0.0.0.0", "--port", "8080"]
