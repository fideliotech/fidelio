# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -tags !dev -ldflags="-w -s" -o fidelio .

# Final stage
FROM scratch AS production
COPY --from=builder /app/fidelio /fidelio
ENV PORT=8080
EXPOSE 8080
CMD ["/fidelio"]