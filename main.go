package main

import(
	"net"
	"fmt"
	"github.com/opesun/trigga/rooms"
	"github.com/opesun/trigga/connection"
	"github.com/opesun/trigga/commands"
)

func main() {
	fmt.Println("Starting server on :8912")
	r := rooms.New()
	ln, err := net.Listen("tcp", "127.0.0.1:8912")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			//fmt.Println("Accept err: ", err)
			continue
		}
		//fmt.Println("Client connected.")
		c := connection.New(conn)
		r.RegConn(c)
		go commands.Process(r, c)
	}
}