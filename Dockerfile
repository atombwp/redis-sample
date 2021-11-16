FROM golang AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o red

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/red ./ok

CMD ["/app/ok"]