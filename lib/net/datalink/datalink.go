package datalink

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"time"
)

const (
	ReadQueryLength = 10000 * 2
	SendQueryLength = 10000 * 2
)

var packReadQuery = make(chan gopacket.Layer, ReadQueryLength)
var packSendQuery = make(chan gopacket.Layer, SendQueryLength)

var (
	deviceName  string
	snapShotLen int32 = 1024
	promiscuous bool
	err         error
	timeout     = 30 * time.Second
	handle      pcap.Handle
	ports 		[65536]bool
)

func Work() {
	go Send()
	go Read()
}

func Send() {}

func Read() {
	ifs, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}
	deviceName = ifs[2].Name // 每台电脑都是不一样的，先打印全部，找到对应要监听的网卡，我的网卡是第三个

	handle, err := pcap.OpenLive(deviceName, snapShotLen, promiscuous, timeout)
	if err != nil {
		panic(err)
	}
	defer handle.Close()

	pSrc := gopacket.NewPacketSource(handle, handle.LinkType())
	for p := range pSrc.Packets() {
		// 数据链路    p.Layers()[0]
		// 网络层    p.Layers()[1]
		// 传输层    p.Layers()[2]
		fmt.Println(2)
		if len(p.Layers()) <= 1 {
			continue
		}
		packReadQuery <- p.Layers()[1]
	}
}

func ReadPack() gopacket.Layer {
	return <-packReadQuery
}

func byte2Dec(b [2]byte) int {
	return int(b[0])<<8 + int(b[1])
}

// 解帧，发送给网络层
func Send2NetWork() {

}
// 封帧，发送数据出去
func SendData() {

}

// 心有数据结构和算法，以设计技能为经验，以Net、OS支撑，SQL做数据存储，Redis，MQ缓存，Docker、K8s架构支撑项目，用Go表达与描绘