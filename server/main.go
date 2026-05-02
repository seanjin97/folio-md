package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	addr           = ":8081"
	stylePath      = "/style.tex"
	convertTimeout = 30 * time.Second
	maxUploadSize  = 32 << 20
)

// statusRecorder wraps http.ResponseWriter to capture the status code
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rec, r)
		duration := time.Since(start)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "status", rec.status, "duration", duration.String())
	}
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "Invalid multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing 'file' field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(header.Filename, ".md") {
		http.Error(w, "File must be a .md file", http.StatusBadRequest)
		return
	}

	tmpdir, err := os.MkdirTemp("", "folio-md-")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpdir)

	inputPath := filepath.Join(tmpdir, "input.md")
	outputPath := filepath.Join(tmpdir, "output.pdf")

	in, err := os.Create(inputPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(in, file); err != nil {
		in.Close()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	in.Close()

	ctx, cancel := context.WithTimeout(r.Context(), convertTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "pandoc",
		inputPath,
		"--pdf-engine=lualatex",
		"-H", stylePath,
		"-o", outputPath,
	)
	stderr, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("pandoc failed", "error", err, "output", string(stderr))
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			http.Error(w, string(stderr), http.StatusInternalServerError)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := os.Open(outputPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	stem := strings.TrimSuffix(filepath.Base(header.Filename), ".md")
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="`+stem+`.pdf"`)
	if _, err := io.Copy(w, out); err != nil {
		slog.Error("failed to write response", "error", err)
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	slog.Info("hehe")
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", loggingMiddleware(convertHandler))

	slog.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server failed", "error", err)
	}
}
