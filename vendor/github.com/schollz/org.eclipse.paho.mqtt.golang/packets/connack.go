package packets

import (
	"bytes"
	"fmt"
	"io"
)

//ConnackPacket is an internal representation of the fields of the
//Connack MQTT packet
type ConnackPacket struct {
	FixedHeader
	TopicNameCompression byte
	ReturnCode           byte
}

func (ca *ConnackPacket) String() string {
	str := fmt.Sprintf("%s\n", ca.FixedHeader)
	str += fmt.Sprintf("returncode: %d", ca.ReturnCode)
	return str
}

func (ca *ConnackPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.WriteByte(ca.TopicNameCompression)
	body.WriteByte(ca.ReturnCode)
	ca.FixedHeader.RemainingLength = 2
	packet := ca.FixedHeader.pack()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

//Unpack decodes the details of a ControlPacket after the fixed
//header has been read
func (ca *ConnackPacket) Unpack(b io.Reader) {
	ca.TopicNameCompression = decodeByte(b)
	ca.ReturnCode = decodeByte(b)
}

//Details returns a Details struct containing the Qos and
//MessageID of this ControlPacket
func (ca *ConnackPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0}
}
