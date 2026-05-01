#!/bin/sh
if [ -z "$1" ]; then
  echo "Usage: docker run --rm -v \$(pwd):/data folio-md <input.md>" >&2
  exit 1
fi

BASENAME="$1"
FILENAME="${BASENAME%.*}"

pandoc "/data/${BASENAME}" \
  --pdf-engine=lualatex \
  -H /style.tex \
  -o "/data/${FILENAME}.pdf"
