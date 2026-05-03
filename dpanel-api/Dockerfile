# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dpanel-api .

# Runtime stage
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /app/dpanel-api /dpanel-api
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/dpanel-api"]
