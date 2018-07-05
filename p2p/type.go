package p2p

//UniMsgPacket it is a unicast packet , Index is peer id to send to
type UniMsgPacket struct {
	Index uint16
	P     Packet
}

//BcastMsgPacket it is a multicast packet , Indexs is filter peers index which not send to
type BcastMsgPacket struct {
	Indexs []uint16
	P      Packet
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
