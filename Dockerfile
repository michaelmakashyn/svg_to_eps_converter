# Use the official Golang image as the base image
FROM golang:1.16-alpine

# Install Inkscape
RUN apk --no-cache add inkscape

# Set the working directory inside the container
WORKDIR /app

# Copy the source code from the current directory to the working directory inside the container
COPY .. .

# Build the Go application
RUN go build -o main .

# Expose port 8080 for the web server
EXPOSE 8080

# Set the command to run the binary executable
CMD ["./main"]