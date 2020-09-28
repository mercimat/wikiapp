FROM golang:latest
WORKDIR /go/src/github.com/mercimat/wikiapp/
RUN go get -d -v go.mongodb.org/mongo-driver/mongo
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/src/github.com/mercimat/wikiapp/app .
COPY tmpl /root/tmpl/
COPY static /root/static/
CMD ["./app"]
