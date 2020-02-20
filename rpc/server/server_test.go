package main

import (
	"bufio"
	bytes2 "bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
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

func Test_UploadInServerOk(t *testing.T) {
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
	_ = bufio.NewWriter(conn)
	options := "123.txt"
	_ = rpc.Upd + ":" + options
	file, err := os.OpenFile("files/" + options, os.O_RDONLY, 0666)
	if err != nil {
		t.Fatalf("Can't open file: %v",err)
	}
	openFile, err := os.OpenFile("testFile/"+options,os.O_CREATE|os.O_TRUNC|os.O_RDONLY, 0666)
	if err != nil {
		t.Fatalf("can't create file: %v", err)
	}
	defer func() {
		err = openFile.Close()
		if err != nil {
			t.Fatalf("can't close: %v",err)
		}
	}()
	bytes, err := io.Copy(openFile, file)
	if err != nil {
		t.Fatalf("Can't copy file: %v", err)
	}
	log.Print(bytes)

	fileClient, err := ioutil.ReadFile(rpc.WayForClient + options)
	if err != nil {
		log.Fatalf("can't Read file: %v",err)
	}
	fileServer, err := ioutil.ReadFile("testFile/"+ options)
	if err != nil {
		log.Fatalf("can't Read file: %v",err)
	}
	if !bytes2.Equal(fileClient,fileServer) {
		t.Fatalf("Плохо %s %s %v", fileClient,fileServer, err)
	}
	//conn.Close()
}

func Test_UploadInServerError(t *testing.T) {
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
	_ = bufio.NewWriter(conn)
	options := "1234.txt"
	_ = rpc.Upd + ":" + options
	_, err = os.Open(rpc.WayForClient + options)
	if err == nil {
		t.Fatal("We should not go here.")
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