Convert $ARGUMENTS to a styled PDF using folio-md.

Run the following from the directory containing the file:

```sh
docker run --rm -v "$(pwd):/data" folio-md $ARGUMENTS
```

Prerequisites:
- Docker must be running
- The folio-md image must already be built (`docker build . -t folio-md:latest` in the folio-md directory)

The PDF is written to the same directory as the input file with the same base name (e.g. `report.md` → `report.pdf`).
