// https://github.com/bwmarrin/snowflake
package snowflake

import (
	"strconv"
	"sync"
	"time"
)

type Node struct {
	mu       sync.Mutex
	lastTime int64 //上次时间
	machine  int64 //机器id
	step     int64 //序列号
	epoch    int64 //开始时间
	nodeBits uint8
	nodeMax  int64
	stepBits uint8
	stepMax  int64
}
type ID struct {
	id       int64
	epoch    int64
	nodeBits uint8
	stepBits uint8
}

func NewNode(epoch, machine int64, nodeBits, stepBits uint8) Node {
	if nodeBits+stepBits != 22 {
		panic("nodeBits加stepBits不等于22位")
	}
	if machine > 1<<nodeBits-1 {
		panic("超出最大机器码")
	}
	return Node{
		epoch:    epoch,
		machine:  machine,
		nodeBits: nodeBits,
		stepBits: stepBits,
		nodeMax:  1<<nodeBits - 1,
		stepMax:  1<<stepBits - 1,
	}
}

func (n *Node) Generate() ID {
	n.mu.Lock()
	now := time.Now().UnixNano()/1e6 - n.epoch
	if now == n.lastTime {
		n.step = n.step + 1&n.stepMax
		if n.step == 0 {
			for now <= n.lastTime {
				now = time.Now().UnixNano()/1e6 - n.epoch
			}
		}
	} else if now < n.lastTime {
		//时间回拨持续等待
		for now <= n.lastTime {
			now = time.Now().UnixNano()/1e6 - n.epoch
		}
		//如果第一次则不增加step
		if n.lastTime != 0 {
			n.step = n.step + 1&n.stepMax
		}
		if n.step == 0 {
			for now <= n.lastTime {
				now = time.Now().UnixNano()/1e6 - n.epoch
			}
		}
	} else {
		n.step = 0
	}
	n.lastTime = now
	r := now<<(n.nodeBits+n.stepBits) | (n.machine << n.stepBits) | (n.step)
	n.mu.Unlock()
	return ID{id: r, epoch: n.epoch, nodeBits: n.nodeBits, stepBits: n.stepBits}
}

func (f ID) Time() int64 {
	return (f.id >> (f.nodeBits + f.stepBits)) + f.epoch
}
func (f ID) Node() int64 {
	return f.id & ((1<<f.nodeBits - 1) << f.stepBits) >> f.stepBits
}
func (f ID) Step() int64 {
	return f.id & (1<<f.stepBits - 1)
}
func (f ID) Int64() int64 {
	return f.id
}
func (f ID) String() string {
	return strconv.FormatInt(f.id, 10)
}
