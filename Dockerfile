# Etapa de build
FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/finder ./cmd/finder

# Etapa final
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/finder /finder
EXPOSE 8085
USER nonroot:nonroot
ENTRYPOINT ["/finder"]
