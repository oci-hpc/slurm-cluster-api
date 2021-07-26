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
static char**makeCharArray(int size) {
  return calloc(sizeof(char*), size);
}
static void setArrayString(char **a, char *s, int n) {
	a[n] = s;
}
*/
import "C"

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func InitializeJobsEndpoint(r *gin.Engine) {
	r.GET("/jobs", getJobsEndpoint)
	r.POST("/jobs/submit", submitJob)
}

func getJobsEndpoint(cnx *gin.Context) {
	var slres *C.job_info_msg_t
	now := time.Now().Unix()
	ret := C.slurm_load_jobs(C.long(now-1000), &slres, 0)
	count := int(slres.record_count)
	jobs := convertJobInfoArray(slres.job_array, count)

	if ret == 0 {
		cnx.JSON(200, jobs)
	} else {
		cnx.JSON(500, gin.H{})
	}

	C.slurm_free_job_info_msg(slres)
}

func getJobInfos() (jobs []JobInfo) {
	var slres *C.job_info_msg_t
	defer C.slurm_free_job_info_msg(slres)
	now := time.Now().Unix()
	ret := C.slurm_load_jobs(C.long(now-1000), &slres, 0)
	if ret == -1 {
		log.Printf("WARN: getJobs: error getting jobs from SLURM API")
	}
	count := int(slres.record_count)
	jobs = convertJobInfoArray(slres.job_array, count)
	return
}

func RunJobMonitor() {
	for {
		jobs := getJobInfos()
		for _, job := range jobs {
			upsertJobStatus(job)
		}
		time.Sleep(15 * time.Second)
	}
}

func upsertJobStatus(jobInfo JobInfo) {
	job := convertJobInfoToJob(jobInfo)
	res := queryJobsBySlurmJobId(jobInfo.JobId)
	if res.Id == 0 {
		insertJob(job)
	} else {
		job.Id = res.Id
		updateJob(job)
	}
}

func convertJobInfoToJob(jobInfo JobInfo) (job Job) {
	job.AccrueTime = jobInfo.AccrueTime
	job.Command = jobInfo.Command
	job.EligibleTime = jobInfo.EligibleTime
	job.EndTime = jobInfo.EndTime
	job.JobId = jobInfo.JobId
	job.JobState = jobInfo.JobState
	job.JobStateDescription = jobInfo.StateDesc
	job.JobStateReason = jobInfo.StateReason
	job.NTasksPerBoard = jobInfo.NTasksPerBoard
	job.NTasksPerCore = jobInfo.NTasksPerCore
	job.NTasksPerNode = jobInfo.NTasksPerNode
	job.NTasksPerSocket = jobInfo.NTasksPerSocket
	job.NTasksPerTres = jobInfo.NTasksPerTres
	job.NumCpus = jobInfo.NumCpus
	job.NumNodes = jobInfo.NumNodes
	job.PreemptTime = jobInfo.PreemptTime
	job.PreemptableTime = jobInfo.PreemptableTime
	job.ResizeTime = jobInfo.ResizeTime
	job.StartTime = jobInfo.StartTime
	job.SuspendTime = jobInfo.SuspendTime
	job.WorkDir = jobInfo.WorkDir
	return job
}

func createJobFromRequest(req JobDescriptorRequest, slurmJobId int) (job Job) {
	job.JobId = slurmJobId
	job.Script = req.Script
	return job
}

func convertRowsToJobs(rows *sql.Rows, jobs *[]Job) {
	defer rows.Close()
	if rows == nil {
		return
	}
	for rows.Next() {
		var job Job
		err := rows.Scan(
			&job.Id,
			&job.JobId,
			&job.ClusterUserId,
			&job.AccrueTime,
			&job.EligibleTime,
			&job.EndTime,
			&job.PreemptTime,
			&job.PreemptableTime,
			&job.ResizeTime,
			&job.StartTime,
			&job.SuspendTime,
			&job.WorkDir,
			&job.NTasksPerCore,
			&job.NTasksPerTres,
			&job.NTasksPerNode,
			&job.NTasksPerSocket,
			&job.NTasksPerBoard,
			&job.NumCpus,
			&job.NumNodes,
			&job.Script,
			&job.Command,
			&job.JobState,
			&job.JobStateReason,
			&job.JobStateDescription,
		)
		if err != nil {
			log.Printf("WARN: convertRowsToJobs: " + err.Error())
		}
		*jobs = append(*jobs, job)
	}
}

func insertJob(job Job) {
	sqlString := `
		INSERT INTO t_job (
			m_job_id,
			m_user_id,
			m_accrue_time,
			m_eligible_time,
			m_end_time,
			m_preempt_time,
			m_preemptable_time,
			m_resize_time,
			m_start_time,
			m_suspend_time,
			m_work_dir,
			m_n_tasks_per_core,
			m_n_tasks_per_tres,
			m_n_tasks_per_node,
			m_n_tasks_per_socket,
			m_n_tasks_per_board,
			m_num_cpus,
			m_num_nodes,
			m_script,
			m_command,
			m_job_state,
			m_job_state_reason,
			m_job_state_description
		) values (
			:m_job_id,
			:m_user_id,
			:m_accrue_time,
			:m_eligible_time,
			:m_end_time,
			:m_preempt_time,
			:m_preemptable_time,
			:m_resize_time,
			:m_start_time,
			:m_suspend_time,
			:m_work_dir,
			:m_n_tasks_per_core,
			:m_n_tasks_per_tres,
			:m_n_tasks_per_node,
			:m_n_tasks_per_socket,
			:m_n_tasks_per_board,
			:m_num_cpus,
			:m_num_nodes,
			:m_script,
			:m_command,
			:m_job_state,
			:m_job_state_reason,
			:m_job_state_description
		)
	`
	db := GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("m_job_id", job.JobId),
		sql.Named("m_user_id", job.ClusterUserId),
		sql.Named("m_accrue_time", job.AccrueTime),
		sql.Named("m_eligible_time", job.EligibleTime),
		sql.Named("m_end_time", job.EndTime),
		sql.Named("m_preempt_time", job.PreemptTime),
		sql.Named("m_preemptable_time", job.PreemptableTime),
		sql.Named("m_resize_time", job.ResizeTime),
		sql.Named("m_start_time", job.StartTime),
		sql.Named("m_suspend_time", job.SuspendTime),
		sql.Named("m_work_dir", job.WorkDir),
		sql.Named("m_n_tasks_per_core", job.NTasksPerCore),
		sql.Named("m_n_tasks_per_tres", job.NTasksPerTres),
		sql.Named("m_n_tasks_per_node", job.NTasksPerNode),
		sql.Named("m_n_tasks_per_socket", job.NTasksPerSocket),
		sql.Named("m_n_tasks_per_board", job.NTasksPerBoard),
		sql.Named("m_num_cpus", job.NumCpus),
		sql.Named("m_num_nodes", job.NumNodes),
		sql.Named("m_script", job.Script),
		sql.Named("m_command", job.Command),
		sql.Named("m_job_state", job.JobState),
		sql.Named("m_job_state_reason", job.JobStateReason),
		sql.Named("m_job_state_description", job.JobStateDescription),
	)
	if err != nil {
		log.Printf("WARN: insertJob: " + err.Error())
	}
}

func updateJob(job Job) {
	sqlString := `
		UPDATE t_job
		SET
			m_accrue_time = :m_accrue_time,
			m_eligible_time = :m_eligible_time,
			m_end_time = :m_end_time,
			m_preempt_time = :m_preempt_time,
			m_preemptable_time = :m_preemptable_time,
			m_resize_time = :m_resize_time,
			m_start_time = :m_start_time,
			m_suspend_time = :m_suspend_time,
			m_work_dir = :m_work_dir,
			m_n_tasks_per_core = :m_n_tasks_per_core,
			m_n_tasks_per_tres = :m_n_tasks_per_core,
			m_n_tasks_per_node = :m_n_tasks_per_node,
			m_n_tasks_per_socket = :m_n_tasks_per_socket,
			m_n_tasks_per_board = :m_n_tasks_per_board,
			m_num_cpus = :m_num_cpus,
			m_num_nodes = :m_num_nodes,
			m_command = :m_command,
			m_job_state = :m_job_state,
			m_job_state_reason = :m_job_state_reason,
			m_job_state_description = :m_job_state_description
		WHERE id = :id
	`
	db := GetDbConnection()
	defer db.Close()
	_, err := db.Exec(
		sqlString,
		sql.Named("m_accrue_time", job.AccrueTime),
		sql.Named("m_eligible_time", job.EligibleTime),
		sql.Named("m_end_time", job.EndTime),
		sql.Named("m_preempt_time", job.PreemptTime),
		sql.Named("m_preemptable_time", job.PreemptableTime),
		sql.Named("m_resize_time", job.ResizeTime),
		sql.Named("m_start_time", job.StartTime),
		sql.Named("m_suspend_time", job.SuspendTime),
		sql.Named("m_work_dir", job.WorkDir),
		sql.Named("m_n_tasks_per_core", job.NTasksPerCore),
		sql.Named("m_n_tasks_per_tres", job.NTasksPerTres),
		sql.Named("m_n_tasks_per_node", job.NTasksPerNode),
		sql.Named("m_n_tasks_per_socket", job.NTasksPerSocket),
		sql.Named("m_n_tasks_per_board", job.NTasksPerBoard),
		sql.Named("m_num_cpus", job.NumCpus),
		sql.Named("m_num_nodes", job.NumNodes),
		sql.Named("m_command", job.Command),
		sql.Named("m_job_state", job.JobState),
		sql.Named("m_job_state_reason", job.JobStateReason),
		sql.Named("m_job_state_description", job.JobStateDescription),
		sql.Named("id", job.Id),
	)
	if err != nil {
		log.Printf("WARN: updateJob: " + err.Error())
	}
}

func queryAllJobs() (jobs []Job) {
	sqlString := `
		SELECT 
			m_job_id,
			m_user_id,
			m_accrue_time,
			m_eligible_time,
			m_end_time,
			m_preempt_time,
			m_preemptable_time,
			m_resize_time,
			m_start_time,
			m_suspend_time,
			m_work_dir,
			m_n_tasks_per_core,
			m_n_tasks_per_tres,
			m_n_tasks_per_node,
			m_n_tasks_per_socket,
			m_n_tasks_per_board,
			m_num_cpus,
			m_num_nodes,
			m_script,
			m_command,
			m_job_state,
			m_job_state_reason,
			m_job_state_description
		FROM t_job;
	`
	db := GetDbConnection()
	defer db.Close()
	rows, err := db.Query(sqlString)
	if err != nil {
		log.Printf("WARN: queryAllJobs: " + err.Error())
	}
	convertRowsToJobs(rows, &jobs)
	err = rows.Err()
	if err != nil {
		log.Printf("WARN: queryAllJobs: " + err.Error())
	}
	return jobs
}

func queryJobsByUser(clusterUserId int) (jobs []Job) {
	sqlString := `
		SELECT 
		  id
			m_job_id,
			m_user_id,
			m_accrue_time,
			m_eligible_time,
			m_end_time,
			m_preempt_time,
			m_preemptable_time,
			m_resize_time,
			m_start_time,
			m_suspend_time,
			m_work_dir,
			m_n_tasks_per_core,
			m_n_tasks_per_tres,
			m_n_tasks_per_node,
			m_n_tasks_per_socket,
			m_n_tasks_per_board,
			m_num_cpus,
			m_num_nodes,
			m_script,
			m_command,
			m_job_state,
			m_job_state_reason,
			m_job_state_description
		FROM t_job
		WHERE m_user_id = :m_user_id;
	`
	db := GetDbConnection()
	defer db.Close()
	rows, err := db.Query(sqlString, sql.Named("m_user_id", clusterUserId))
	if err != nil {
		log.Printf("WARN: queryJobsByUser: " + err.Error())
	}
	convertRowsToJobs(rows, &jobs)
	err = rows.Err()
	if err != nil {
		log.Printf("WARN: queryJobsByUser: " + err.Error())
	}
	return jobs
}

func queryJobsBySlurmJobId(slurmJobId int) (job Job) {
	sqlString := `
		SELECT 
		  id,
			m_job_id,
			m_user_id,
			m_accrue_time,
			m_eligible_time,
			m_end_time,
			m_preempt_time,
			m_preemptable_time,
			m_resize_time,
			m_start_time,
			m_suspend_time,
			m_work_dir,
			m_n_tasks_per_core,
			m_n_tasks_per_tres,
			m_n_tasks_per_node,
			m_n_tasks_per_socket,
			m_n_tasks_per_board,
			m_num_cpus,
			m_num_nodes,
			m_script,
			m_command,
			m_job_state,
			m_job_state_reason,
			m_job_state_description
		FROM t_job
		WHERE m_job_id = :m_job_id;
	`
	db := GetDbConnection()
	defer db.Close()
	rows, err := db.Query(sqlString, sql.Named("m_job_id", slurmJobId))
	if err != nil {
		log.Printf("WARN: queryJobsBySlurmJobId: " + err.Error())
	}
	var jobs []Job
	convertRowsToJobs(rows, &jobs)
	err = rows.Err()
	if err != nil {
		log.Printf("WARN: queryJobsBySlurmJobId: " + err.Error())
	}
	if len(jobs) == 0 {
		return
	}
	return jobs[0]
}

func submitJob(cnx *gin.Context) {
	jsonData, err := ioutil.ReadAll(cnx.Request.Body)

	if err != nil {
		cnx.JSON(400, err.Error())
		return
	}

	var req JobDescriptorRequest
	json.Unmarshal([]byte(jsonData), &req)

	if req.Username == "" {
		req.Username = "DefaultUser"
	}

	pathEnv := os.Getenv("PATH")
	pathSetString := "PATH=" + pathEnv

	outputDirPath := os.Getenv("OUTPUT_DIR")
	outputUserPath := path.Join(outputDirPath, req.Username)

	if _, err := os.Stat(outputUserPath); os.IsNotExist(err) {
		err = os.Mkdir(outputUserPath, 0755)
		if err != nil {
			print("WARN: submitJob: " + err.Error())
			cnx.JSON(500, "Error creating path for user")
			return
		}
	}

	//job_desc_msg := convertJobDescriptor(jobDescriptor)
	var job_desc_msg C.job_desc_msg_t
	C.slurm_init_job_desc_msg(&job_desc_msg)
	source := []string{pathSetString}
	job_desc_msg.env_size = C.uint(len(source))
	job_desc_msg.environment = getCStringArray(source)
	job_desc_msg.work_dir = C.CString(outputUserPath)
	job_desc_msg.name = C.CString(req.JobName)
	job_desc_msg.account = C.CString(req.Account)
	job_desc_msg.script = C.CString(req.Script)
	job_desc_msg.user_id = 1000
	job_desc_msg.group_id = 1000

	var slres *C.submit_response_msg_t
	defer C.slurm_free_submit_response_response_msg(slres)

	ret := C.slurm_submit_batch_job(&job_desc_msg, &slres)
	if ret == 0 {
		res := convertSubmitResponse(slres)
		job := createJobFromRequest(req, res.JobId)
		insertJob(job)
		cnx.JSON(200, res)
	} else {
		errno := C.slurm_get_errno()
		errno_str := "SLURM-" + strconv.Itoa(int(errno)) + " " + C.GoString(C.slurm_strerror(errno))
		cnx.JSON(500, errno_str)
	}

}

func getCStringArray(source []string) **C.char {
	cArray := C.malloc(C.size_t(len(source)) * C.size_t(unsafe.Sizeof(uintptr(0))))

	// convert the C array to a Go Array so we can index it
	a := (*[1<<30 - 1]*C.char)(cArray)

	for idx, val := range source {
		a[idx] = C.CString(val)
	}

	return (**C.char)(cArray)
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

func convertJobDescriptor(jobDescriptor JobDescriptor) C.job_desc_msg_t {
	var job_desc_msg C.job_desc_msg_t
	C.slurm_init_job_desc_msg(&job_desc_msg)

	job_desc_msg.name = C.CString(jobDescriptor.Name)
	job_desc_msg.script = C.CString(jobDescriptor.Script)
	return job_desc_msg
}

func convertSubmitResponse(slres *C.submit_response_msg_t) SubmitResponse {
	var submitResponse SubmitResponse
	submitResponse.JobId = int(slres.job_id)
	submitResponse.ErrorCode = int(slres.error_code)
	submitResponse.StepId = int(slres.step_id)
	submitResponse.JobSubmitUserMsg = C.GoString(slres.job_submit_user_msg)
	return submitResponse
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

type JobAllocationRequest struct {
	JobId      int
	SubmitTime time.Time
	OutputPath string
	JobName    string
	Account    string
	Script     string
}

type SubmitResponse struct {
	JobId            int
	StepId           int
	ErrorCode        int
	JobSubmitUserMsg string
}

type Job struct {
	Id                  int
	JobId               int
	ClusterUserId       int
	AccrueTime          time.Time
	EligibleTime        time.Time
	EndTime             time.Time
	PreemptTime         time.Time
	PreemptableTime     time.Time
	ResizeTime          time.Time
	StartTime           time.Time
	SuspendTime         time.Time
	WorkDir             string
	NTasksPerCore       int
	NTasksPerTres       int
	NTasksPerNode       int
	NTasksPerSocket     int
	NTasksPerBoard      int
	NumCpus             int
	NumNodes            int
	Script              string
	Command             string
	JobState            int
	JobStateReason      int
	JobStateDescription string
}

type JobDescriptorRequest struct {
	Username string
	Account  string
	JobName  string
	Script   string
}

type JobDescriptor struct {
	Account       string
	AcctgFreq     string
	AdminComment  string
	AllocNode     string
	AllocRespPort int
	AllocSid      int
	Argc          int
	Argv          string
	//ArrayBitmap unsafe.Pointer
	BatchFeatures     string
	BitFlags          int
	BeginTime         time.Time
	BurstBuffer       string
	Clusters          string
	ClusterFeatures   string
	Comment           string
	Contiguous        int
	CoreSpec          int
	CpuBind           string
	CpuBindType       int
	CpuFreqMin        int
	CpuFreqMax        int
	CpuFreqGov        int
	CpusPerTres       string
	CrontabEntry      unsafe.Pointer
	Deadline          time.Time
	DelayBoot         int
	Dependency        string
	EndTime           time.Time
	Environment       string
	EnvSize           int
	Extra             string
	ExcNodes          string
	Features          string
	FedSiblingsActive int64
	FedSiblingsViable int64
	GroupId           int
	HetJobOffset      int
	Immediate         int
	JobId             int
	JobIdStr          string
	KillNodeOnFail    int
	Licenses          string
	MailType          int
	MailUser          string
	McsLabel          string
	MemBind           string
	MemBindType       int
	MemPerTres        string
	Name              string
	Network           string
	Nodes             string
	Nice              int
	NumTasks          int
	OpenMode          int
	OriginCluster     string
	OtherPort         int
	Overcommit        int
	Partition         string
	PlaneSize         int
	PowerFlags        int
	Priority          int
	Profile           int
	Qos               string
	Reboot            int
	RespHost          string
	ReqNodes          string
	RestartCnt        int
	Requeue           int
	Reservation       string
	//Script is the execution, full batch file written to a string
	Script    string
	ScriptBuf unsafe.Pointer
	//SelectJobInfo
	Shared          int
	SiteFactor      int
	SpankJobEnv     string
	SpankJobEnvSize int
	TaskDist        int
	TimeLimit       int
	TimeMin         int
	Tresbind        string
	TresFreq        string
	TresPerJob      string
	TresPerNode     string
	TresPerSocket   string
	TresPerTask     string
	ClusterUserId   int
	SlurmUserId     int
	WaitAllNodes    int
	WarnFlags       int
	WarnSignal      int
	WarnTime        int
	CpusPerTask     int
	MinCpus         int
	MaxCpus         int
	MinNodes        int
	MaxNodes        int
	BoardsPerNode   int
	SocketsPerBoard int
	SocketsPerNode  int
	CoresPerSocket  int
	ThreadsPerCore  int
	NTasksPerCore   int
	NTasksPerTres   int
	NTasksPerNode   int
	NTasksPerSocket int
	NTasksPerBoard  int
	PnMinMemory     int64
	PnMinCpus       int
	PnMinTmpDisk    int
	ReqSwitch       int
	StdErr          string
	StdIn           string
	StdOut          string
	TresReqCnt      int64
	Wait4Switch     int
	WcKey           string
	WorkDir         string
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
