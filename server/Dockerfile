FROM golang:1.24 AS base

# Set working directory inside the container
WORKDIR /build

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o go-server

# Expose your server's port
EXPOSE 8080

# Run the binary when the container starts
CMD ["./go-server"]
