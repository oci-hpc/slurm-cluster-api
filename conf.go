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

	"github.com/gin-gonic/gin"
)

func InitializeConfEndpoint(r *gin.Engine) {
	r.GET("/conf", getConf)
}

func getConf(cnx *gin.Context) {
	var slres *C.slurm_conf_t
	//now := time.Now().Unix()
	ret := C.slurm_load_ctl_conf(0, &slres)

	if ret != 0 {
		cnx.JSON(500, gin.H{})
		C.slurm_free_ctl_conf(slres)
		return
	}

	res := convertConf(slres)

	C.slurm_free_ctl_conf(slres)

	if ret == 0 {
		cnx.JSON(200, res)
	} else {
		cnx.JSON(500, gin.H{})
	}
}

func convertConf(slres *C.slurm_conf_t) ConfOptions {
	var conf ConfOptions
	conf.LastUpdate = time.Unix(0, int64(slres.last_update))
	conf.AccountingStorageTres = C.GoString(slres.accounting_storage_tres)
	conf.AccountingStorageEnforce = int(slres.accounting_storage_enforce)
	conf.AccoutingStorageBackupHost = C.GoString(slres.accounting_storage_backup_host)
	conf.AccountingStorageExtHost = C.GoString(slres.accounting_storage_ext_host)
	conf.AccountingStorageHost = C.GoString(slres.accounting_storage_host)
	conf.AccountingStorageParams = C.GoString(slres.accounting_storage_params)
	conf.AccountingStoragePass = C.GoString(slres.accounting_storage_pass)
	conf.AccountingStoragePort = int(slres.accounting_storage_port)
	conf.AccountingStorageType = C.GoString(slres.accounting_storage_type)
	conf.AccountingStorageUser = C.GoString(slres.accounting_storage_user)
	conf.AcctGatherEnergyType = C.GoString(slres.acct_gather_energy_type)
	conf.AcctGatherProfileType = C.GoString(slres.acct_gather_profile_type)
	conf.AcctGatherInterconnectType = C.GoString(slres.acct_gather_interconnect_type)
	conf.AcctGatherFilesystemType = C.GoString(slres.acct_gather_filesystem_type)
	conf.AcctGatherNodeFreq = int(slres.acct_gather_node_freq)
	conf.AuthAltTypes = C.GoString(slres.authalttypes)
	conf.AuthInfo = C.GoString(slres.authinfo)
	conf.AuthHaltParams = C.GoString(slres.authalt_params)
	conf.BatchStartTimeout = int(slres.batch_start_timeout)
	conf.BBType = C.GoString(slres.bb_type)
	conf.BootTime = time.Unix(0, int64(slres.boot_time))
	conf.CliFilterPlugins = C.GoString(slres.cli_filter_plugins)
	conf.CoreSpecPlugin = C.GoString(slres.core_spec_plugin)
	conf.ClusterName = C.GoString(slres.cluster_name)
	conf.CommParams = C.GoString(slres.comm_params)
	conf.CompleteWait = int(slres.complete_wait)
	conf.ConfFlags = int(slres.conf_flags)
	conf.ControlAddr = C.GoString(*slres.control_addr)
	conf.ControlCnt = int(slres.control_cnt)
	conf.ControlMachine = C.GoString(*slres.control_machine)
	conf.CpuFreqDef = int(slres.cpu_freq_def)
	conf.CpuFreqGovs = int(slres.cpu_freq_govs)
	conf.CredType = C.GoString(slres.cred_type)
	conf.DebugFlags = int64(slres.debug_flags)
	conf.DefMemPerCpu = int64(slres.def_mem_per_cpu)
	conf.DependencyParams = C.GoString(slres.dependency_params)
	conf.EioTimeout = int(slres.eio_timeout)
	conf.EnforcePartLimits = int(slres.enforce_part_limits)
	conf.Epilog = C.GoString(slres.epilog)
	conf.EpilogMsgTime = int(slres.epilog_msg_time)
	conf.EpilogSlurmctld = C.GoString(slres.epilog_slurmctld)
	conf.ExtSensorsType = C.GoString(slres.ext_sensors_type)
	conf.ExtSensorsFreq = int(slres.ext_sensors_freq)
	conf.FedParams = C.GoString(slres.fed_params)
	conf.FirstJobId = int(slres.first_job_id)
	conf.FsDampeningFactor = int(slres.fs_dampening_factor)
	conf.GetEnvTimeout = int(slres.get_env_timeout)
	conf.GresPlugins = C.GoString(slres.gres_plugins)
	conf.GpuFreqDef = C.GoString(slres.gpu_freq_def)
	conf.HashVal = int(slres.hash_val)
	conf.HealthCheckInterval = int(slres.health_check_interval)
	conf.HealthCheckNodeState = int(slres.health_check_node_state)
	conf.HealthCheckProgram = C.GoString(slres.health_check_program)
	conf.InactiveLimit = int(slres.inactive_limit)
	conf.InteractiveStepOpts = C.GoString(slres.interactive_step_opts)
	conf.JobAcctGatherFreq = C.GoString(slres.job_acct_gather_freq)
	conf.JobAcctGatherType = C.GoString(slres.job_acct_gather_type)
	conf.JobAcctGatherParams = C.GoString(slres.job_acct_gather_params)
	conf.JobAcctOomKill = int(slres.job_acct_oom_kill)
	conf.JobCompHost = C.GoString(slres.job_comp_host)
	conf.JobCompLoc = C.GoString(slres.job_comp_loc)
	conf.JobCompParams = C.GoString(slres.job_comp_params)
	conf.JobCompPass = C.GoString(slres.job_comp_pass)
	conf.JobCompPort = int(slres.job_comp_port)
	conf.JobCompType = C.GoString(slres.job_comp_type)
	conf.JobCompUser = C.GoString(slres.job_comp_user)
	conf.JobContainerPlugin = C.GoString(slres.job_container_plugin)
	conf.JobCredentialPrivateKey = C.GoString(slres.job_credential_private_key)
	conf.JobCredentialPublicCertificate = C.GoString(slres.job_credential_public_certificate)
	conf.JobFileAppend = int(slres.job_file_append)
	conf.JobRequeue = int(slres.job_requeue)
	conf.JobSubmitPlugins = C.GoString(slres.job_submit_plugins)
	conf.KeepAliveTime = int(slres.keep_alive_time)
	conf.KillOnBadExit = int(slres.kill_on_bad_exit)
	conf.KillWait = int(slres.kill_wait)
	conf.LaunchParams = C.GoString(slres.launch_params)
	conf.LaunchType = C.GoString(slres.launch_type)
	conf.Licenses = C.GoString(slres.licenses)
	conf.LogFmt = int(slres.log_fmt)
	conf.MailDomain = C.GoString(slres.mail_domain)
	conf.MailProg = C.GoString(slres.mail_prog)
	conf.MaxArraySz = int(slres.max_array_sz)
	conf.MaxDbdMsgs = int(slres.max_dbd_msgs)
	conf.MaxJobCnt = int(slres.max_job_cnt)
	conf.MaxMemPerCpu = int(slres.max_mem_per_cpu)
	conf.MaxJobId = int(slres.max_job_id)
	conf.MaxStepCnt = int(slres.max_step_cnt)
	conf.McsPlugin = C.GoString(slres.mcs_plugin)
	conf.McsPluginParams = C.GoString(slres.mcs_plugin_params)
	conf.MinJobAge = int(slres.min_job_age)
	conf.MpiDefault = C.GoString(slres.mpi_default)
	conf.MpiParams = C.GoString(slres.mpi_params)
	conf.MsgTimeout = int(slres.msg_timeout)
	conf.NextJobId = int(slres.next_job_id)
	conf.NodeFeaturePlugins = C.GoString(slres.node_features_plugins)
	conf.NodePrefix = C.GoString(slres.node_prefix)
	conf.OverTimeLimit = int(slres.over_time_limit)
	conf.PluginDir = C.GoString(slres.plugindir)
	conf.PlugStack = C.GoString(slres.plugstack)
	conf.PowerParameters = C.GoString(slres.power_parameters)
	conf.PowerPlugin = C.GoString(slres.power_plugin)
	conf.PreemptExemptTime = int(slres.preempt_exempt_time)
	conf.PreemptMode = int(slres.preempt_mode)
	conf.PreemptType = C.GoString(slres.preempt_type)
	conf.PrepParams = C.GoString(slres.prep_params)
	conf.PrepPlugins = C.GoString(slres.prep_plugins)
	conf.PriorityDecayHl = int(slres.priority_decay_hl)
	conf.PriorityCalcPeriod = int(slres.priority_calc_period)
	conf.PriorityFavorSmall = int(slres.priority_favor_small)
	conf.PriorityFlags = int(slres.priority_flags)
	conf.PriorityResetPeriod = int(slres.priority_reset_period)
	conf.PriorityType = C.GoString(slres.priority_type)
	conf.PriorityWeightAge = int(slres.priority_weight_age)
	conf.PriorityWeightFs = int(slres.priority_weight_fs)
	conf.PriorityWeightJs = int(slres.priority_weight_js)
	conf.PriorityWeightQos = int(slres.priority_weight_qos)
	conf.PriorityWeightPart = int(slres.priority_weight_part)
	conf.PrivateData = int(slres.private_data)
	conf.Prolog = C.GoString(slres.prolog)
	conf.PrologEpilogTimeout = int(slres.prolog_epilog_timeout)
	conf.PrologSlurmctld = C.GoString(slres.prolog_slurmctld)
	conf.PropigatePrioProcess = int(slres.propagate_prio_process)
	conf.PrologFlags = int(slres.prolog_flags)
	conf.PropigateRLimits = C.GoString(slres.propagate_rlimits)
	conf.PropigateRLimitsExcept = C.GoString(slres.propagate_rlimits_except)
	conf.RebootProgram = C.GoString(slres.reboot_program)
	conf.ReconfigFlags = int(slres.reconfig_flags)
	conf.RequeueExit = C.GoString(slres.requeue_exit)
	conf.RequeueExitHold = C.GoString(slres.requeue_exit_hold)
	conf.ResumeFailProgram = C.GoString(slres.resume_fail_program)
	conf.ResumeProgram = C.GoString(slres.resume_program)
	conf.ResumeRate = int(slres.resume_rate)
	conf.ResumeTimeout = int(slres.resume_timeout)
	conf.ResvEpilog = C.GoString(slres.resv_epilog)
	conf.ResvOverRun = int(slres.resv_over_run)
	conf.ResvProlog = C.GoString(slres.resv_prolog)
	conf.Ret2Service = int(slres.ret2service)
	conf.RoutePlugin = C.GoString(slres.route_plugin)
	conf.SbCastParameters = C.GoString(slres.sbcast_parameters)
	conf.SchedLogLevel = int(slres.sched_log_level)
	conf.SchedLogfile = C.GoString(slres.sched_logfile)
	conf.SchedParams = C.GoString(slres.sched_params)
	conf.SchedTimeSlice = int(slres.sched_time_slice)
	conf.SchedType = C.GoString(slres.schedtype)
	conf.ScronParams = C.GoString(slres.scron_params)
	conf.SelectType = C.GoString(slres.select_type)
	conf.SelectTypeParam = int(slres.select_type_param)
	conf.SiteFactorPlugin = C.GoString(slres.site_factor_plugin)
	conf.SiteFactorParams = C.GoString(slres.site_factor_plugin)
	conf.SlurmConf = C.GoString(slres.slurm_conf)
	conf.SlurmUserId = int(slres.slurm_user_id)
	conf.SlurmUserName = C.GoString(slres.slurm_user_name)
	conf.SlurmdUserId = int(slres.slurmd_user_id)
	conf.SlurmdUserName = C.GoString(slres.slurmd_user_name)
	conf.SlurmctldAddr = C.GoString(slres.slurmctld_addr)
	conf.SlurmctldDebug = int(slres.slurmctld_debug)
	conf.SlurmctldLogfile = C.GoString(slres.slurmctld_logfile)
	conf.SlurmctldPidFile = C.GoString(slres.slurmctld_pidfile)
	conf.SlurmctldPlugstack = C.GoString(slres.slurmctld_plugstack)
	conf.SlurmctldPort = int(slres.slurmctld_port)
	conf.SlurmctldPortCount = int(slres.slurmctld_port_count)
	conf.SlurmctldPrimaryOffProg = C.GoString(slres.slurmctld_primary_off_prog)
	conf.SlurmctldPrimaryOnProg = C.GoString(slres.slurmctld_primary_on_prog)
	conf.SlurmctldSyslogDebug = int(slres.slurmctld_syslog_debug)
	conf.SlurmctldTimeout = int(slres.slurmctld_timeout)
	conf.SlurmctldParams = C.GoString(slres.slurmctld_params)
	conf.SlurmdDebug = int(slres.slurmd_debug)
	conf.SlurmdLogfile = C.GoString(slres.slurmd_logfile)
	conf.SlurmdParams = C.GoString(slres.slurmd_params)
	conf.SlurmdPidFile = C.GoString(slres.slurmd_pidfile)
	conf.SlurmdPort = int(slres.slurmd_port)
	conf.SlurmdSpoolDir = C.GoString(slres.slurmd_spooldir)
	conf.SlurmdSyslogDebug = int(slres.slurmd_syslog_debug)
	conf.SlurmdTimeout = int(slres.slurmd_timeout)
	conf.SrunEpilog = C.GoString(slres.srun_epilog)
	conf.SrunProlog = C.GoString(slres.srun_prolog)
	//conf.SrunPortRange = int(*slres.srun_port_range)
	conf.StateSaveLocation = C.GoString(slres.state_save_location)
	conf.SuspendExcNodes = C.GoString(slres.suspend_exc_nodes)
	conf.SuspendExcParts = C.GoString(slres.suspend_exc_parts)
	conf.SuspendProgram = C.GoString(slres.suspend_program)
	conf.SuspendRate = int(slres.suspend_rate)
	conf.SuspendTime = int(slres.suspend_time)
	conf.SwitchType = C.GoString(slres.switch_type)
	conf.TaskEpilog = C.GoString(slres.task_epilog)
	conf.TaskPlugin = C.GoString(slres.task_plugin)
	conf.TaskPluginParam = int(slres.task_plugin_param)
	conf.TaskProlog = C.GoString(slres.task_prolog)
	conf.TcpTimeout = int(slres.tcp_timeout)
	conf.TmpFs = C.GoString(slres.tmp_fs)
	conf.TopologyParam = C.GoString(slres.topology_param)
	conf.TopologyPlugin = C.GoString(slres.topology_plugin)
	conf.TreeWidth = int(slres.tree_width)
	conf.UnkillableProgram = C.GoString(slres.unkillable_program)
	conf.UnkillableTimeout = int(slres.unkillable_timeout)
	conf.Version = C.GoString(slres.version)
	conf.VsizeFactor = int(slres.vsize_factor)
	conf.WaitTime = int(slres.wait_time)
	conf.X11Params = C.GoString(slres.x11_params)
	return conf
}

type ConfOptions struct {
	AccountingStorageTres          string
	AccountingStorageEnforce       int
	AccoutingStorageBackupHost     string
	AccountingStorageHost          string
	AccountingStorageExtHost       string
	AccountingStorageParams        string
	AccountingStoragePass          string
	AccountingStoragePort          int
	AccountingStorageType          string
	AccountingStorageUser          string
	AcctGatherEnergyType           string
	AcctGatherProfileType          string
	AcctGatherInterconnectType     string
	AcctGatherFilesystemType       string
	AcctGatherNodeFreq             int
	AuthAltTypes                   string
	AuthInfo                       string
	AuthHaltParams                 string
	BatchStartTimeout              int
	BBType                         string
	BootTime                       time.Time
	CliFilterPlugins               string
	CoreSpecPlugin                 string
	ClusterName                    string
	CommParams                     string
	CompleteWait                   int
	ConfFlags                      int
	ControlAddr                    string
	ControlCnt                     int
	ControlMachine                 string
	CpuFreqDef                     int
	CpuFreqGovs                    int
	CredType                       string
	DebugFlags                     int64
	DefMemPerCpu                   int64
	DependencyParams               string
	EioTimeout                     int
	EnforcePartLimits              int
	Epilog                         string
	EpilogMsgTime                  int
	EpilogSlurmctld                string
	ExtSensorsType                 string
	ExtSensorsFreq                 int
	FedParams                      string
	FirstJobId                     int
	FsDampeningFactor              int
	GetEnvTimeout                  int
	GresPlugins                    string
	GroupTime                      int
	GroupForce                     int
	GpuFreqDef                     string
	HashVal                        int
	HealthCheckInterval            int
	HealthCheckNodeState           int
	HealthCheckProgram             string
	InactiveLimit                  int
	InteractiveStepOpts            string
	JobAcctGatherFreq              string
	JobAcctGatherType              string
	JobAcctGatherParams            string
	JobAcctOomKill                 int
	JobCompHost                    string
	JobCompLoc                     string
	JobCompParams                  string
	JobCompPass                    string
	JobCompPort                    int
	JobCompType                    string
	JobCompUser                    string
	JobContainerPlugin             string
	JobCredentialPrivateKey        string
	JobCredentialPublicCertificate string
	//JobDefaultsList []string
	JobFileAppend           int
	JobRequeue              int
	JobSubmitPlugins        string
	KeepAliveTime           int
	KillOnBadExit           int
	KillWait                int
	LastUpdate              time.Time
	LaunchParams            string
	LaunchType              string
	Licenses                string
	LogFmt                  int
	MailDomain              string
	MailProg                string
	MaxArraySz              int
	MaxDbdMsgs              int
	MaxJobCnt               int
	MaxJobId                int
	MaxMemPerCpu            int
	MaxStepCnt              int
	MaxTasksPerNode         int
	McsPlugin               string
	McsPluginParams         string
	MinJobAge               int
	MpiDefault              string
	MpiParams               string
	MsgTimeout              int
	NextJobId               int
	NodeFeaturePlugins      string
	NodePrefix              string
	OverTimeLimit           int
	PluginDir               string
	PlugStack               string
	PowerParameters         string
	PowerPlugin             string
	PreemptExemptTime       int
	PreemptMode             int
	PreemptType             string
	PrepParams              string
	PrepPlugins             string
	PriorityDecayHl         int
	PriorityCalcPeriod      int
	PriorityFavorSmall      int
	PriorityFlags           int
	PriorityMaxAge          int
	PriorityParams          string
	PriorityResetPeriod     int
	PriorityType            string
	PriorityWeightAge       int
	PriorityWeightAssoc     int
	PriorityWeightFs        int
	PriorityWeightJs        int
	PriorityWeightPart      int
	PriorityWeightQos       int
	PrirotiyWeightTres      string
	PrivateData             int
	ProTrackType            string
	Prolog                  string
	PrologEpilogTimeout     int
	PrologSlurmctld         string
	PropigatePrioProcess    int
	PrologFlags             int
	PropigateRLimits        string
	PropigateRLimitsExcept  string
	RebootProgram           string
	ReconfigFlags           int
	RequeueExit             string
	RequeueExitHold         string
	ResumeFailProgram       string
	ResumeProgram           string
	ResumeRate              int
	ResumeTimeout           int
	ResvEpilog              string
	ResvOverRun             int
	ResvProlog              string
	Ret2Service             int
	RoutePlugin             string
	SbCastParameters        string
	SchedLogfile            string
	SchedLogLevel           int
	SchedParams             string
	SchedTimeSlice          int
	SchedType               string
	ScronParams             string
	SelectType              string
	SelectTypeParam         int
	SiteFactorPlugin        string
	SiteFactorParams        string
	SlurmConf               string
	SlurmUserId             int
	SlurmUserName           string
	SlurmdUserId            int
	SlurmdUserName          string
	SlurmctldAddr           string
	SlurmctldDebug          int
	SlurmctldLogfile        string
	SlurmctldPidFile        string
	SlurmctldPlugstack      string
	SlurmctldPort           int
	SlurmctldPortCount      int
	SlurmctldPrimaryOffProg string
	SlurmctldPrimaryOnProg  string
	SlurmctldSyslogDebug    int
	SlurmctldTimeout        int
	SlurmctldParams         string
	SlurmdDebug             int
	SlurmdLogfile           string
	SlurmdParams            string
	SlurmdPidFile           string
	SlurmdPort              int
	SlurmdSpoolDir          string
	SlurmdSyslogDebug       int
	SlurmdTimeout           int
	SrunEpilog              string
	SrunPortRange           int
	SrunProlog              string
	StateSaveLocation       string
	SuspendExcNodes         string
	SuspendExcParts         string
	SuspendProgram          string
	SuspendRate             int
	SuspendTime             int
	SuspendTimeout          int
	SwitchType              string
	TaskEpilog              string
	TaskPlugin              string
	TaskPluginParam         int
	TaskProlog              string
	TcpTimeout              int
	TmpFs                   string
	TopologyParam           string
	TopologyPlugin          string
	TreeWidth               int
	UnkillableProgram       string
	UnkillableTimeout       int
	Version                 string
	VsizeFactor             int
	WaitTime                int
	X11Params               string
}
