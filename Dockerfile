# Start from the latest golang base image
FROM golang:1.23-alpine3.20 as builder

ENV PATH /usr/local/go/bin:$PATH
ENV GOLANG_VERSION 1.23


# Add Maintainer Info
LABEL maintainer="cgil"


# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY cmd/goCloudK8sUserGroupServer ./cmd/goCloudK8sUserGroupServer
COPY cmd/goCloudK8sUserGroupServer/goCloudK8sUserGroupFront/dist ./cmd/goCloudK8sUserGroupServer/goCloudK8sUserGroupFront/dist
COPY pkg ./pkg

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goCloudK8sUserGroupServer ./cmd/goCloudK8sUserGroupServer


######## Start a new stage  #######
# using from scratch for size and security reason
# Containers Are Not VMs! Which Base Container (Docker) Images Should We Use?
# https://youtu.be/82ZCJw9poxM
# https://blog.baeke.info/2021/03/28/distroless-or-scratch-for-go-apps/
# https://github.com/vfarcic/base-container-images-demo
FROM scratch

# Add Maintainer Info
LABEL maintainer="cgil"
LABEL org.opencontainers.image.title="goCloudK8sUserGroup"
LABEL org.opencontainers.image.description="This is a goCloudK8sUserGroup container image, Allows to manage users and groups"
LABEL org.opencontainers.image.authors="cgil"
LABEL org.opencontainers.image.licenses="GPL-3.0"

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

# how to check if container is ok https://docs.docker.com/engine/reference/builder/#healthcheck
HEALTHCHECK --start-period=5s --interval=30s --timeout=3s \
    CMD ["curl", "--fail", "http://localhost:9090/health"]


# Command to run the executable
CMD ["./goCloudK8sUserGroupServer"]
