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
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func InitializeNodesEndpoint(r *gin.Engine) {
	r.GET("/nodes", getNodes)
	//r.POST("/nodes/update", updateNode)
}

func getNodes(cnx *gin.Context) {
	nodes := queryAllNodeStatus()
	cnx.JSON(200, nodes)
}

func updateNode(cnx *gin.Context) {
	jsonData, err := ioutil.ReadAll(cnx.Request.Body)

	if err != nil {
		cnx.JSON(400, err.Error())
		return
	}

	var updateNodeMsg UpdateNodeMsg
	json.Unmarshal([]byte(jsonData), &updateNodeMsg)

	var slreq C.update_node_msg_t
	C.slurm_init_update_node_msg(&slreq)
	slreq = convertUpdateNodeMsg(updateNodeMsg)
	ret := C.slurm_update_node(&slreq)
	print(ret)

}
