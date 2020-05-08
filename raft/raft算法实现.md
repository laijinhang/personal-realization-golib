```golang
package main

// 基于TCP
import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type Node struct {
	mutex       sync.Mutex // 锁
	nodeNo      int        // 节点编号
	currentTerm int        // 当前任期

	votedFor int // 为哪个节点投票
	votes    int // 票数

	role int // 当前节点角色，-1 无任何角色， 0 follower、1 candidate、2 leader

	currentLeader int // 当前节点的领导

	ip          string
	port        string
	candidateCh chan net.Conn // 候选人通信通道
	leaderCh    chan net.Conn // 领导者通信通道

	work           chan func() // 当前调度任务

	nodes                []Node         // 其他节点
	workRun             sync.WaitGroup 	// 当前工作
	currentHeartBeatTime int            // 当前心跳时间
	maxHeartBeatTime     int            // 最大心跳时间
	timeout              int            // 超时时间
}

type Message struct {
	NodeNo  int
	Term    int
	State   int // 0 跟随者，1 候选人，2 领导人
	Message interface{}
}

func main() {
	var NodeNo, _ = strconv.Atoi(os.Args[1])
	// 创建一个节点
	node := &Node{
		mutex:          sync.Mutex{},
		nodeNo:         NodeNo,
		currentTerm:    0,
		votedFor:       -1,
		role:           -1,
		currentLeader:  -1,
		ip:             "127.0.0.1",
		port:           "1000" + strconv.Itoa(NodeNo),
		nodes:          nil,
		candidateCh:    make(chan net.Conn, 20),
		leaderCh:       make(chan net.Conn, 20),
		work: 			make(chan func()),
		timeout:        1000, // 1000毫秒
	}
	fmt.Printf("节点: %d, 网络: %s:%s\n", node.nodeNo, node.ip, node.port)

	// 其他节点信息
	for i := 1; i <= 5; i++ {
		if i == NodeNo {
			continue
		}
		node.nodes = append(node.nodes, Node{
			mutex:         sync.Mutex{},
			nodeNo:        i,
			currentTerm:   0,
			votedFor:      -1,
			role:          -1,
			currentLeader: -1,
			ip:            "127.0.0.1",
			port:          "1000" + strconv.Itoa(i),
			nodes:         nil,
			timeout:       1000,
		})
	}

	go node.TranslateMessage() // 消息分发
	//time.Sleep(10 * time.Second)
	// 开始工作
	go node.TaskSchedule()
	node.Run()
}

// 尝试成为候选人
func (n *Node) CampaignCandidate() {
	n.workRun.Add(1)
	defer n.workRun.Done()

	fmt.Println("尝试获取候选人资格...")
	// 在 1000 ~ 3000 毫秒内，没有收到其它节点要成为候选人的消息时，成为候选人
	t := 1000 + rand.Int31n(2000)
	for i := 0;i < 100;i++ {
		time.Sleep(time.Duration(t / 100) * time.Millisecond)
		if n.role != -1 {
			return
		}
	}
	n.mutex.Lock()
	if n.votedFor == -1 {
		n.role = 1 // 成为候选人
		fmt.Println("节点: ", n.nodeNo, "成为候选人")
		n.votedFor = n.nodeNo // 为自己投一票
		n.votes++
	}
	n.mutex.Unlock()
}

// 成为候选人
func (n *Node) BecomeCandidate() {
	n.workRun.Add(1)
	defer n.workRun.Done()
}

// 尝试成为领导人
func (n *Node) CampaignLeader() {
	n.workRun.Add(1)
	defer n.workRun.Done()

	fmt.Println("CampaignLeader...")
	// 广播，告诉所有节点，我要成为领导者
	if n.votedFor == n.nodeNo {
		var wg sync.WaitGroup
		for i := 0; i < len(n.nodes); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				conn, err := net.Dial("tcp", n.nodes[i].ip+":"+n.nodes[i].port)
				if err != nil {
					return
				}
				defer conn.Close()
				buf, _ := json.Marshal(Message{
					NodeNo: n.nodeNo,
					State:  1,
				})
				conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
				_, err = conn.Write(buf)
				if err != nil {
					return
				}
				conn.SetReadDeadline(time.Now().Add(1 * time.Second))
				num, err := conn.Read(buf)
				if err != nil {
					return
				}
				var msg Message
				err = json.Unmarshal(buf[:num], &msg)
				if msg.State == 0 {
					fmt.Println(msg.Message)
				} else {
					fmt.Println(msg)
					return
				}
				n.mutex.Lock()
				n.votes++
				n.mutex.Unlock()
			}(i)
		}
		wg.Wait()
		// 成为领导者
		n.mutex.Lock()
		if n.votes > (len(n.nodes)+1)/2 {
			n.role = 2
			n.currentTerm++
			fmt.Println("节点: ", n.nodeNo, "成为leader")
		} else {
			n.role = -1
			n.votedFor = -1
			n.votes = 0
		}
		n.mutex.Unlock()
	}
}

// 选举领导人
func (n *Node) ElectionLeader() {
	n.workRun.Add(1)
	defer n.workRun.Done()


	// 1、选举领导者
	/*
		   {
			   "Node": 1,
			   "State": 1,    // 我是候选人，我现在要成为领导者，投我一票
		   }
	*/
	fmt.Println("ElectionLeader...")
	if len(n.candidateCh) == 0 {
		return
	}
	conn := <-n.candidateCh
	// 只投一票，其他直接抛弃
	for i := 0; i < len(n.candidateCh); i++ {
		t := <-n.candidateCh
		t.Close()
	}
	// 告诉第一个
	buf, _ := json.Marshal(Message{
		NodeNo: n.nodeNo,
		State:  0,
	})
	msg := Message{
		NodeNo:  n.nodeNo,
		State:   0,
		Message: "我是节点：" + strconv.Itoa(n.nodeNo) + ",我投你一票",
	}
	buf, _ = json.Marshal(msg)
	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	conn.Write(buf)
	time.Sleep(1 * time.Second)
	conn.Close()
}

// 成为领导人
func (n *Node) BecomeLeader() {
	n.workRun.Add(1)
	defer n.workRun.Done()

	fmt.Println("BecomeLeader...")
	n.mutex.Lock()
	n.currentLeader = n.nodeNo
	n.role = 2
	n.currentTerm++
	n.votedFor = -1
	n.votes = 0
	n.mutex.Unlock()
}

// 领导者工作
func (n *Node) LeaderWork() {
	n.workRun.Add(1)
	defer n.workRun.Done()

	var wg sync.WaitGroup
	if n.role == 2 {
		fmt.Println("LeaderWork...")
		for i := 0; i < len(n.nodes); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				conn, err := net.Dial("tcp", n.nodes[i].ip+":"+n.nodes[i].port)
				if err != nil {
					return
				}
				defer conn.Close()
				var msgW, msgR Message
				msgW = Message{
					NodeNo:  n.nodeNo,
					Term:    n.currentTerm,
					State:   2,
					Message: "我是节点: " + strconv.Itoa(n.nodeNo) + "，我是领导人",
				}
				bufW, _ := json.Marshal(&msgW)
				num := 0
				for num < 20 {
					conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
					if _, err := conn.Write(bufW); err != nil {
						num++
					}
					conn.SetReadDeadline(time.Now().Add(5 * time.Second))
					bufR := make([]byte, 512)
					numL, err := conn.Read(bufR)
					if err != nil {
						num++
						continue
					}
					err = json.Unmarshal(bufR[:numL], &msgR)
					fmt.Println("跟随者的消息：", msgR)
				}
			}(i)
		}
		wg.Wait() // 该节点失去领导人职位
		n.currentLeader = -1
		n.role = -1
	}
}

// 跟随者任务
func (n *Node) Work() {
	n.workRun.Add(1)
	defer n.workRun.Done()

	// 1、领导者消息
	/*
		   {
			   “Node”: 1,
			   "State": 2,            // 我是领导者
			   "Term": 3,            // 我是第三任领导者
			   "Message": "...",    // 消息
		   }
	*/
	conn := <-n.leaderCh
	fmt.Println("Work...")
	for i := 0; i < len(n.leaderCh); i++ {
		t := <-n.leaderCh
		t.Close()
	}
	if n.currentLeader == n.nodeNo {
		return
	}
	num := 0
	bufR := make([]byte, 1024)
	bufW := make([]byte, 1024)
	var msgW, msgR Message
	for num < 10 { // 与领导人失联
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		num, err := conn.Read(bufR)
		if err != nil {
			num++
			continue
		} else {
			json.Unmarshal(bufR[:num], &msgR)
			n.currentLeader = msgR.NodeNo
			n.currentTerm = msgR.Term
		}
		fmt.Println("来自领导者的消息：", msgR)
		msgW = Message{
			NodeNo:  n.nodeNo,
			State:   0,
			Message: "领导者：" + strconv.Itoa(n.currentLeader) + ",我是节点: " + strconv.Itoa(n.nodeNo) + "，我收到你的消息了",
		}
		bufW, _ = json.Marshal(&msgW)
		time.Sleep( 1 * time.Second)
		num, err = conn.Write(bufW)
		if err != nil {
			num++
		}
	}
	conn.Close()
}

// 任务调度
func (n *Node) TaskSchedule() {
	role := 0
	for {
		n.mutex.Lock()
		role = n.role
		n.mutex.Unlock()
		switch role {
		case -1:
			n.work <- n.CampaignCandidate // 尝试成为候选人
			n.work <- n.ElectionLeader    // 为其他节点投票
		case 0:
			//n.work <- n.Work
		case 1:
			n.work <- n.BecomeCandidate // 成为候选人
			n.work <- n.CampaignLeader  // 尝试成为领导者
		case 2:
			n.work <- n.BecomeLeader // 成为领导人
			n.work <- n.LeaderWork   // 领导者工作
		}
		if len(n.leaderCh) != 0 {
			n.work <- n.Work
		}
		n.workRun.Wait()
	}
}

// 消息分发
func (n *Node) TranslateMessage() {
	listen, err := net.Listen("tcp", ":"+n.port)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		// 消息分发
		go func(conn net.Conn) {
			nowStatus := -1 // 0 表示当前没有领导者，1表示当前有领导者

			buf := make([]byte, 1024)
			msg := Message{}
			num, err := conn.Read(buf)
			if err != nil {
				conn.Close()
				return
			}
			json.Unmarshal(buf[:num], &msg)

			// 0 跟随者，1 候选人，2 领导人
			switch msg.State {
			case 1:
				if n.role == -1 {
					nowStatus = 0
					n.currentLeader = -1
					n.votedFor = msg.NodeNo
				}
			case 2:
				nowStatus = 1
				n.role = 0
				n.currentLeader = msg.NodeNo
				n.currentTerm = msg.Term
			}

			switch nowStatus {
			case 0:
				n.candidateCh <- conn // 投票
			case 1:
				n.leaderCh <- conn // 领导者连接
			default:
				// 当前节点与那个节点，竞争成为领导者
				conn.Close()
			}
		}(conn)
	}
}

func (n *Node)Run() {
	for {
		work := <- n.work
		go work()
	}
}
```
