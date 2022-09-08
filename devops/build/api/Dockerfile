# Dockerfile_old References: https://docs.docker.com/engine/reference/builder/

# Build Default Args
ARG APP_NAME=i-drive

# Start from the latest golang base image
FROM golang:alpine as builder

# Build Args
ARG BUILD_PATH=/build
ARG APP_NAME

# Set the Current Working Directory inside the container
WORKDIR ${BUILD_PATH}

# Copy the source from the current directory to the Working Directory inside the container
COPY src/ .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${APP_NAME} .


######## Start a new stage from scratch #######
FROM alpine:3.10

# Add Maintainer Info
LABEL maintainer="Vadym Titov <vad.titov@gmail.com>"

RUN apk update && \
    apk add --upgrade vips-tools && \
    apk add ca-certificates && rm -rf /var/cache/apk/*

# Build Args
ARG DATA_PATH=/data
ARG BUILD_PATH=/build
ARG APP_NAME

WORKDIR /app

######## Start a new volume stage #######
# Create Directory
RUN mkdir -p ${DATA_PATH}

# Copy server configuration
COPY server.cfg .
COPY profiles.cfg .
COPY migrations ./migrations
COPY resources .${DATA_PATH}

# Copy dump util
COPY devops/bin/boltdb-dump /usr/local/bin/
RUN chmod +x /usr/local/bin/boltdb-dump

# Copy the Pre-built binary file from the previous stage
COPY --from=builder ${BUILD_PATH}/${APP_NAME} /usr/local/bin/

##### install cert
ADD ./devops/build/certs /usr/local/share/ca-certificates
RUN update-ca-certificates

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/usr/local/bin/i-drive"]
