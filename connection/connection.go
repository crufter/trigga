package connection

import(
	"net"
	"crypto/rand"
)

type Connection struct {
	Conn			net.Conn
	Id				[16]byte
}

func New(c net.Conn) *Connection {
	b := make([]byte, 16)
    n, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	if n != 16 {
		panic("No good.")
	}
	id := [16]byte{}
	for i, v := range b {
		id[i] = v
	}
	ret := &Connection{
		c,
		id,
	}
	return ret
}