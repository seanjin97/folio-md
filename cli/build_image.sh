#!/bin/sh
cd "$(dirname "$0")/.."
docker build -f cli/Dockerfile . -t folio-md:latest
