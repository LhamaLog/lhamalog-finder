package main

import (
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServidorAtendeHealth(t *testing.T) {
	s := servidor(":8085")
	if s.Addr != ":8085" {
		t.Fatalf("addr = %s", s.Addr)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = s.Serve(ln) }()
	t.Cleanup(func() { _ = s.Close() })

	var resp *http.Response
	url := "http://" + ln.Addr().String() + "/health"
	for i := 0; i < 20; i++ {
		resp, err = http.Get(url)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	corpo, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 || string(corpo) != `{"status":"UP"}` {
		t.Fatalf("status=%d corpo=%s", resp.StatusCode, corpo)
	}
}
