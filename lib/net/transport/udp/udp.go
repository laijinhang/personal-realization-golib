package udp

import "lib/net/ipv4"

const ProtocolVersion = 17

type UDPPseudoHead struct {
	SourceIP 		ipv4.IPv4
	DestinationIP 	ipv4.IPv4
	Reserved 		int8
	ProtoclVersion 	int8
	Length 			int16
}

type header struct {
	SourcePort uint16
	DestinationPort uint16
	Length uint16
	CheckSum uint16
}

type UDP struct {
	PseudoHead UDPPseudoHead	// 伪头部
	Header header
	Data []byte
}

// UDP 传输层，封包操作
func (this *UDP) Packet() []byte {
	buf := make([]byte, 0, len(this.Data)+8)
	buf = append(buf, byte(this.Header.SourcePort>>8), byte((this.Header.SourcePort << 8) >> 8))
	buf = append(buf, byte(this.Header.DestinationPort>>8), byte((this.Header.DestinationPort << 8) >> 8))
	buf = append(buf, byte(this.Header.Length>>8), byte((this.Header.Length << 8) >> 8))
	buf = append(buf, byte(this.Header.CheckSum>>8), byte((this.Header.CheckSum << 8) >> 8))

	buf = append(buf, this.Data...)
	return buf
}
// UDP 传输层，解包操作
func (this *UDP) UnPacket(d []byte) UDP {
	udpP := UDP{}
	udpP.Header = header{
		SourcePort:      byte2Dec([2]byte{d[0], d[1]}),
		DestinationPort: byte2Dec([2]byte{d[2], d[3]}),
		Length:          byte2Dec([2]byte{d[4], d[5]}),
		CheckSum:        byte2Dec([2]byte{d[6], d[7]}),
	}
	udpP.Data = d[8:]
	return udpP
}

// 计算出校验和
func (this *UDP)calculationChecksums() uint16 {
	checksums := int32(0)
	// 1、计算UDP伪头部
	checksums += int32(byte2Dec([2]byte{this.PseudoHead.SourceIP[0], this.PseudoHead.SourceIP[1]}))
	checksums += int32(byte2Dec([2]byte{this.PseudoHead.SourceIP[2], this.PseudoHead.SourceIP[3]}))
	checksums = (checksums << 16) >> 16 + checksums >> 16

	checksums += int32(byte2Dec([2]byte{this.PseudoHead.DestinationIP[0], this.PseudoHead.DestinationIP[1]}))
	checksums = (checksums << 16) >> 16 + checksums >> 16
	checksums += int32(byte2Dec([2]byte{this.PseudoHead.DestinationIP[2], this.PseudoHead.DestinationIP[3]}))
	checksums = (checksums << 16) >> 16 + checksums >> 16

	checksums += int32(byte2Dec([2]byte{byte(this.PseudoHead.Reserved), byte(this.PseudoHead.ProtoclVersion)}))
	checksums = (checksums << 16) >> 16 + checksums >> 16

	checksums += int32(this.PseudoHead.Length)
	checksums = (checksums << 16) >> 16 + checksums >> 16
	// 2、计算UDP头部
	checksums += int32(this.Header.SourcePort)
	checksums = (checksums << 16) >> 16 + checksums >> 16
	checksums += int32(this.Header.DestinationPort)
	checksums = (checksums << 16) >> 16 + checksums >> 16
	checksums += int32(this.Header.Length)
	checksums = (checksums << 16) >> 16 + checksums >> 16
	// 3、计算UDP数据部分
	for i := 0;i <= len(this.Data)-2;i += 2 {
		checksums += int32(byte2Dec([2]byte{byte(this.Data[i]), byte(this.Data[i+1])}))
		checksums = (checksums << 16) >> 16 + checksums >> 16
	}
	if len(this.Data) % 2 != 0 {
		checksums += int32(byte2Dec([2]byte{byte(this.Data[len(this.Data)-1]), byte(0)}))
		checksums = (checksums << 16) >> 16 + checksums >> 16
	}

	return uint16(checksums)	// 截取低16位
}

func byte2Dec(b [2]byte) uint16 {
	return uint16(b[0]) << 8 + uint16(b[1])
}

func Work() {}