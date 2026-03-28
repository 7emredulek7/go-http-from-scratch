package main

import (
	"fmt"
	"httpserver/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", err)
		}

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", err)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		headers := req.Headers.GetAll()
		fmt.Printf("Headers:\n")
		for key, value := range headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

	}

}
