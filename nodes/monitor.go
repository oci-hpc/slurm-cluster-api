package nodes

import (
	"fmt"
	"time"
)

func RunNodeMonitor() {
	fmt.Println("Running node monitor")
	for {
		fmt.Println("Running node monitor")
		nodes := getNodeStatus()
		for _, node := range nodes {
			upsertNodeStatus(node)
		}
		time.Sleep(15 * time.Second)
	}
}
