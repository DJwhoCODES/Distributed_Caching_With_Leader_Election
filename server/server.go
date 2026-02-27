package server

import (
	"bytes"
	"distributed_caching_with_leader_election/protocol"
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

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		cmdAny, err := protocol.ParseCommand(conn)
		if err != nil {
			return
		}

		switch cmd := cmdAny.(type) {
		case *protocol.CommandSet:
			s.store.Set(string(cmd.Key), cmd.Value, int64(cmd.TTL))
			s.writeOK(conn)

		case *protocol.CommandGet:
			value, ok := s.store.Get(string(cmd.Key))
			if ok {
				s.writeValue(conn, value)
			} else {
				s.writeNil(conn)
			}
		}
	}
}

func (s *Server) writeOK(conn net.Conn) {
	conn.Write([]byte{1})
}

func (s *Server) writeNil(conn net.Conn) {
	conn.Write([]byte{0})
}

func (s *Server) writeValue(conn net.Conn, value []byte) {
	buf := new(bytes.Buffer)

	buf.WriteByte(2)

	protocol.WriteInt32(buf, int32(len(value)))

	buf.Write(value)

	conn.Write(buf.Bytes())
}
