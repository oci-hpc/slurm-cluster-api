CREATE TABLE IF NOT EXISTS t_user (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  m_username TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS t_node (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  m_cluster_name TEXT,
  m_cores INTEGER,
  m_cpus INTEGER,
  m_gres TEXT,
  m_name TEXT UNIQUE,
  m_node_addr TEXT,
  m_node_hostname TEXT,
  m_node_state INTEGER,
  m_port INTEGER,
  m_reason TEXT,
  m_reason_time DATETIME,
  m_sockets INTEGER,
  m_threads INTEGER,
  m_version TEXT,
  m_last_seen_time DATETIME
);

CREATE TABLE IF NOT EXISTS t_job (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  m_job_id INTEGER,
  m_user_id INTEGER,
  m_accrue_time DATETIME,
  m_eligible_time DATETIME,
  m_end_time DATETIME,
  m_preempt_time DATETIME,
  m_preemptable_time DATETIME,
  m_resize_time DATETIME,
  m_start_time DATETIME,
  m_submit_time DATETIME,
  m_suspend_time DATETIME,
  m_work_dir TEXT,
  m_n_tasks_per_core INTEGER,
  m_n_tasks_per_tres INTEGER,
  m_n_tasks_per_node INTEGER,
  m_n_tasks_per_socket INTEGER,
  m_n_tasks_per_board INTEGER,
  m_num_cpus INTEGER,
  m_num_nodes INTEGER,
  m_script TEXT,
  m_command TEXT,
  m_job_state INT,
  m_job_state_reason INT,
  m_job_state_description TEXT,
  FOREIGN KEY (m_user_id) REFERENCES t_user (id)
);