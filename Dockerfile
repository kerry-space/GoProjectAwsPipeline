# Start with a base image that includes the Go toolchain.
# Specifying 'alpine' as a variant for a smaller image size.
FROM golang:1.18-alpine as builder

# Install necessary packages like 'git'.
# 'ca-certificates' is often required by applications.
RUN apk update && apk add --no-cache git ca-certificates

# Set the working directory inside the container.
WORKDIR /app

# Copy the local package files to the container's workspace.
COPY . .

# Download all the dependencies that are specified in the go.mod file.
# Using 'go mod tidy' to add missing and remove unused modules.
RUN go mod tidy

# Build the Go app for the target platform specified by the GOARCH environment variable.
# This line assumes you have a main package in the root of your project directory.
# Adjust the path to where your main package's main.go is located if necessary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH:-amd64} go build -o cmd/site

# Use a minimal Alpine image for the production container.
FROM alpine:latest

# Import the compiled binary from the previous stage.
COPY --from=builder /app/cmd/site /app/cmd/site

# Set the port the container listens on.
EXPOSE 8080

# Run the binary.
ENTRYPOINT ["/app/cmd/site"]