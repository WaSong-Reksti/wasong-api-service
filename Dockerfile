# Use the official Golang image as the base image
FROM golang:1.18-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Install any dependencies required
RUN go mod tidy

# Build the Go application
RUN go build -o main main.go
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]