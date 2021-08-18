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
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	templateRepo "github.com/oci-hpc/slurm-cluster-api/src/template"
)

func InitializeJobsEndpoint(r *gin.Engine) {
	r.GET("/jobs", getJobsEndpoint)
	r.GET("/jobs/template", getJobsTemplateSubmit)
	r.POST("/jobs/submit", submitJob)
	r.POST("/jobs/submit/template", templateJobSubmit)
}

func getJobsEndpoint(cnx *gin.Context) {
	jobs := queryAllJobs()
	cnx.JSON(200, jobs)
}

func getJobsTemplateSubmit(cnx *gin.Context) {
	query := cnx.Request.URL.Query()
	if val, ok := query["id"]; ok || len(val) > 0 {
		intVal, err := strconv.Atoi(val[0])
		if err != nil {
			cnx.JSON(400, "Invalid value for id")
		}
		jobs := queryJobTemplateSubmission(intVal)
		cnx.JSON(200, jobs)
		return
	}
	cnx.JSON(400, "id field is required")
}

func templateJobSubmit(cnx *gin.Context) {
	jsonData, err := ioutil.ReadAll(cnx.Request.Body)

	if err != nil {
		cnx.JSON(400, err.Error())
		return
	}

	var req TemplateJobSubmitRequest
	json.Unmarshal([]byte(jsonData), &req)

	selectedTemplate := templateRepo.QueryTemplateById(req.TemplateId)
	combined := selectedTemplate.Header + "\n" + selectedTemplate.Body
	t := template.Must(template.New("t2").Parse(combined))
	var result string
	buf := bytes.NewBufferString(result)
	templateKeyValues := map[string]string{}
	for key, value := range req.KeyValues {
		templateKeyValues[strings.ToLower(key)] = value
	}
	t.Execute(buf, templateKeyValues)
	var jobReq JobDescriptorRequest
	jobReq.Script = buf.String()
	print(jobReq.Script)
	jobReq.Account = "default"
	jobReq.JobName = "test-job"
	jobId, err := slurmSubmitJob(jobReq)
	if err != nil {
		cnx.JSON(500, "Internal error submitting job")
	} else {
		var templateJobSubmission JobTemplateSubmission
		templateJobSubmission.TemplateId = selectedTemplate.Id
		templateJobSubmission.JobId = jobId
		kv, _ := json.Marshal(req.KeyValues)
		templateJobSubmission.TemplateKeyValues = string(kv)
		insertJobTemplateSubmission(templateJobSubmission)
		cnx.JSON(200, "Job successfully submitted")
	}
	//Pass in a map[string]string with keys in lowercase
	//t.Execute(os.Stdout, passIn)
}

func submitJob(cnx *gin.Context) {
	jsonData, err := ioutil.ReadAll(cnx.Request.Body)

	if err != nil {
		cnx.JSON(400, err.Error())
		return
	}

	var req JobDescriptorRequest
	json.Unmarshal([]byte(jsonData), &req)

	_, err = slurmSubmitJob(req)
	if err != nil {
		cnx.JSON(500, "Internal error submitting job")
	} else {
		cnx.JSON(200, "Job successfully submitted")
	}
}

func slurmSubmitJob(req JobDescriptorRequest) (jobId int, err error) {
	if req.Username == "" {
		req.Username = "DefaultUser"
		req.ClusterUserId = 1
	}
	//TODO: Remove below when users are implemented
	req.ClusterUserId = 1

	pathEnv := os.Getenv("PATH")
	pathSetString := "PATH=/usr:" + pathEnv

	outputDirPath := os.Getenv("OUTPUT_DIR")
	outputUserPath := path.Join(outputDirPath, req.Username)

	if _, err = os.Stat(outputUserPath); os.IsNotExist(err) {
		err = os.Mkdir(outputUserPath, 0755)
		if err != nil {
			log.Println("WARN: slurmSubmitJob: " + err.Error())
			return
		}
	}
	outputWorkDir, err := createOutputDirectory(req.Username)
	if err != nil {
		log.Println("WARN: slurmSubmitJob: " + err.Error())
		return
	}

	//job_desc_msg := convertJobDescriptor(jobDescriptor)
	var job_desc_msg C.job_desc_msg_t
	C.slurm_init_job_desc_msg(&job_desc_msg)
	source := []string{pathSetString}
	print(pathSetString)
	job_desc_msg.env_size = C.uint(len(source))
	job_desc_msg.environment = getCStringArray(source)
	job_desc_msg.work_dir = C.CString(outputWorkDir)
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
		jobId = res.JobId
		return
	} else {
		errno := C.slurm_get_errno()
		errno_str := "SLURM-" + strconv.Itoa(int(errno)) + " " + C.GoString(C.slurm_strerror(errno))
		log.Println("WARN: slurmSubmitJob: " + errno_str)
		err = errors.New(errno_str)
		return
	}
}

func createOutputDirectory(username string) (ouputDirectoryPath string, err error) {
	dateString := time.Now().Format("2006-01-02")

	outputDirPath := os.Getenv("OUTPUT_DIR")
	outputUserDatePath := path.Join(outputDirPath, username, dateString)

	if _, err = os.Stat(outputUserDatePath); os.IsNotExist(err) {
		err = os.Mkdir(outputUserDatePath, 0755)
		if err != nil {
			print("WARN: createOutputDirectory: " + err.Error())
			return
		}
	}

	fileCount, err := countDirectoryFiles(outputUserDatePath)
	if err != nil {
		print("WARN: createOutputDirectory: " + err.Error())
	}
	fileCount = fileCount + 1

	ouputDirectoryPath = path.Join(outputUserDatePath, strconv.Itoa(fileCount))
	err = os.Mkdir(ouputDirectoryPath, 0755)
	if err != nil {
		print("WARN: createOutputDirectory: " + err.Error())
		return
	}

	return
}

func countDirectoryFiles(path string) (files int, err error) {
	d, err := os.ReadDir(path)
	if err != nil {
		return
	}
	files = len(d)
	return
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
