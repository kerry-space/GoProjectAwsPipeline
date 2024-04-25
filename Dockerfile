# Use an official Golang runtime as a parent image
FROM golang:alpine AS builder

# Install git, required for fetching Go dependencies.
RUN apk update && apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/mypackage/myapp/

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Fetch dependencies. Using go get.
RUN go get -d -v

# Build the Go app
RUN go build -o /app/cmd/site

# Start a new stage from scratch
FROM scratch

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/cmd/site /site

# Copy other necessary files like templates and configurations
COPY templates/ /templates/
COPY *.yml /

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/site"]