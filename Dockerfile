FROM golang:1.22.1-alpine3.18 

WORKDIR /app
RUN apk add --no-cache make nodejs npm

# Copy the Go module files
COPY . ./

# Download the Go module dependencies
RUN make install
RUN go mod vendor
RUN make build

EXPOSE 3000

ENTRYPOINT [ "./bin/pixelvista" ]