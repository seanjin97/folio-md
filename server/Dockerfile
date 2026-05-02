FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY server/go.mod ./
COPY server/main.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/server .

FROM pandoc/extra:3.9.0.2-debian AS runtime
RUN apt-get update && apt-get install -y --no-install-recommends \
      fonts-noto-color-emoji \
    && rm -rf /var/lib/apt/lists/*
RUN tlmgr install tex-gyre && luaotfload-tool -u
COPY style.tex /style.tex
COPY --from=build /out/server /usr/local/bin/server
EXPOSE 8080
ENTRYPOINT []
CMD ["/usr/local/bin/server"]

FROM runtime AS dev
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH
RUN curl -fsSL https://go.dev/dl/go1.23.6.linux-amd64.tar.gz | tar -C /usr/local -xzf -
WORKDIR /src
RUN go install github.com/air-verse/air@latest
