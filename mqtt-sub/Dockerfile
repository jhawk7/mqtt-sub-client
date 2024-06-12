FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /go/app
COPY . ./
RUN go mod download
RUN mkdir bin
RUN cd cmd/sub/ && go build -o ../../bin/sub-client

FROM golang:1.22.1-alpine3.19
WORKDIR /go
COPY --from=builder /go/app/bin/sub-client .
ENTRYPOINT ["./sub-client"]