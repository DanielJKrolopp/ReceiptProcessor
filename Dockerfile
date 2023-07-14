FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-receipt-processor
EXPOSE 8080

# Run
CMD ["/docker-receipt-processor"]