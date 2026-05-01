import subprocess
import tempfile
import shutil
from pathlib import Path

from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.responses import FileResponse
from starlette.background import BackgroundTask

app = FastAPI()


@app.post("/convert")
async def convert(file: UploadFile = File(...)):
    if not file.filename.endswith(".md"):
        raise HTTPException(status_code=400, detail="File must be a .md file")

    tmpdir = tempfile.mkdtemp()
    try:
        input_path = Path(tmpdir) / "input.md"
        output_path = Path(tmpdir) / "output.pdf"

        input_path.write_bytes(await file.read())

        result = subprocess.run(
            [
                "pandoc", str(input_path),
                "--pdf-engine=lualatex",
                "-H", "/style.tex",
                "-o", str(output_path),
            ],
            capture_output=True,
            text=True,
            timeout=30
        )

        if result.returncode != 0:
            shutil.rmtree(tmpdir, ignore_errors=True)
            print(result.stderr)
            raise HTTPException(status_code=500, detail=result.stderr)

        return FileResponse(
            path=str(output_path),
            media_type="application/pdf",
            filename=Path(file.filename).stem + ".pdf",
            background=BackgroundTask(shutil.rmtree, tmpdir, ignore_errors=True),
        )
    except HTTPException:
        raise
    except Exception as e:
        shutil.rmtree(tmpdir, ignore_errors=True)
        raise HTTPException(status_code=500, detail=str(e))
