# Start from the latest golang base image
FROM golang:1-alpine3.17 as builder

# Add Maintainer Info
LABEL maintainer="cgil"


# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . ./

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goCloudK8sUserGroupServer ./cmd/goCloudK8sUserGroupServer


######## Start a new stage  #######
FROM scratch
# to comply with security best practices
# Running containers with 'root' user can lead to a container escape situation (the default with Docker...).
# It is a best practice to run containers as non-root users
# https://docs.docker.com/develop/develop-images/dockerfile_best-practices/
# https://docs.docker.com/engine/reference/builder/#user
USER 1221:1221
WORKDIR /goapp

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/goCloudK8sUserGroupServer .

ENV PORT="${PORT}"
ENV DB_DRIVER="${DB_DRIVER}"
ENV DB_HOST="${DB_HOST}"
ENV DB_PORT="${DB_PORT}"
ENV DB_NAME="${DB_NAME}"
ENV DB_USER="${DB_USER}"
ENV DB_PASSWORD="${DB_PASSWORD}"
ENV DB_SSL_MODE="${DB_SSL_MODE}"
ENV ADMIN_USER="${ADMIN_USER}"
ENV ADMIN_PASSWORD="${ADMIN_PASSWORD}"
ENV JWT_SECRET="${JWT_SECRET}"
ENV JWT_DURATION_MINUTES="${JWT_DURATION_MINUTES}"
# Expose port  to the outside world, goCloudK8sUserGroup will use the env PORT as listening port or 8080 as default
EXPOSE 8080

# Command to run the executable
CMD ["./goCloudK8sUserGroupServer"]