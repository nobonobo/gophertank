FROM golang:1.11.1-alpine
RUN apk add -U --no-cache git make
RUN go get github.com/gopherjs/gopherjs
RUN go get -d -u github.com/nobonobo/gophertank
WORKDIR /go/src/github.com/nobonobo/gophertank
RUN make depends
EXPOSE 8080
ENTRYPOINT [ "make", "run" ]
