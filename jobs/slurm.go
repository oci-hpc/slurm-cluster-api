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
*/
import "C"

import (
	"log"
	"time"
)

func getJobInfos() (jobs []JobInfo) {
	var slres *C.job_info_msg_t
	defer C.slurm_free_job_info_msg(slres)
	now := time.Now().Unix()
	ret := C.slurm_load_jobs(C.long(now-1000), &slres, 0)
	if ret == -1 {
		log.Printf("WARN: getJobs: error getting jobs from SLURM API")
	}
	if slres != nil {
		count := int(slres.record_count)
		jobs = convertJobInfoArray(slres.job_array, count)
	}
	return
}
