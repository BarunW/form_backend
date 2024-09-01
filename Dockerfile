# Go base image
FROM golang:alpine as builder

# Make a dir app
RUN mkdir /app

# Add the working dir to the container
ADD . /app

# Set the working directory inside the container
WORKDIR /app

# to clean any cache
RUN go clean --modcache

# download all required dependencies 
RUN go mod download

# Build the Go application using cgo enabled and ldflags[ this will provide additional information to Go linker]
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo  -ldflags="-extldflags=-static" -o formpoc .

# Our production used
# Copy your Go source code to the container
FROM alpine:latest AS certificates
RUN apk --no-cache add ca-certificates

#- 
FROM scratch
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-crtificates.crt

COPY --from=builder /app/formpoc .

# Configure go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

EXPOSE 8080
CMD ["./formpoc"]
