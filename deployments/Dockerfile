FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Set the environment variable for go. This is to disable the version control system.
#? This is to avoid the error: error obtaining VCS status: exit status 128. Use -buildvcs=false to disable VCS stamping.
# RUN go env -w GOFLAGS="-buildvcs=false"
COPY . .

# COPY go.mod go.sum ./ 
RUN go mod download && go mod verify

# Copy the application's code
# User air for hot reload
# RUN go install github.com/cosmtrek/air@latest

# Build the application
# RUN go build -o ./bin/tamra ./cmd/tamra/

# Command to run the executable
CMD ["go", "run", "cmd/tamra/main.go"]
# ENTRYPOINT ["./bin/tamra", "-port=8080", "-db=postgres://postgres:mysecretpassword@db:5432/tamra-postgis?sslmode=disable", "-log-level=debug", "-firebase-config=firebaseConfig.json"]
# CMD ["air", "-c", ".air.toml"]
