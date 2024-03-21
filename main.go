package main

import (
	"fmt"
	"time"

	"github.com/ipoluianov/cc_node/logger"
	"github.com/ipoluianov/cc_node/node"
)

var nodes []*node.Node

func runNode(id string) {
	n := node.NewNode()
	err := n.Start(logger.CurrentExePath() + "/data/" + id)
	if err != nil {
		fmt.Println("Run Node Error:", err)
	}
	nodes = append(nodes, n)
}

func main() {
	logger.InitNearExe()
	runNode("01")
	runNode("02")
	runNode("03")
	fmt.Println("Started")
	time.Sleep(1000 * time.Hour)
}
