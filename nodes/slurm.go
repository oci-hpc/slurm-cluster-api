package nodes

/*
#cgo pkg-config: slurm
#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <stdlib.h>
#include <stdint.h>
#include "slurm/slurm.h"
#include "slurm/slurm_errno.h"
*/
import "C"
import (
	"log"
	"time"
)

func getNodeStatus() (nodes []NodeInfo) {
	var slres *C.node_info_msg_t
	defer C.slurm_free_node_info_msg(slres)

	now := time.Now().Unix()
	ret := C.slurm_load_node(C.long(now-1000), &slres, 0)
	if ret == -1 {
		log.Printf("WARN: getNodeStatus: unable to retrieve node status")
		return
	}
	count := int(slres.record_count)
	nodes = convertNodeInfoArray(slres.node_array, count)

	return nodes
}
