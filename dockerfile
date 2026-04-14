# ---- Build Stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency files first (layer caching)
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary (statically linked)
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ---- Run Stage ----
FROM alpine:3.19

WORKDIR /app

# Copy only the compiled binary
COPY --from=builder /app/server .

# Expose the app port
EXPOSE 8080

CMD ["./server"]