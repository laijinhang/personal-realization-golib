package net

import (
	"lib/net/datalink"
	"lib/net/ipv4"
	"lib/net/transport/udp"
)


// 模拟os调度
func init() {
	// 数据链路层工作
	go datalink.Work()
	// 网络层工作
	go ipv4.Work()
	// 传输层工作
	go udp.Work()
}
