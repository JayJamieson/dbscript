FROM golang:1.24.1-bookworm

# Install mysql-client and ca-certificates
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        default-mysql-client \
        ca-certificates

# Set working directory
WORKDIR /app

# TODO: implement air hot reloading
