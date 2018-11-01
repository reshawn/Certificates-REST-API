
# create static binary and put it in from scratch container
FROM golang:1.9

WORKDIR /go/src/Certificates-REST-API
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -a main.go

FROM scratch
COPY --from=0 /go/src/Certificates-REST-API/main /main
CMD ["/main"]