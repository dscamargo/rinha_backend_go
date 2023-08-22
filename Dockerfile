FROM golang:1.21 as build
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -v -o /bin/app ./cmd/main.go

FROM alpine:3.14.10
RUN apk add dumb-init
WORKDIR /app
COPY --from=build /bin/app ./

ENV GOGC 12300
ENV GOMAXPROCS 3

CMD ["./app"]