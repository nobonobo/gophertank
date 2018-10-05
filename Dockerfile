FROM golang:1.11.1-alpine
RUN apk add -U --no-cache git make
RUN go get -d github.com/gopherjs/gopherjs
RUN cd $GOPATH/src/github.com/gopherjs/gopherjs;\
    git checkout --track origin/go1.11.1-reflect \
    && go install
RUN go get -d github.com/nobonobo/gophertank
WORKDIR /go/src/github.com/nobonobo/gophertank
RUN make depends
EXPOSE 8000
ENTRYPOINT [ "make", "run" ]
