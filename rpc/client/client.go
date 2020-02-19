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

func operationsLoop(commands string, loop func(cmd string, file string) bool) {
	for {
		fmt.Println(commands)
		var cmd string
		_, err := fmt.Scan(&cmd)
		var file string
		_, err = fmt.Scan(&file)
		if err != nil {
			log.Fatalf("Can't read input: %v", err) // %v - natural ...
		}
		if exit := loop(strings.TrimSpace(cmd), file); exit {
			return
		}
	}
}

func main() {
	operationsLoop(operations, StartingOperationsLoop)
}

func StartingOperationsLoop(cmd string, fileName string) (exit bool) {
	switch cmd {
	case "Download":
		addr := "localhost:9999"
		log.Print("client connecting")
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatalf("can't connect to %s: %v", addr, err)
		}
		defer conn.Close()
		log.Print("client connected")
		writer := bufio.NewWriter(conn)
		line := cmd + ":" + fileName
		log.Print("command sending")
		err = rpc.WriteLine(line, writer)
		if err != nil {
			log.Fatalf("can't send command %s to server: %v", line, err)
		}
		log.Print("command sent")
		downloadFromServer(conn, fileName)
	case "Upload":
		addr := "localhost:9999"
		log.Print("client connecting")
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatalf("can't connect to %s: %v", addr, err)
		}
		defer conn.Close()
		log.Print("client connected")
		writer := bufio.NewWriter(conn)
		line := cmd + ":" + fileName
		log.Print("command sending")
		err = rpc.WriteLine(line, writer)
		if err != nil {
			log.Fatalf("can't send command %s to server: %v", line, err)
		}
		log.Print("command sent")
		uploadInServer(conn, fileName)
	case "List":
		addr := "localhost:9999"
		log.Print("client connecting")
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatalf("can't connect to %s: %v", addr, err)
		}
		defer conn.Close()
		log.Print("client connected")
		writer := bufio.NewWriter(conn)
		line := cmd + ":" + fileName
		log.Print("command sending")
		err = rpc.WriteLine(line, writer)
		if err != nil {
			log.Fatalf("can't send command %s to server: %v", line, err)
		}
		log.Print("command sent")
		listFile(conn)
	case "quit":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}

func downloadFromServer(conn net.Conn, fileName string) {
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {
		log.Printf("can't read: %v", err)
		return
	}
	if line == "result: error\n" {
		log.Printf("file not such: %v", err)
		return
	}
	log.Print(line)
	bytes, err := ioutil.ReadAll(reader) // while not EOF
	if err != nil {
		if err != io.EOF {
			log.Printf("can't read data: %v", err)
		}
	}
	log.Print(len(bytes))
	err = ioutil.WriteFile("rpc/client/files/"+fileName, bytes, 0666)
	if err != nil {
		log.Printf("can't write file: %v", err)
	}

}

func uploadInServer(conn net.Conn, fileName string) {
	options := strings.TrimSuffix(fileName, "\n")
	file, err := os.Open("rpc/client/files/" + options)
	writer := bufio.NewWriter(conn)
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
	log.Print(fileName)

	//name := file.Name()
	fileByte, err := io.Copy(writer, file)
	log.Print(fileByte)
	//writer.Flush()
}

func listFile(conn net.Conn) {
	reader := bufio.NewReader(conn)
	line, err := rpc.ReadLine(reader)
	if err != nil {
		log.Printf("can't read: %v", err)
		return
	}
	var list string
	for i:=0; i< len(line); i++{
		if string(line[i]) == " " || string(line[i]) == "\n"{
			fmt.Println(list)
			list = ""
		} else {
			list = list + string(line[i])
		}
	}
	_, err = ioutil.ReadAll(reader) // while not EOF
	if err != nil {
		if err != io.EOF {
			log.Printf("can't read data: %v", err)
		}
	}
}
