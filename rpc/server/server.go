package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"server/pkg/rpc"
	"strings"
)

func main() {
	const addr = "0.0.0.0:9999"
	log.Print("server starting")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("can't listen on %s: %v", addr, err)
	}
	defer listener.Close()
	log.Print("server started")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't accept connection: %v", err)
			continue
		}
		go handleConn(conn)
		//defer conn.Close() -> defer в цикле не делается!
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {
		log.Printf("error while reading: %v", err)
		return
	}
	index := strings.IndexByte(line, ':')
	writer := bufio.NewWriter(conn)
	if index == -1 {
		// TODO: неправильно ввёл команду
		//conn.Write([]byte("error: invalid line"))
		log.Printf("invalid line received %s", line)
		err := rpc.WriteLine("error: invalid line", writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return
		}
		return
	}

	cmd, options := line[:index], line[index+1:]
	log.Printf("command received: %s", cmd)
	log.Printf("options received: %s", options)

	switch cmd {
	case "Upload":
		options := strings.TrimSuffix(options, "\n")
		line, err := rpc.ReadLine(reader)
		if err != nil {
			log.Printf("can't read: %v", err)
			return
		}
		if line == "result: error\n"{
			log.Printf("file not such: %v", err)
			return
		}

		bytes, err := ioutil.ReadAll(reader) // while not EOF
		if err != nil {
			if err != io.EOF {
				log.Printf("can't read data: %v", err)
			}
		}
		err = ioutil.WriteFile("rpc/server/files/" + options, bytes, 0666)
		if err != nil {
			log.Printf("can't write file: %v", err)
		}
	case "Download":
		options = strings.TrimSuffix(options, "\n")
		file, err := os.Open("rpc/server/files/" + options)
		if err != nil {
			log.Print("file does not exist")
			err = rpc.WriteLine("result: error", writer)
			return
		}
		err = rpc.WriteLine("result: ok", writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return
		}
		//name := file.Name()
		_, err = io.Copy(writer, file)
		//writer.Flush()
		case "List":
			options = strings.TrimSuffix(options,"\n")
			fileName := rpc.ReadDir("rpc/"+ options + "/files")
			err := rpc.WriteLine(fileName, writer)
			if err != nil {
			log.Printf("error while writing: %v", err)
			return
			}
	default:
		err := rpc.WriteLine("result: error", writer)
		if err != nil {
			log.Printf("error while writing: %v", err)
			return
		}
	}
}
