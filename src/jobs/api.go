package jobs

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
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"unsafe"

	"github.com/gin-gonic/gin"
)

func InitializeJobsEndpoint(r *gin.Engine) {
	r.GET("/jobs", getJobsEndpoint)
	r.POST("/jobs/submit", submitJob)
}

func getJobsEndpoint(cnx *gin.Context) {
	jobs := queryAllJobs()
	cnx.JSON(200, jobs)
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
		req.ClusterUserId = 1
	}
	//TODO: Remove below when users are implemented
	req.ClusterUserId = 1

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

func createJobFromRequest(req JobDescriptorRequest, slurmJobId int) (job Job) {
	job.ClusterUserId = req.ClusterUserId
	job.JobId = slurmJobId
	job.Script = req.Script
	return job
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
