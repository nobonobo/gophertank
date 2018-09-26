package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nytimes/gziphandler"

	"github.com/nobonobo/gophertank/api"
	"github.com/nobonobo/gophertank/logger"
)

func main() {
	var dev string
	var listen string
	flag.StringVar(&dev, "dev", "http://localhost:8080", "reverseproxy to gopherjs server")
	flag.StringVar(&listen, "listen", ":8000", "listen address")
	flag.Parse()

	//mime.AddExtensionType(".wasm", "application/wasm")
	ctx := context.Background()

	mux := http.NewServeMux()
	mux.Handle("/ws/", api.New())
	if len(dev) > 0 {
		logger.Println("development mode")
		u, _ := url.Parse(dev)
		rp := httputil.NewSingleHostReverseProxy(u)
		mux.Handle("/", gziphandler.GzipHandler(rp))
	} else {
		logger.Println("normal mode")
		mux.Handle("/", gziphandler.GzipHandler((http.FileServer(http.Dir("./app")))))
	}

	l, err := net.Listen("tcp", listen)
	if err != nil {
		logger.Fatal(err)
	}
	server := &http.Server{
		Handler:  mux,
		ErrorLog: logger.GetErrLogger(),
	}
	serveErrCh := make(chan error)
	go func(e chan<- error) {
		logger.Print("start server listen:", l.Addr())
		if err := server.Serve(l); err != nil {
			e <- err
		}
	}(serveErrCh)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	logger.Println(<-sig, "signal detected ...")
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalln(err)
	}
	if err := <-serveErrCh; err != nil {
		if err != http.ErrServerClosed {
			logger.Fatalln(err)
		}
	}
	logger.Println("terminated")
}
