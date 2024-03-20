FROM golang:1.22

ARG TEST_DB_USER
ARG TEST_DB_PW
ARG TEST_DB_HOST
ARG TEST_DB_PORT
ARG TEST_DB_NAME

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

COPY ./scripts/wait-for-it.sh /usr/local/bin/wait-for-it.sh
COPY ./scripts/tests-entrypoint.sh /usr/local/bin/tests-entrypoint.sh

RUN chmod +x /usr/local/bin/wait-for-it.sh /usr/local/bin/tests-entrypoint.sh

ENTRYPOINT ["/usr/local/bin/tests-entrypoint.sh"]
