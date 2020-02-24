package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"server/pkg/rpc"
	"strings"
)

func main() {
	file, err := os.Create("server-log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Can't close file: %v", err)
		}
	}()
	log.Print("server starting")
	host := "0.0.0.0"
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "9999"
	}
	err = start(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatal(err)
	}
}

// filepath

func start(addr string) (err error) {
	listener, err := net.Listen(rpc.Tcp, addr)
	if err != nil {
		log.Fatalf("can't listen %s: %v", addr, err)
		return err
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			log.Fatalf("Can't close conn: %v", err)
		}
	}()
	for {
		conn, err := listener.Accept()
		log.Print("accept connection")
		if err != nil {
			log.Fatalf("can't accept: %v", err)
			continue
		}
		log.Print("handle connection")
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) error{
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Can't close conn: %v", err)
		}
	}()
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {

		log.Fatalf("error while reading: %v", err)
		return err
	}
	index := strings.IndexByte(line, ':')
	writer := bufio.NewWriter(conn)
	if index == -1 {
		log.Printf("invalid line received %s", line)
		err := rpc.WriteLine("error: invalid line", writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return err
		}
		return err
	}
	cmd, options := line[:index], line[index+1:]
	log.Printf("command received: %s", cmd)
	log.Printf("options received: %s", options)
	switch cmd {
	case rpc.Upd:
		options := strings.TrimSuffix(options, rpc.Suffix)
		line, err := rpc.ReadLine(reader)
		if err != nil {
			log.Printf("can't read: %v", err)
			return err
		}
		if line == rpc.CheckError + rpc.Suffix {
			log.Printf("file not such: %v", err)
			return err
		}
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("can't read data: %v", err)
				return err
			}
		}
		err = ioutil.WriteFile(rpc.WayForServer+options, bytes, 0666)
		if err != nil {
			log.Printf("can't write file: %v", err)
			return err
		}
		err = rpc.WriteLine(rpc.CheckOk, writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return err
		}
	case rpc.Dwn:
		options = strings.TrimSuffix(options, rpc.Suffix)
		file, err := os.Open(rpc.WayForServer + options)
		if err != nil {
			log.Print("file does not exist")
			err = rpc.WriteLine(rpc.CheckError, writer)
			return err
		}
		err = rpc.WriteLine(rpc.CheckOk, writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return err
		}
		_, err = io.Copy(writer, file)
		err = writer.Flush()
		if err != nil {
			log.Printf("Can't flush: %v", err)
			return err
		}
	case rpc.List:
		options = strings.TrimSuffix(options, rpc.Suffix)
		fileName := rpc.ReadDir(rpc.WayInServer)
		err := rpc.WriteLine(fileName, writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return err
		}
	default:
		err := rpc.WriteLine(rpc.CheckError, writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return err
		}
	}
	return nil
}