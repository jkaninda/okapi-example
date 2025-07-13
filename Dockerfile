########################
# Builder Stage
########################
FROM golang:1.24.4 AS build

WORKDIR /app

# Copy source code
COPY . .

# Download Go dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/okapi-example

########################
# Final Stage
########################
FROM alpine:3.22.0

ENV TZ=UTC

# Install runtime dependencies and set up directories
RUN apk --update --no-cache add tzdata ca-certificates curl

# Copy built binary
COPY --from=build /app/okapi-example /usr/local/bin/okapi-example
RUN chmod a+x /usr/local/bin/okapi-example && ln -s /usr/local/bin/okapi-example /okapi-example

# Expose HTTP Port
EXPOSE 8080

ENTRYPOINT ["/okapi-example"]