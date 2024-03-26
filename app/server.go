package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	buff := make([]byte, 1024)
	_, err = conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	var res []byte
	path := strings.Split(string(buff), " ")[1]
	if path == "/" || strings.Contains(path, "echo") {
		res = []byte("HTTP/1.1 200 OK\r\n")
	} else {
		res = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	content := []byte(strings.Split(path, "/")[2])
	contentType := []byte("Content-Type: text/plain\r\n")
	contentSize := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))
	contentSizeBytes := []byte(contentSize)
	res = append(res, contentType...)
	res = append(res, contentSizeBytes...)
	res = append(res, content...)

	_, err = conn.Write(res)
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}

}
