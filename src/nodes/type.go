package nodes

import "time"

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
	Id            int
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
	LastSeenTime time.Time
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
