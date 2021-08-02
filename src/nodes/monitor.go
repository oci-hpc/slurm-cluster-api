package nodes

import "time"

func RunNodeMonitor() {
	for {
		nodes := getNodeStatus()
		for _, node := range nodes {
			upsertNodeStatus(node)
		}
		time.Sleep(15 * time.Second)
	}
}
