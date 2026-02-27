package server

import (
	"distributed_caching_with_leader_election/store"
	"log"
	"net"
)

type Server struct {
	addr  string
	store *store.Store
}

func New(addr string) *Server {
	return &Server{
		addr:  addr,
		store: store.New(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	log.Println("Server started on:", s.addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {}
