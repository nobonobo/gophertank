PKG:=github.com/nobonobo/gophertank
FRONT:=$(PKG)/app
GOPATH_ORIG:=$(GOPATH)
export GOPATH:=$(shell cd ../../../.. && pwd):$(GOPATH)


build:
	gopherjs build -o ./app/app.js $(FRONT)

run:
	-gopherjs serve --http localhost:8888 $(FRONT) &
	go run ./main.go -dev=http://localhost:8888 -listen=:8080

gopath:
	@echo $(GOPATH)

depends:
	GOPATH=$(GOPATH_ORIG) go get \
		github.com/gopherjs/gopherjs \
		github.com/gopherjs/vecty \
		github.com/google/uuid \
		github.com/gopherjs/websocket \
		golang.org/x/net//websocket \
		honnef.co/go/js/dom \
		github.com/lngramos/three \
		github.com/vecty/vthree \
		github.com/nytimes/gziphandler

docker:
	docker build --rm -t nobonobo/gophertank .
