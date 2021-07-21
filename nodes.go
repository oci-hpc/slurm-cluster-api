package main

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
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func InitializeNodesEndpoint(r *gin.Engine) {
	r.GET("/nodes", getNodes)
	r.POST("/nodes/update", updateNode)
}

func getNodes(cnx *gin.Context) {
	var slres *C.node_info_msg_t
	now := time.Now().Unix()
	ret := C.slurm_load_node(C.long(now-1000), &slres, 0)
	count := int(slres.record_count)
	nodes := convertNodeInfoArray(slres.node_array, count)

	if ret == 0 {
		cnx.JSON(200, nodes)
	} else {
		cnx.JSON(500, gin.H{})
	}

	C.slurm_free_node_info_msg(slres)
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

func convertNodeInfoArray(node_array *C.node_info_t, count int) []NodeInfo {
	node_info_slice := (*[1 << 28]C.node_info_t)(unsafe.Pointer(node_array))[:count:count]
	var nodes []NodeInfo
	// Not sure why, but array is often large and filled with empty objects.
	// Iterate backwards and break on empty object to avoid unnecessary cycles
	// Checking arch seems to work, might not be the best field to check
	for i := count - 1; i >= 0; i-- {
		if node_info_slice[i].arch == nil {
			break
		}
		nodeInfo := convertNodeInfo(node_info_slice[i])
		nodes = append(nodes, nodeInfo)
	}
	return nodes
}

func convertNodeInfo(node_info C.node_info_t) NodeInfo {
	var nodeInfo NodeInfo
	nodeInfo.Arch = C.GoString(node_info.arch)
	nodeInfo.BcastAddress = C.GoString(node_info.bcast_address)
	nodeInfo.Boards = int(node_info.boards)
	nodeInfo.BootTime = time.Unix(0, int64(node_info.boot_time))
	nodeInfo.ClusterName = C.GoString(node_info.cluster_name)
	nodeInfo.Cores = int(node_info.cores)
	nodeInfo.CpuBind = int(node_info.cpu_bind)
	nodeInfo.CpuLoad = int(node_info.cpu_load)
	nodeInfo.FreeMem = int(node_info.free_mem)
	nodeInfo.Cpus = int(node_info.cpus)
	nodeInfo.CpuSpecList = C.GoString(node_info.cpu_spec_list)
	nodeInfo.Features = C.GoString(node_info.features)
	nodeInfo.FeaturesAct = C.GoString(node_info.features_act)
	nodeInfo.Gres = C.GoString(node_info.gres)
	nodeInfo.GresDrain = C.GoString(node_info.gres_drain)
	nodeInfo.GresUsed = C.GoString(node_info.gres_used)
	nodeInfo.McsLabel = C.GoString(node_info.mcs_label)
	nodeInfo.MemSpecLimit = int(node_info.mem_spec_limit)
	nodeInfo.Name = C.GoString(node_info.name)
	nodeInfo.NextState = int(node_info.next_state)
	nodeInfo.NodeAddr = C.GoString(node_info.node_addr)
	nodeInfo.NodeHostname = C.GoString(node_info.node_hostname)
	nodeInfo.NodeState = int(node_info.node_state)
	nodeInfo.OS = C.GoString(node_info.os)
	nodeInfo.Owner = int(node_info.owner)
	nodeInfo.Partitions = C.GoString(node_info.partitions)
	nodeInfo.Port = int(node_info.port)
	nodeInfo.RealMemory = int(node_info.real_memory)
	nodeInfo.Comment = C.GoString(node_info.comment)
	nodeInfo.Reason = C.GoString(node_info.reason)
	nodeInfo.ReasonTime = time.Unix(0, int64(node_info.reason_time))
	nodeInfo.ReasonUID = int(node_info.reason_uid)
	nodeInfo.SlurmdStartTime = time.Unix(0, int64(node_info.slurmd_start_time))
	nodeInfo.Sockets = int(node_info.sockets)
	nodeInfo.Threads = int(node_info.threads)
	nodeInfo.TmpDisk = int(node_info.tmp_disk)
	nodeInfo.Weight = int(node_info.weight)
	nodeInfo.TresFmtStr = C.GoString(node_info.tres_fmt_str)
	nodeInfo.Version = C.GoString(node_info.version)
	return nodeInfo
}

func convertUpdateNodeMsg(updateNodeMsg UpdateNodeMsg) C.update_node_msg_t {
	var slreq C.update_node_msg_t
	slreq.comment = C.CString(updateNodeMsg.Comment)
	slreq.cpu_bind = C.uint(updateNodeMsg.CpuBind)
	slreq.features = C.CString(updateNodeMsg.Features)
	slreq.features_act = C.CString(updateNodeMsg.FeaturesAct)
	slreq.gres = C.CString(updateNodeMsg.Gres)
	slreq.node_addr = C.CString(updateNodeMsg.NodeAddr)
	slreq.node_hostname = C.CString(updateNodeMsg.NodeHostname)
	slreq.node_names = C.CString(updateNodeMsg.NodeNames)
	slreq.reason = C.CString(updateNodeMsg.Reason)
	slreq.reason_uid = C.uint(updateNodeMsg.ReasonUID)
	slreq.weight = C.uint(updateNodeMsg.Weight)
	return slreq
}

type UpdateNodeMsg struct {
	Comment      string
	CpuBind      int
	Features     string
	FeaturesAct  string
	Gres         string
	NodeAddr     string
	NodeHostname string
	NodeNames    string
	NodeState    int
	Reason       string
	ReasonUID    int
	Weight       int
}

type NodeInfo struct {
	Arch          string
	BcastAddress  string
	Boards        int
	BootTime      time.Time
	ClusterName   string
	Cores         int
	CoreSpecCount int
	CpuBind       int
	CpuLoad       int
	FreeMem       int
	Cpus          int
	CpuSpecList   string
	//energy			*_Ctype_struct_acct_gather_energy
	//ext_sensors		*_Ctype_struct_ext_sensors_data
	//power			*_Ctype_struct_power_mgmt_data
	Features     string
	FeaturesAct  string
	Gres         string
	GresDrain    string
	GresUsed     string
	McsLabel     string
	MemSpecLimit int
	Name         string
	NextState    int
	NodeAddr     string
	NodeHostname string
	NodeState    int
	OS           string
	Owner        int
	Partitions   string
	Port         int
	RealMemory   int
	Comment      string
	Reason       string
	ReasonTime   time.Time
	ReasonUID    int
	//select_nodeinfo		*_Ctype_struct_dynamic_plugin_data
	SlurmdStartTime time.Time
	Sockets         int
	Threads         int
	TmpDisk         int
	Weight          int
	TresFmtStr      string
	Version         string
}
