# Go 1.27.0 fica na imagem: o projeto roda sem a linguagem instalada na máquina.
FROM golang:1.27.0

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /usr/local/bin/finder ./cmd/finder

EXPOSE 8085
USER nobody
ENTRYPOINT ["/usr/local/bin/finder"]
