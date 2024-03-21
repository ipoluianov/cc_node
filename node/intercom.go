package node

import (
	"errors"
	"fmt"
	"sync"
)

type Frame struct {
	SrcAddress  string
	DestAddress string
	Type        string
	Content     string
	Content1    string
	Content2    string
	Content3    string
}

func NewFrame(src string, dest string, tp string, context string) Frame {
	var c Frame
	c.SrcAddress = src
	c.DestAddress = dest
	c.Type = tp
	c.Content = context
	return c
}

type Intercom struct {
	mtx    sync.Mutex
	routes map[string]*Node
}

var ErrHostNotFound error
var com *Intercom

func init() {
	ErrHostNotFound = errors.New("ErrHostNotFound")

	var c Intercom
	c.init()
	com = &c
}

func (c *Intercom) init() {
	c.routes = make(map[string]*Node)
}

func (c *Intercom) RegisterReceiver(address string, node *Node) {
	c.mtx.Lock()
	c.routes[address] = node
	c.mtx.Unlock()
	fmt.Println("RegisterReceiver", address)
}

func (c *Intercom) Send(router string, frame Frame) (err error) {
	c.mtx.Lock()
	n, ok := c.routes[frame.DestAddress]
	c.mtx.Unlock()

	if ok && n != nil {
		n.Receive(frame)
	} else {
		err = ErrHostNotFound
	}
	return
}
