# Use the official Go image from the DockerHub
FROM golang:1.23 as builder
RUN go clean -modcache
# Set the working directory in the image
WORKDIR /app
RUN apt-get install openssl

RUN apt-get update && apt-get install --reinstall ca-certificates && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Copy the go module and sum files
COPY go.mod go.sum ./

# ask for the argument and set gitlab credentials that came frome the docker compose and the .env file
ARG GITLAB_TOKEN

#RUN echo "Token: ${GITLAB_TOKEN}"
RUN git config --global url."https://gitlab-ci-token:${GITLAB_TOKEN}@git.omnicloud.mx/".insteadOf "https://git.omnicloud.mx/"
RUN echo "machine git.omnicloud.mx login gitlab-ci-token password ${GITLAB_TOKEN}" > ~/.netrc


# Set go private for our go modules
ENV GOPRIVATE=git.omnicloud.mx/omnicloud/development/go-modules/*

#Clear cache of go mod
RUN go clean -modcache

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application (with the static c++ binarys)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/migrate-bin migrate/main.go

COPY migrate/migrations /app/migrate/migrations
# assign executable permissions to binaries
RUN chmod +x /app/migrate
RUN chmod +x /app/main

# Use a fresh image to reduce size and only carry the compiled application forward
FROM debian:bullseye-slim

#Set .env variables
ARG GITLAB_TOKEN
ARG DATABASE_DSN
ARG KAFKA_SERVERS
ARG OMNIVIEW_SERVER
ARG OLD-DATABASE_DSN

ENV DATABASE_DSN=$DATABASE_DSN
ENV KAFKA_SERVERS=$KAFKA_SERVERS
ENV OMNIVIEW_SERVER=$OMNIVIEW_SERVER
ENV SENTRY_DSN=$SENTRY_DSN
ENV SENTRY_ENVIRONMENT=$SENTRY_ENVIRONMENT
ENV OLD-DATABASE_DSN=$OLD-DATABASE_DSN

#update aptget and install pgtools to use later in the entrypoint
RUN apt-get update && apt-get install -y postgresql-client \
    && apt-get install -y ca-certificates \
    && apt-get install -y openssl \
    && rm -rf /var/lib/apt/lists/*

ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt

# Set working directory
WORKDIR /opt/ms-collector

# Copy the binary from the builder step and also your config file
COPY --from=builder /app/main /app/main
COPY --from=builder /app/migrate-bin /app/migrate-bin
COPY --from=builder /app/resources/config.yml ./resources/config.yml
COPY --from=builder /app/migrate/migrations /app/migrate/migrations

# assign executable permissions to binaries
RUN chmod +x /app/main
RUN chmod +x /app/migrate-bin

# Copy the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set the binary as the entrypoint of the container
ENTRYPOINT ["/entrypoint.sh"]

# Set the binary as the entrypoint of the container
ENTRYPOINT ["/app/main"]