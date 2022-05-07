##
## Build
##
FROM golang:1.18.1-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /syncovery_collector

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /syncovery_collector /syncovery_collector

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/syncovery_collector"]