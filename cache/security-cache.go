package cache
// 并发安全的cache实现

import "sync"

type LRUCache struct {
	head, tail    *LinkedList         // 头结点、尾节点
	mapKey        map[int]*LinkedList // key => 节点
	mapNode       map[*LinkedList]int // 节点 => key
	len, capacity int                 // 已使用长度、总长度
	lock sync.Locker
}

type LinkedList struct {
	prior, next *LinkedList // 前驱节点、后续节点
	val         int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		mapNode:  make(map[*LinkedList]int),
		mapKey:   make(map[int]*LinkedList),
		capacity: capacity,
	}
}

func (this *LRUCache) Get(key int) int {
	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.mapKey[key]; ok {
		val := this.mapKey[key].val

		// 把key对应的节点 移到头部
		switch this.mapKey[key] {
		case this.head: // 头部
		case this.tail: // 尾部
			this.tailMoveToHead()
		default:
			this.moveToHead(this.mapKey[key])
		}

		return val
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	
	if this.capacity <= 0 {
		return
	}
	// 如果存在
	if _, ok := this.mapKey[key]; ok {
		// 置换新值
		this.mapKey[key].val = value
		// 把key对应的节点 移到头部
		switch this.mapKey[key] {
		case this.head: // 头部
		case this.tail: // 尾部
			this.tailMoveToHead()
		default:
			this.moveToHead(this.mapKey[key])
		}
		return
	}
	// 如果不存在
	if this.len == this.capacity { // 如果上限了
		if this.capacity == 1 { // 将第一个值替换掉
			// 淘汰末尾的key
			delete(this.mapKey, this.mapNode[this.head])
			this.head.val = value
		} else { // 将未节点移到头部
			// 淘汰末尾的key
			delete(this.mapKey, this.mapNode[this.tail])
			this.tailMoveToHead()
			this.head.val = value
		}
		this.mapKey[key] = this.head
		this.mapNode[this.head] = key
		return
	} else if this.len != this.capacity { // 没有达到上限
		t := &LinkedList{
			prior: nil,
			next:  this.head,
			val:   value,
		}
		if this.len >= 1 {
			this.head.prior = t
		}
		this.head = t
		this.len++
		if this.len == 2 {
			this.tail = this.head.next
			this.tail.prior = this.head
		}
		this.mapNode[t] = key
		this.mapKey[key] = t
	}
}

// 其他位置移到头部
func (this *LRUCache) moveToHead(n *LinkedList) {
	// 1、将要移除的节点从原链表中移出
	n.next.prior = n.prior
	n.prior.next = n.next

	n.prior = nil
	n.next = nil
	// 2、将n移到头部
	n.next = this.head
	this.head.prior = n
	// 3、head指向新的第一个节点
	this.head = n
}

// 末尾移到头部
func (this *LRUCache) tailMoveToHead() {
	// 1、保留倒数第二个节点
	t := this.tail.prior
	// 2、把倒数第二个节点变成最后一个节点
	t.next = nil
	// 3、把原最后一个位置变成第一个节点
	this.tail.prior = nil // 前驱指向nil
	this.tail.next = this.head
	// 4、把原第一个变成第二个节点
	this.head.prior = this.tail
	// 5、head指向新的第一个节点
	this.head = this.tail
	// 6、tail指向新的最后一个节点
	this.tail = t
}
