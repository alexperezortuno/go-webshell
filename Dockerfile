FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go mod download

# Build the Go app
RUN go build -a -o /go/bin/app cmd/api/main.go

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["/go/bin/app"]
