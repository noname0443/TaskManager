FROM golang:latest AS build
ADD . /go/src
WORKDIR /go/src

RUN go get ./...
ENV CGO_ENABLED=0
RUN go build -o app main.go

FROM scratch AS runtime
COPY --from=build /go/src/app /go/app
EXPOSE 3000/tcp
WORKDIR /go
ENTRYPOINT ["./app"]
