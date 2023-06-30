/*

	FILE CONTAINS IMPLEMENTATION FOR WEBSOCKETS for ws://{name}/directs

*/

package handlers 

import (
	"net"

	"log"
	"encoding/json"

	"bufio"
)

// DirectsReceive : format for receiving data from websocket
// type DirectsReceive struct {
// 	ReceiverUsername 		string 		`json:"receiver"` 		//username of the receiver
// 	Content 				string 		`json:"content"`		//content of the message as text
// }
// // DirectsSend : format for sending data from websocket
// type DirectsSend struct {
// 	Directs 				RenderedDirects					`json:"directs"` 		// all directs to be received
// 	Status 					int 							`json:"status"`			// status of operation
// }

// Channel wraps user connection.
type Channel struct {
    conn net.Conn    // WebSocket connection.
    send chan DirectsSend // Outgoing packets queue.
}

func NewChannel(conn net.Conn) *Channel {
    c := &Channel{
        conn: conn,
        send: make(chan DirectsSend, 1024),
    }

    go c.reader()
    go c.writer()

    return c
}

func readPacket(buf *bufio.Reader) (DirectsReceive , error){
	var raw []byte
	_ , err := buf.Read(raw)
	if err != nil {
		return DirectsReceive{} , err
	}	

	var data DirectsReceive
	err = json.Unmarshal(raw , &data)
	if err != nil {	
		return DirectsReceive{} , err
	}

	return data , nil
}
func writePacket(buf *bufio.Writer , packet DirectsSend) error{
	raw , err := json.Marshal(packet)
	if err != nil {
		return err
	}

	_ , err = buf.Write(raw)
	if err != nil {
		return err
	}

	return nil
}

func (c *Channel) handle(packet DirectsReceive){
	log.Println("Packet = ", packet)
}

func (c *Channel) reader() {
    // We make a buffered read to reduce read syscalls.
    buf := bufio.NewReader(c.conn)

    for {
        pkt, _ := readPacket(buf)
        c.handle(pkt)
    }
}


func (c *Channel) writer() {
    // We make buffered write to reduce write syscalls. 
    buf := bufio.NewWriter(c.conn)

    for pkt := range c.send {
        _ = writePacket(buf, pkt)
        buf.Flush()
    }
}