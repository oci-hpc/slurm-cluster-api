package jobs

import (
	"time"
)

func RunJobMonitor() {
	for {
		jobs := getJobInfos()
		for _, job := range jobs {
			upsertJobStatus(job)
		}
		time.Sleep(15 * time.Second)
	}
}
