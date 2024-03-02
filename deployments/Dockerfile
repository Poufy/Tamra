FROM golang:1.22

# Set the working directory inside the container
WORKDIR /usr/src/

# Copy the application's code
COPY . .

# User air for hot reload
RUN go install github.com/cosmtrek/air@latest

# Copy go.mod and go.sum and download dependencies
# TODO: copy go.sum as well
COPY go.mod ./ 
RUN go mod download && go mod verify

# Set the working directory to the location of main.go
WORKDIR /usr/src/cmd/tamra

# Build the application
# RUN go build -o tamra

# Command to run the executable
# CMD ["./tamra"]
CMD ["air", "-c", ".air.toml"]
