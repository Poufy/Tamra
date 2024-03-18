FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Set the environment variable for go. This is to disable the version control system.
#? This is to avoid the error: error obtaining VCS status: exit status 128. Use -buildvcs=false to disable VCS stamping.
# RUN go env -w GOFLAGS="-buildvcs=false"
COPY . .

# COPY go.mod go.sum ./ 
RUN go mod download && go mod verify

# Install go migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install psql
RUN apt-get update && apt-get install -y postgresql-client

COPY ./scripts/wait-for-db.sh /usr/local/bin/wait-for-db.sh

# Run the script and the go tests
RUN chmod +x /usr/local/bin/wait-for-db.sh
RUN ./usr/local/bin/wait-for-db.sh && go test -v ./...

# Run the migrations
# RUN migrate -path ./migrations -database "postgres://postgres:mysecretpassword@db:5432/tamra-postgis-test?sslmode=disable" up

# Run the seeder
# RUN psql -h db -U postgres -d tamra-postgis-test -a -f ./seeds/seed.sql

# Build the application
# RUN go build -o ./bin/tamra ./cmd/tamra/

# Command to run the executable
# CMD ["go", "run", "cmd/tamra/main.go"]
# ENTRYPOINT ["./bin/tamra", "-port=8080", "-db=postgres://postgres:mysecretpassword@db:5432/tamra-postgis?sslmode=disable", "-log-level=debug", "-firebase-config=firebaseConfig.json", "restaurant-logos-bucket=dev-tamra-restaurant-logos"]
# CMD ["air", "-c", ".air.toml"]
