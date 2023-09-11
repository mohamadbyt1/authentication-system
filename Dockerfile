# Use the official Go image as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files into the container
COPY . .

# Build the Go binary
RUN go build -o main .
# Expose the API port
EXPOSE 8080 8080
# Command to run the API when the container starts
CMD ["./main"]
