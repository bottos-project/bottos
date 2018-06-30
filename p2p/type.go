package p2p

//Packet if it is a unicast packet , ids is peer id to send to , and if
// it is a boradcast packet, ids is filter peers id
type MsgPacket struct {
	Index []uint16
	P     Packet
}

const (
	MAX_PACKET_LEN = 10000000
)

type Head struct {
	ProtocolType uint16
	PacketType   uint16
	Pad          uint16
}

type Packet struct {
	H    Head
	Data []byte
}
