package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"server/pkg/rpc"
	"testing"
	"time"
)

func Test_DownloadInServerOk(t *testing.T) {
	host := "localhost"
	port := rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err := start(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(rpc.TimeSleep)
	conn, err := net.Dial(rpc.Tcp, addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	writer := bufio.NewWriter(conn)
	options := "123.txt"
	line := rpc.Dwn + ":" + options
	err = rpc.WriteLine(line, writer)
	if err != nil {
		t.Fatalf("can't send command %s to server: %v", line, err)
	}
	reader := bufio.NewReader(conn)
	line, err = rpc.ReadLine(reader)
	log.Print(line)
	if line != "result: ok\n" {
		t.Fatalf("result not ok: %s %v", line, err)
	}
}

func Test_DownloadInServerError(t *testing.T) {
	host := "localhost"
	port := rand.Intn(999) + 9000
	addr := fmt.Sprintf("%s:%d", host, port)
	go func() {
		err := start(addr)
		if err != nil {
			t.Fatalf("can't start server: %v", err)
		}
	}()
	time.Sleep(rpc.TimeSleep)
	conn, err := net.Dial(rpc.Tcp, addr)
	if err != nil {
		t.Fatalf("can't connect to server: %v", err)
	}
	writer := bufio.NewWriter(conn)
	options := "1234.txt"
	line := rpc.Dwn + ":" + options
	err = rpc.WriteLine(line, writer)
	if err != nil {
		t.Fatalf("can't send command %s to server: %v", line, err)
	}
	reader := bufio.NewReader(conn)
	line, err = rpc.ReadLine(reader)
	log.Print(line)
	if line != "result: error\n" {
		t.Fatalf("result not ok: %s %v", line, err)
	}
}