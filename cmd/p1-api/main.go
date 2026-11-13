package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"example.com/p1-reliable-academy-workflows/internal/httpapi"
)

func main() {
	if err := run(net.Listen); err != nil {
		log.Fatal(err)
	}
}

// The listener seam keeps startup-error tests deterministic and socket-free.
func run(listen func(string, string) (net.Listener, error)) error {
	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           httpapi.NewRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	listener, err := listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", server.Addr, err)
	}
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("serve on %s: %w", server.Addr, err)
	}
	return nil
}
