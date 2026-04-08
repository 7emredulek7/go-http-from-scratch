package server

import (
	"fmt"
	"httpserver/internal/headers"
	"httpserver/internal/request"
	"httpserver/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	isClosed *atomic.Bool
	handler  Handler
}

func newServer(listener net.Listener, handler Handler) *Server {
	return &Server{
		listener: listener,
		isClosed: &atomic.Bool{},
		handler:  handler,
	}
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := newServer(listener, handler)
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	if s.isClosed.Load() {
		return nil
	}
	s.isClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if s.isClosed.Load() {
			return
		}

		if err != nil {
			return
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	w := response.NewWriter(conn)

	req, err := request.RequestFromReader(conn)
	if err != nil {
		h := headers.NewHeaders()
		h.Set("Connection", "close")
		h.Set("Content-Length", "0")
		w.WriteStatusLine(response.StatusBadRequest)
		w.WriteHeaders(h)
		return
	}

	s.handler(w, req)
}
