FROM golang AS build
WORKDIR /go/src
COPY . .

RUN go get ./...
ENV CGO_ENABLED=0
RUN go build -o app main.go

FROM scratch AS runtime
COPY --from=build /go/src/app /go/app
EXPOSE 8080/tcp
WORKDIR /go
ENTRYPOINT ["./app"]
