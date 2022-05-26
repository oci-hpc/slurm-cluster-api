package jobs

import (
	"database/sql"
	"log"

	db "github.com/oci-hpc/database"
)

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

func convertRowsToJobs(rows *sql.Rows, jobs *[]Job) {
	if rows == nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var job Job
		var nullTemplateSubmissionId sql.NullInt32
		var nullTemplateId sql.NullInt32
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
			&nullTemplateSubmissionId,
			&nullTemplateId,
		)
		if err != nil {
			log.Printf("WARN: convertRowsToJobs: " + err.Error())
		}
		if nullTemplateSubmissionId.Valid {
			job.JobTemplateSubmission = int(nullTemplateSubmissionId.Int32)
			job.JobTemplateId = int(nullTemplateId.Int32)
		}
		*jobs = append(*jobs, job)
	}
}

func insertJobTemplateSubmission(submission JobTemplateSubmission) (id int64) {
	sqlString := `
		INSERT INTO t_job_template_submission (
			m_job_id,
			m_template_id,
			m_template_key_values
		) VALUES (
			:m_job_id,
			:m_template_id,
			:m_template_key_values
		)
	`
	db := db.GetDbConnection()
	defer db.Close()
	res, err := db.Exec(
		sqlString,
		sql.Named("m_job_id", submission.JobId),
		sql.Named("m_template_id", submission.TemplateId),
		sql.Named("m_template_key_values", submission.TemplateKeyValues),
	)
	if err != nil {
		log.Printf("WARN: insertJobTemplateSubmission: " + err.Error())
	}
	id, _ = res.LastInsertId()
	return
}

func queryJobTemplateSubmission(id int) (jobTemplateSubmission JobTemplateSubmission) {
	sqlString := `
		SELECT 
		  id,
			m_job_id,
			m_template_id,
			m_template_key_values
		FROM t_job_template_submission
		WHERE id = :id
	`
	db := db.GetDbConnection()
	defer db.Close()
	rows, err := db.Query(
		sqlString,
		sql.Named("id", id),
	)
	if err != nil {
		log.Printf("WARN: queryJobTemplateSubmission: " + err.Error())
	}
	defer rows.Close()
	if rows == nil {
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&jobTemplateSubmission.Id,
			&jobTemplateSubmission.JobId,
			&jobTemplateSubmission.TemplateId,
			&jobTemplateSubmission.TemplateKeyValues,
		)
		if err != nil {
			log.Printf("WARN: queryJobTemplateSubmission: " + err.Error())
		}
	}
	return jobTemplateSubmission
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
	db := db.GetDbConnection()
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
	db := db.GetDbConnection()
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
			t_job.id,
			t_job.m_job_id,
			t_job.m_user_id,
			t_job.m_accrue_time,
			t_job.m_eligible_time,
			t_job.m_end_time,
			t_job.m_preempt_time,
			t_job.m_preemptable_time,
			t_job.m_resize_time,
			t_job.m_start_time,
			t_job.m_suspend_time,
			t_job.m_work_dir,
			t_job.m_n_tasks_per_core,
			t_job.m_n_tasks_per_tres,
			t_job.m_n_tasks_per_node,
			t_job.m_n_tasks_per_socket,
			t_job.m_n_tasks_per_board,
			t_job.m_num_cpus,
			t_job.m_num_nodes,
			t_job.m_script,
			t_job.m_command,
			t_job.m_job_state,
			t_job.m_job_state_reason,
			t_job.m_job_state_description,
			t_job_template_submission.id,
			t_job_template_submission.m_template_id
		FROM t_job
		LEFT JOIN t_job_template_submission ON t_job_template_submission.m_job_id = t_job.m_job_id
		ORDER BY t_job.id DESC;
	`
	db := db.GetDbConnection()
	defer db.Close()
	rows, err := db.Query(sqlString)
	if err != nil {
		log.Printf("WARN: queryAllJobs: " + err.Error())
	}
	if rows == nil {
		return jobs
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
			t_job.id,
			t_job.m_job_id,
			t_job.m_user_id,
			t_job.m_accrue_time,
			t_job.m_eligible_time,
			t_job.m_end_time,
			t_job.m_preempt_time,
			t_job.m_preemptable_time,
			t_job.m_resize_time,
			t_job.m_start_time,
			t_job.m_suspend_time,
			t_job.m_work_dir,
			t_job.m_n_tasks_per_core,
			t_job.m_n_tasks_per_tres,
			t_job.m_n_tasks_per_node,
			t_job.m_n_tasks_per_socket,
			t_job.m_n_tasks_per_board,
			t_job.m_num_cpus,
			t_job.m_num_nodes,
			t_job.m_script,
			t_job.m_command,
			t_job.m_job_state,
			t_job.m_job_state_reason,
			t_job.m_job_state_description,
			t_job_template_submission.id,
			t_job_template_submission.m_template_id
		FROM t_job
		LEFT JOIN t_job_template_submission ON t_job_template_submission.m_job_id = t_job.m_job_id
		WHERE m_user_id = :m_user_id;
	`
	db := db.GetDbConnection()
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
		  t_job.id,
			t_job.m_job_id,
			t_job.m_user_id,
			t_job.m_accrue_time,
			t_job.m_eligible_time,
			t_job.m_end_time,
			t_job.m_preempt_time,
			t_job.m_preemptable_time,
			t_job.m_resize_time,
			t_job.m_start_time,
			t_job.m_suspend_time,
			t_job.m_work_dir,
			t_job.m_n_tasks_per_core,
			t_job.m_n_tasks_per_tres,
			t_job.m_n_tasks_per_node,
			t_job.m_n_tasks_per_socket,
			t_job.m_n_tasks_per_board,
			t_job.m_num_cpus,
			t_job.m_num_nodes,
			t_job.m_script,
			t_job.m_command,
			t_job.m_job_state,
			t_job.m_job_state_reason,
			t_job.m_job_state_description,
			t_job_template_submission.id,
			t_job_template_submission.m_template_id
		FROM t_job
		LEFT JOIN t_job_template_submission ON t_job_template_submission.m_job_id = t_job.m_job_id
		WHERE t_job.m_job_id = :m_job_id;
	`
	db := db.GetDbConnection()
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
