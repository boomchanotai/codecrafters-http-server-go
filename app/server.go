package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleRequest(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte, 1024)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	var res []byte
	path := strings.Split(string(buff), " ")[1]
	userAgent := strings.TrimPrefix(strings.Split(string(buff), "\r\n")[2], "User-Agent: ")
	if path == "/" {
		res = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else if strings.Contains(path, "echo") {
		body := strings.TrimPrefix(path, "/echo/")
		header := fmt.Sprintf("200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d", len(body))
		res = []byte("HTTP/1.1 " + header + "\r\n\r\n" + body)
	} else if path == "/user-agent" {
		header := fmt.Sprintf("200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d", len(userAgent))
		res = []byte("HTTP/1.1 " + header + "\r\n\r\n" + userAgent)
	} else if strings.Contains(path, "files") {
		filePath := strings.TrimPrefix(path, "/files/")
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file: ", err.Error())
		}
		defer file.Close()

		// Read file
		fileInfo, _ := file.Stat()
		fileSize := fileInfo.Size()
		fileBytes := make([]byte, fileSize)
		file.Read(fileBytes)

		header := fmt.Sprintf("200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d", fileSize)
		res = []byte("HTTP/1.1 " + header + "\r\n\r\n" + string(fileBytes))
	} else {
		res = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	_, err = conn.Write(res)
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn)
	}

}
