package main

import (
	"distributed_caching_with_leader_election/server"
	"log"
)

func main() {
	s := server.New(":9000")
	log.Fatal(s.Start())
}
