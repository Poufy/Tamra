FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Set the environment variable for go. This is to disable the version control system.
#? This is to avoid the error: error obtaining VCS status: exit status 128. Use -buildvcs=false to disable VCS stamping.
# RUN go env -w GOFLAGS="-buildvcs=false"

COPY go.mod go.sum ./ 
RUN go mod download && go mod verify

# Copy the application's code
COPY . .
# User air for hot reload
# RUN go install github.com/cosmtrek/air@latest

# Copy go.mod and go.sum and download dependencies
# TODO: copy go.sum as well

# Set the working directory to the location of main.go
WORKDIR /app/cmd/tamra

RUN pwd

# Build the application
RUN go build -o tamra

# Command to run the executable
CMD ["./tamra", "-port=8080", "-db=postgres://postgres:mysecretpassword@host.docker.internal:5432/tamra-postgis?sslmode=disable"]
# CMD ["air", "-c", ".air.toml"]
