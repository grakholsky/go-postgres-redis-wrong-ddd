FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/go-postgres-redis/
COPY cmd/ ./cmd/
COPY pkg/ ./pkg
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod vendor
RUN go build -o /go/bin/metric-app ./cmd/metric-app/main.go

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 9090
CMD ["metric-app"]
