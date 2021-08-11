package jobs

import (
	"time"
	"unsafe"
)

type JobAllocationRequest struct {
	JobId      int
	SubmitTime time.Time
	OutputPath string
	JobName    string
	Account    string
	Script     string
}

type JobTemplateSubmission struct {
	Id                int
	JobId             int
	TemplateId        int
	TemplateKeyValues map[string]string
}

type SubmitResponse struct {
	JobId            int
	StepId           int
	ErrorCode        int
	JobSubmitUserMsg string
}

type Job struct {
	Id                    int
	JobId                 int
	ClusterUserId         int
	AccrueTime            time.Time
	EligibleTime          time.Time
	EndTime               time.Time
	PreemptTime           time.Time
	PreemptableTime       time.Time
	ResizeTime            time.Time
	StartTime             time.Time
	SuspendTime           time.Time
	WorkDir               string
	NTasksPerCore         int
	NTasksPerTres         int
	NTasksPerNode         int
	NTasksPerSocket       int
	NTasksPerBoard        int
	NumCpus               int
	NumNodes              int
	Script                string
	Command               string
	JobState              int
	JobStateReason        int
	JobStateDescription   string
	JobTemplateSubmission int
}

type JobDescriptorRequest struct {
	Username      string
	ClusterUserId int
	Account       string
	JobName       string
	Script        string
}

type TemplateJobSubmitRequest struct {
	TemplateId int
	KeyValues  map[string]string
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
