package rpc

import (
	"bufio"
	"io/ioutil"
	"log"
)

func ReadLine(reader *bufio.Reader) (line string, err error) {
	return reader.ReadString('\n')
}

func WriteLine(line string, writer *bufio.Writer) (err error) {
	_, err = writer.WriteString(line + "\n")
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	return
}

func ReadDir(line string) (fileList string) {
	files, err := ioutil.ReadDir(line)
	if err != nil {
		log.Printf("Can't read dir: %v", err)
	}
	for _, file := range files {
		if fileList == "" {
			fileList = fileList + file.Name()
		} else {
			fileList = fileList + " " + file.Name()
		}
	}
	fileList = fileList + "\n"
	return fileList
}


const Dwn = "Download"
const Upd = "Upload"
const List = "List"
const Quit = "quit"
const CheckError = "result: error"
const WayForServer = "files/"
const CheckOk = "result: ok"
const Addr = "0.0.0.0:9999"
const Tcp = "tcp"
const Suffix = "\n"
const WayInServer = "files"
const AddrClient = "localhost:9999"
const WayForClient = "files/"

const TimeSleep = 1_000_000_000