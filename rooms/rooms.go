package rooms

import(
	cn "github.com/opesun/trigga/connection"
	"github.com/opesun/trigga/binhelper"
	"sync"
	"fmt"
	"encoding/json"
)

type Rooms struct {
	conns			map[[16]byte]*cn.Connection
	roomToConns		map[string]map[[16]byte]struct{}	// This tells us what connection.Connections should we forward the publish to.
	connToRooms		map[[16]byte]map[string]struct{}	// Bookkeping field, at connection close this tells us from which rooms should we remove the connection.Connection.
	mut				*sync.RWMutex
}

func (c *Rooms) RegConn(conn *cn.Connection) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.conns[conn.Id] = conn
	c.connToRooms[conn.Id] = map[string]struct{}{}
}

func (c *Rooms) Subscribe(roomName string, conn *cn.Connection) {
	c.mut.Lock()
	defer c.mut.Unlock()
	_, exists := c.roomToConns[roomName]
	if !exists {
		c.roomToConns[roomName] = map[[16]byte]struct{}{}
	}
	c.roomToConns[roomName][conn.Id] = struct{}{}
	_, ex := c.connToRooms[conn.Id]
	if !ex {
		c.connToRooms[conn.Id] = map[string]struct{}{}
	}
	c.connToRooms[conn.Id][roomName] = struct{}{}
}

func (c *Rooms) Unsubscribe(roomName string, conn *cn.Connection) {
	c.mut.Lock()
	defer c.mut.Unlock()
	room, roomExists := c.roomToConns[roomName]
	if !roomExists {
		return
	}
	_, connIsInRoom := room[conn.Id]
	if !connIsInRoom {
		return
	}
	connToRooms, connHasRoom := c.connToRooms[conn.Id]
	if !connHasRoom {
		panic("Inconsistent state: conn should have room" + roomName)
	}
	delete(room, conn.Id)
	delete(connToRooms, roomName)
}

func (c *Rooms) Remove(conn *cn.Connection) {
	c.mut.Lock()
	defer c.mut.Unlock()
	delete(c.conns, conn.Id)
	for i := range c.connToRooms[conn.Id] {
		delete(c.roomToConns[i], conn.Id)
	}
	delete(c.connToRooms, conn.Id)
	conn.Conn.Close()
}

func (c *Rooms) Publish(roomName string, msg []byte, except *cn.Connection) {
	c.mut.Lock()
	defer c.mut.Unlock()
	all := 0
	errors := 0
	m := map[string]interface{}{
		"r": roomName,
		"m": string(msg),		// I am so unhappy about this.
	}
	smsg, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	for i := range c.roomToConns[roomName] {
		if i == except.Id {
			continue
		}
		all++
		conn, exists := c.conns[i]
		if !exists {
			fmt.Println("Bookkeeping is buggy.")
			continue
		}
		err := binhelper.WriteMsg(conn.Conn, smsg)
		if err != nil {
			fmt.Println(err)
			errors++
		}
	}
	if errors > 0 {
		fmt.Println(errors, "errors while publishing to", all, "in room \"", roomName, "\"")
	}
}

func New() *Rooms {
	return &Rooms{
		map[[16]byte]*cn.Connection{},
		map[string]map[[16]byte]struct{}{},
		map[[16]byte]map[string]struct{}{},
		&sync.RWMutex{},
	}
}