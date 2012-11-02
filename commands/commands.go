package commands

import(
	"net"
	"encoding/json"
	"github.com/opesun/trigga/binhelper"
	cn	"github.com/opesun/trigga/connection"
	ro "github.com/opesun/trigga/rooms"
	"fmt"
)

func readAndDecode(c net.Conn) (map[string]interface{}, error) {
	b, err := binhelper.ReadMsg(c)
	if err != nil {
		return nil, err
	}
	var v interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}
	ma, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Msg is not a map.")
	}
	return ma, nil
}

func Process(r *ro.Rooms, c *cn.Connection) {
	defer func() {
		recover()
		r.Remove(c)
	}()
	for {
		m, err := readAndDecode(c.Conn)
		if err != nil {
			panic(err)
		}
		switch m["c"] {
		case "p":
			r.Publish(m["r"].(string), []byte(m["m"].(string)))
		case "s":
			r.Subscribe(m["r"].(string), c)
		case "u":
			r.Unsubscribe(m["r"].(string), c)
		default:
			panic("Unkown command" + m["c"].(string))
		}
	}
}