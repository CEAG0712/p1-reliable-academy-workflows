package main

import (
	"errors"
	"net"
	"strings"
	"testing"
)

func TestRunListenFailure(t *testing.T) {
	cause := errors.New("listener unavailable")
	err := run(func(network, address string) (net.Listener, error) {
		if network != "tcp" || address != "127.0.0.1:8080" {
			t.Fatalf("listener target = %q %q, want tcp loopback:8080", network, address)
		}
		return nil, cause
	})
	if !errors.Is(err, cause) || !strings.Contains(err.Error(), "listen") || !strings.Contains(err.Error(), "127.0.0.1:8080") {
		t.Fatalf("startup error lacks cause/address/context: %v", err)
	}
}

type failedListener struct {
	cause  error
	closed bool
}

func (l *failedListener) Accept() (net.Conn, error) { return nil, l.cause }
func (l *failedListener) Close() error              { l.closed = true; return nil }
func (l *failedListener) Addr() net.Addr            { return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080} }

func TestRunServeFailure(t *testing.T) {
	cause := errors.New("accept failed")
	listener := &failedListener{cause: cause}
	err := run(func(string, string) (net.Listener, error) { return listener, nil })
	if !errors.Is(err, cause) || !strings.Contains(err.Error(), "serve") || !strings.Contains(err.Error(), "127.0.0.1:8080") {
		t.Fatalf("serve error lacks cause/address/context: %v", err)
	}
	if !listener.closed {
		t.Fatal("failed server did not close listener")
	}
}
