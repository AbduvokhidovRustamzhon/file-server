package main

import (
	"bufio"
	bytes2 "bytes"
	"fmt"
	"io/ioutil"
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

func Test_UploadToServerOk(t *testing.T) {
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
	_ = rpc.Upd + ":" + options
	src, err := ioutil.ReadFile("files/123.txt")
	if err != nil {
		log.Fatalf("Can't read file: %v",err)
	}
	_, err = writer.Write(src)
	if err != nil {
		log.Fatalf("Can't write: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Can't flush: %v", err)
	}
	err = conn.Close()
	if err != nil {
		log.Fatalf("Can't close conn: %v", err)
	}
	dst, err := ioutil.ReadFile(rpc.WayForClient + options)
	if err != nil {
		log.Fatalf("can't Read file: %v",err)
	}
	if !bytes2.Equal(src, dst) {
		t.Fatalf("files are not equal: %v", err)
	}
}

func Test_ListInServerOk(t *testing.T)  {
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
	options := ""
	line := rpc.List + ":" + options
	err = rpc.WriteLine(line, writer)
	if err != nil {
		t.Fatalf("can't send command %s to server: %v", line, err)
	}
	reader := bufio.NewReader(conn)
	line, err = rpc.ReadLine(reader)
	if line != "123.txt\n" {
		t.Fatalf("result not ok: %s %v", line, err)
	}
}