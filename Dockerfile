# Use an official Golang runtime as the base image
FROM golang:1.21.1

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main

# Expose the port your application will listen on
EXPOSE 8080

# Define the command to run your application
CMD ["./main"]
