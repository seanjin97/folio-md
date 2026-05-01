FROM pandoc/extra:3.9.0.2-debian

RUN apt-get update && apt-get install -y --no-install-recommends \
    fonts-noto-color-emoji \
    && rm -rf /var/lib/apt/lists/*

RUN tlmgr install tex-gyre \
    && luaotfload-tool -u

COPY style.tex /style.tex
COPY cli/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN useradd --create-home --shell /bin/sh appuser
USER appuser

ENTRYPOINT ["/entrypoint.sh"]
