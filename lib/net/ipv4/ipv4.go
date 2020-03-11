package ipv4

import (
	"fmt"
	"github.com/google/gopacket"
)

const (
	TCPLen = 10000
	UDPLen = 10000
)

var TCPC = make(chan gopacket.Layer, TCPLen)	// TCP缓存队列
var UDPC = make(chan gopacket.Layer, UDPLen)	// UDP缓存队列

type IPv4 [4]byte

type IPv4Data struct {
	header Header
	data ipv4Data
}

type Header struct {
	Version 	byte	// 版本
	HeadLength 	byte	// 首部长度
	TOS 		byte	// 服务类型(TOS)
	TotalLength		uint16	// 总长度
	Identification	uint16	//
}

type ipv4Data string

// IPv4 网络层，封包操作
func IPDataPacket() {

}

// Ipv4 网络层，解包操作
func IPDataUnPacket() {

}

func Work() {
	switch p.Layers()[2].LayerType().String() {
	case "TCP":

	case "UDP":
		UDPC <- p.Layers()[2]
	default:
		continue
	}
	if p.Layers()[2].LayerType().String() != "UDP" {
		continue
	}


	udpP := p.Layers()[2]
	// 源端口   目的端口
	// 数据长度 校验码
	fmt.Println("-------------------------")
	hD := []byte(udpP.LayerContents())
	fmt.Printf("src port: %d, des port: %d\n", byte2Dec([2]byte{hD[0], hD[1]}), byte2Dec([2]byte{hD[2], hD[3]}))
	fmt.Printf("data len: %d, checksum: %d\n", byte2Dec([2]byte{hD[4], hD[5]}), byte2Dec([2]byte{hD[6], hD[7]}))
	fmt.Println(udpP.LayerPayload())
	fmt.Println(len(udpP.LayerPayload()), byte2Dec([2]byte{hD[4], hD[5]}))
	fmt.Println("-------------------------")
}
