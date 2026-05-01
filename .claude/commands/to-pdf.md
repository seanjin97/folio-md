Convert $ARGUMENTS to a styled PDF using folio-md.

Run the following shell command:

```sh
/path/to/folio-md/convert.sh $ARGUMENTS
```

Replace `/path/to/folio-md` with the actual path to this repo on disk.

Prerequisites:
- Docker must be running
- The folio-md image must already be built (`./build_image.sh` in the folio-md directory)

The PDF is written to the same directory as the input file with the same base name (e.g. `report.md` → `report.pdf`).
