FROM golang:1.21.2 AS builder

WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application inside the container
RUN go build -o easypix cmd/main.go

# Stage 2: Create the final lightweight image
FROM scratch

# Copy the built Go binary from the builder stage
COPY --from=builder /app/easypix /app/easypix

# Expose the port on which the web server will listen
EXPOSE 8888

# Define the command to run your Go web server
CMD ["/app/easypix"]