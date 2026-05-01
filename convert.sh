#!/bin/sh
if [ -z "$1" ]; then
  echo "Usage: $0 <input.md>" >&2
  exit 1
fi

INPUT="$1"
INPUT_DIR="$(dirname "$(realpath "$INPUT")")"
BASENAME="$(basename "$INPUT")"
FILENAME="${BASENAME%.*}"
SCRIPT_DIR="$(dirname "$(realpath "$0")")"

docker run --rm \
  --volume "${INPUT_DIR}:/data" \
  --volume "${SCRIPT_DIR}:/scripts" \
  folio-md:latest \
  "/data/${BASENAME}" \
  --pdf-engine=lualatex \
  -H /scripts/style.tex \
  -o "/data/${FILENAME}.pdf"
