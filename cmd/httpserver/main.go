package main

import (
	"httpserver/internal/headers"
	"httpserver/internal/request"
	"httpserver/internal/response"
	"httpserver/internal/server"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const port = 42069

func main() {
	srv, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
			h := headers.NewHeaders()
			h.Set("Content-Type", "text/html")
			h.Set("Content-Length", strconv.Itoa(len(body)))
			h.Set("Connection", "close")
			w.WriteStatusLine(response.StatusBadRequest)
			w.WriteHeaders(h)
			w.WriteBody(body)
		case "/myproblem":
			body := []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
			h := headers.NewHeaders()
			h.Set("Content-Type", "text/html")
			h.Set("Content-Length", strconv.Itoa(len(body)))
			h.Set("Connection", "close")
			w.WriteStatusLine(response.StatusInternalServerError)
			w.WriteHeaders(h)
			w.WriteBody(body)
		default:
			body := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
			h := headers.NewHeaders()
			h.Set("Content-Type", "text/html")
			h.Set("Content-Length", strconv.Itoa(len(body)))
			h.Set("Connection", "close")
			w.WriteStatusLine(response.StatusOK)
			w.WriteHeaders(h)
			w.WriteBody(body)
		}
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer srv.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
