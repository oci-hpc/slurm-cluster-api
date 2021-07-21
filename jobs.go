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
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func InitializeJobsEndpoint(r *gin.Engine) {
	r.GET("/jobs", getJobs)
}

func getJobs(cnx *gin.Context) {
	var slres *C.job_info_msg_t
	now := time.Now().Unix()
	ret := C.slurm_load_jobs(C.long(now-1000), &slres, 0)
	count := int(slres.record_count)
	nodes := convertJobInfoArray(slres.job_array, count)

	if ret == 0 {
		cnx.JSON(200, nodes)
	} else {
		cnx.JSON(500, gin.H{})
	}

	C.slurm_free_job_info_msg(slres)
}

func convertJobInfoArray(job_array *C.job_info_t, count int) []JobInfo {
	var jobInfoSlice []JobInfo
	if count == 0 {
		return jobInfoSlice
	}
	job_info_slice := (*[1 << 28]C.job_info_t)(unsafe.Pointer(job_array))[:count:count]
	for i := 0; i < count; i++ {
		jobInfo := convertJobInfo(job_info_slice[i])
		jobInfoSlice = append(jobInfoSlice, jobInfo)
	}
	return jobInfoSlice
}

func convertJobInfo(slres C.job_info_t) JobInfo {
	var jobInfo JobInfo
	jobInfo.Account = C.GoString(slres.account)
	jobInfo.AccrueTime = time.Unix(int64(slres.accrue_time), 0)
	jobInfo.AdminComment = C.GoString(slres.admin_comment)
	jobInfo.AllocNode = C.GoString(slres.alloc_node)
	jobInfo.AllocSid = int(slres.alloc_sid)
	jobInfo.ArrayJobId = int(jobInfo.ArrayJobId)
	jobInfo.ArrayTaskId = int(jobInfo.ArrayTaskId)
	jobInfo.ArrayMaxTasks = int(jobInfo.ArrayMaxTasks)
	jobInfo.ArrayTaskStr = C.GoString(slres.array_task_str)
	jobInfo.AssocId = int(jobInfo.AllocSid)
	jobInfo.BatchFeatures = C.GoString(slres.batch_features)
	jobInfo.BatchFlag = int(slres.batch_flag)
	jobInfo.BitFlags = int(slres.bitflags)
	jobInfo.BoardsPerNode = int(slres.boards_per_node)
	jobInfo.BurstBuffer = C.GoString(slres.burst_buffer)
	jobInfo.BurstBufferState = C.GoString(slres.burst_buffer_state)
	jobInfo.Cluster = C.GoString(slres.cluster)
	jobInfo.ClusterFeatures = C.GoString(slres.cluster_features)
	jobInfo.Command = C.GoString(slres.command)
	jobInfo.Comment = C.GoString(slres.comment)
	jobInfo.Contiguous = int(slres.contiguous)
	jobInfo.CoreSpec = int(slres.core_spec)
	jobInfo.CoresPerSocket = int(slres.cores_per_socket)
	jobInfo.BillableTres = int(slres.billable_tres)
	jobInfo.CpusPerTask = int(slres.cpus_per_task)
	jobInfo.CpuFreqMin = int(slres.cpu_freq_min)
	jobInfo.CpuFreqMax = int(slres.cpu_freq_max)
	jobInfo.CpuFreqGov = int(slres.cpu_freq_gov)
	jobInfo.CpusPerTres = C.GoString(slres.cpus_per_tres)
	jobInfo.Cronspec = C.GoString(slres.cronspec)
	jobInfo.Deadline = time.Unix(int64(slres.deadline), 0)
	jobInfo.DelayBoot = int(slres.delay_boot)
	jobInfo.Dependency = C.GoString(slres.dependency)
	jobInfo.DerivedEC = int(slres.derived_ec)
	jobInfo.EligibleTime = time.Unix(int64(slres.eligible_time), 0)
	jobInfo.EndTime = time.Unix(int64(slres.end_time), 0)
	jobInfo.ExcNodes = C.GoString(slres.exc_nodes)
	jobInfo.ExcNodeInx = int(*slres.exc_node_inx)
	jobInfo.ExitCode = int(slres.exit_code)
	jobInfo.Features = C.GoString(slres.features)
	jobInfo.FedOriginStr = C.GoString(slres.fed_origin_str)
	jobInfo.FedSiblingsActive = int64(slres.fed_siblings_active)
	jobInfo.FedSiblingsActiveStr = C.GoString(slres.fed_siblings_active_str)
	jobInfo.FedSiblingsViable = int64(slres.fed_siblings_viable)
	jobInfo.FedSibilingsViableStr = C.GoString(slres.fed_siblings_viable_str)
	jobInfo.GresDetailCnt = int(slres.gres_detail_cnt)
	//jobInfo.GresDetailStr = C.GoString(*slres.gres_detail_str)
	jobInfo.GresTotal = C.GoString(slres.gres_total)
	jobInfo.GroupId = int(slres.group_id)
	jobInfo.HetJobId = int(slres.het_job_id)
	jobInfo.HetJobIdSet = C.GoString(slres.het_job_id_set)
	jobInfo.HetJobOffset = int(slres.het_job_offset)
	jobInfo.JobId = int(slres.job_id)
	jobInfo.JobState = int(slres.job_state)
	jobInfo.LastSchedEval = time.Unix(int64(slres.last_sched_eval), 0)
	jobInfo.Licenses = C.GoString(slres.licenses)
	jobInfo.MailType = int(slres.mail_type)
	jobInfo.MailUser = C.GoString(slres.mail_user)
	jobInfo.MaxCpus = int(slres.max_cpus)
	jobInfo.MaxNodes = int(slres.max_nodes)
	jobInfo.McsLabel = C.GoString(slres.mcs_label)
	jobInfo.MemPerTres = C.GoString(slres.mem_per_tres)
	jobInfo.Name = C.GoString(slres.name)
	jobInfo.Network = C.GoString(slres.network)
	jobInfo.Nodes = C.GoString(slres.nodes)
	jobInfo.Nice = int(slres.nice)
	jobInfo.NodeInx = int(*slres.node_inx)
	jobInfo.NTasksPerCore = int(slres.ntasks_per_core)
	jobInfo.NTasksPerTres = int(slres.ntasks_per_tres)
	jobInfo.NTasksPerNode = int(slres.ntasks_per_node)
	jobInfo.NTasksPerSocket = int(slres.ntasks_per_socket)
	jobInfo.NTasksPerBoard = int(slres.ntasks_per_board)
	jobInfo.NumCpus = int(slres.num_cpus)
	jobInfo.NumNodes = int(slres.num_nodes)
	jobInfo.NumTasks = int(slres.num_tasks)
	jobInfo.Partition = C.GoString(slres.partition)
	jobInfo.PnMinMemory = int64(slres.pn_min_memory)
	jobInfo.PnMinCpus = int(slres.pn_min_cpus)
	jobInfo.PnMinTmpDisk = int(slres.pn_min_tmp_disk)
	jobInfo.PowerFlags = int(slres.power_flags)
	jobInfo.PreemptTime = time.Unix(int64(slres.preempt_time), 0)
	jobInfo.PreSusTime = time.Unix(int64(slres.pre_sus_time), 0)
	jobInfo.Priority = int(slres.priority)
	jobInfo.Profile = int(slres.profile)
	jobInfo.Qos = C.GoString(slres.qos)
	jobInfo.Reboot = int(slres.reboot)
	jobInfo.ReqNodes = C.GoString(slres.req_nodes)
	jobInfo.ReqNodeInx = int(*slres.req_node_inx)
	jobInfo.ReqSwitch = int(slres.req_switch)
	jobInfo.Requeue = int(slres.requeue)
	jobInfo.ResizeTime = time.Unix(int64(slres.resize_time), 0)
	jobInfo.RestartCnt = int(slres.restart_cnt)
	jobInfo.ResvName = C.GoString(slres.resv_name)
	jobInfo.SchedNodes = C.GoString(slres.sched_nodes)
	jobInfo.Shared = int(slres.shared)
	jobInfo.ShowFlags = int(slres.show_flags)
	jobInfo.SiteFactor = int(slres.site_factor)
	jobInfo.SocketsPerBoard = int(slres.sockets_per_board)
	jobInfo.SocketsPerNode = int(slres.sockets_per_node)
	jobInfo.StartTime = time.Unix(int64(slres.start_time), 0)
	jobInfo.StartProtocolVer = int(slres.start_protocol_ver)
	jobInfo.StateDesc = C.GoString(slres.state_desc)
	jobInfo.StateReason = int(slres.state_reason)
	jobInfo.StdErr = C.GoString(slres.std_err)
	jobInfo.StdIn = C.GoString(slres.std_in)
	jobInfo.StdOut = C.GoString(slres.std_out)
	jobInfo.SubmitTime = time.Unix(int64(slres.submit_time), 0)
	jobInfo.SuspendTime = time.Unix(int64(slres.suspend_time), 0)
	jobInfo.SystemComment = C.GoString(slres.system_comment)
	jobInfo.TimeLimit = int(slres.time_limit)
	jobInfo.TimeMin = int(slres.time_min)
	jobInfo.ThreadsPerCore = int(slres.threads_per_core)
	jobInfo.Tresbind = C.GoString(slres.tres_bind)
	jobInfo.TresFreq = C.GoString(slres.tres_freq)
	jobInfo.TresPerJob = C.GoString(slres.tres_per_job)
	jobInfo.TresPerNode = C.GoString(slres.tres_per_node)
	jobInfo.TresPerSocket = C.GoString(slres.tres_per_socket)
	jobInfo.TresPerTask = C.GoString(slres.tres_per_task)
	jobInfo.TresReqStr = C.GoString(slres.tres_req_str)
	jobInfo.TresAllocStr = C.GoString(slres.tres_alloc_str)
	jobInfo.UserId = int(slres.user_id)
	jobInfo.UserName = C.GoString(slres.user_name)
	jobInfo.Wait4Switch = int(slres.wait4switch)
	jobInfo.WcKey = C.GoString(slres.wckey)
	jobInfo.WorkDir = C.GoString(slres.work_dir)
	return jobInfo
}

type JobInfo struct {
	Account      string
	AccrueTime   time.Time
	AdminComment string
	AllocNode    string
	AllocSid     int
	//ArrayBitmap unsafe.Pointer
	ArrayJobId            int
	ArrayTaskId           int
	ArrayMaxTasks         int
	ArrayTaskStr          string
	AssocId               int
	BatchFeatures         string
	BatchFlag             int
	BatchHost             string
	BitFlags              int
	BoardsPerNode         int
	BurstBuffer           string
	BurstBufferState      string
	Cluster               string
	ClusterFeatures       string
	Command               string
	Comment               string
	Contiguous            int
	CoreSpec              int
	CoresPerSocket        int
	BillableTres          int
	CpusPerTask           int
	CpuFreqMin            int
	CpuFreqMax            int
	CpuFreqGov            int
	CpusPerTres           string
	Cronspec              string
	Deadline              time.Time
	DelayBoot             int
	Dependency            string
	DerivedEC             int
	EligibleTime          time.Time
	EndTime               time.Time
	ExcNodes              string
	ExcNodeInx            int
	ExitCode              int
	Features              string
	FedOriginStr          string
	FedSiblingsActive     int64
	FedSiblingsActiveStr  string
	FedSiblingsViable     int64
	FedSibilingsViableStr string
	GresDetailCnt         int
	GresDetailStr         string
	GresTotal             string
	GroupId               int
	HetJobId              int
	HetJobIdSet           string
	HetJobOffset          int
	JobId                 int
	//JobResrcs JobResources
	JobState        int
	LastSchedEval   time.Time
	Licenses        string
	MailType        int
	MailUser        string
	MaxCpus         int
	MaxNodes        int
	McsLabel        string
	MemPerTres      string
	Name            string
	Network         string
	Nodes           string
	Nice            int
	NodeInx         int
	NTasksPerCore   int
	NTasksPerTres   int
	NTasksPerNode   int
	NTasksPerSocket int
	NTasksPerBoard  int
	NumCpus         int
	NumNodes        int
	NumTasks        int
	Partition       string
	PnMinMemory     int64
	PnMinCpus       int
	PnMinTmpDisk    int
	PowerFlags      int
	PreemptTime     time.Time
	PreemptableTime time.Time
	PreSusTime      time.Time
	Priority        int
	Profile         int
	Qos             string
	Reboot          int
	ReqNodes        string
	ReqNodeInx      int
	ReqSwitch       int
	Requeue         int
	ResizeTime      time.Time
	RestartCnt      int
	ResvName        string
	SchedNodes      string
	//SelectJobInfo
	Shared           int
	ShowFlags        int
	SiteFactor       int
	SocketsPerBoard  int
	SocketsPerNode   int
	StartTime        time.Time
	StartProtocolVer int
	StateDesc        string
	StateReason      int
	StdErr           string
	StdIn            string
	StdOut           string
	SubmitTime       time.Time
	SuspendTime      time.Time
	SystemComment    string
	TimeLimit        int
	TimeMin          int
	ThreadsPerCore   int
	Tresbind         string
	TresFreq         string
	TresPerJob       string
	TresPerNode      string
	TresPerSocket    string
	TresPerTask      string
	TresReqStr       string
	TresAllocStr     string
	UserId           int
	UserName         string
	Wait4Switch      int
	WcKey            string
	WorkDir          string
}
