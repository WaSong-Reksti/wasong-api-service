# Use the official Golang image as the base image
FROM golang:1.18-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code into the container
COPY . .

# Install any dependencies required
RUN go mod tidy

# Build the Go application
RUN go build -o main main.go

# Expose the port on which the Go application will run
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]